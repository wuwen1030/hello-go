package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wuwen/hello-go/internal/handler"
	"github.com/wuwen/hello-go/internal/middleware"
	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/pkg/auth"
	"github.com/wuwen/hello-go/internal/pkg/config"
	"github.com/wuwen/hello-go/internal/pkg/database"
	"github.com/wuwen/hello-go/internal/repository"
	"github.com/wuwen/hello-go/internal/router"
	"github.com/wuwen/hello-go/internal/router/api"
	"github.com/wuwen/hello-go/internal/service"
)

type App struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	a.config = cfg

	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	dbClient, err := database.NewDBClient(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	db, err := dbClient.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %v", err)
	}

	// 自动迁移数据库表
	// 1. 先迁移 Role 表，因为它被 User 表引用
	if err := db.AutoMigrate(&model.Role{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 创建默认角色
	defaultRole := model.Role{
		ID:          1,
		Name:        "user",
		Permissions: model.Permissions{"article.create", "article.update", "article.delete", "article.read"},
	}
	// 如果不存在则创建
	result := db.FirstOrCreate(&defaultRole, model.Role{ID: 1})
	if result.Error != nil {
		return fmt.Errorf("failed to create default role: %v", result.Error)
	}

	// 2. 再迁移 User 表，因为它依赖 Role 表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 3. 最后迁移其他表
	if err := db.AutoMigrate(&model.Article{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 初始化认证
	auth.Initialize(a.config.JWT.Secret, a.config.JWT.ExpireTime)

	// 初始化依赖
	a.setupDependencies(db)

	return nil
}

func (a *App) setupDependencies(db *gorm.DB) {
	// 创建 gin 引擎
	r := gin.Default()

	// 应用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LimiterMiddleware(time.Second, 100)) // 每秒最多 100 个请求

	// role
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	// article
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, roleRepo)
	userHandler := handler.NewUserHandler(userService, roleService)

	// 注册路由
	a.setupRoutes(r, articleHandler, userHandler, roleHandler)

	// 创建 HTTP 服务器
	a.router = r
	a.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func (a *App) setupRoutes(r *gin.Engine, articleHandler *handler.ArticleHandler,
	userHandler *handler.UserHandler, roleHandler *handler.RoleHandler) {
	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 分组路由
	publicGroup := r.Group("/api/v1")
	authGroup := r.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware())

	// 路由注册
	routers := []router.Router{
		api.NewHealthRouter(),
		api.NewUserRouter(userHandler, roleHandler),
		api.NewRoleRouter(roleHandler),
		api.NewArticleRouter(articleHandler),
	}
	for _, r := range routers {
		r.Register(publicGroup, authGroup)
	}
}

func (a *App) Run() error {
	// 启动服务器
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}
