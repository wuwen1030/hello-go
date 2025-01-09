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
	db, err := database.NewDBClient(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.Article{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// user
	if err := db.AutoMigrate(&model.User{}); err != nil {
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

	// article
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 注册路由
	a.setupRoutes(r, articleHandler, userHandler)

	// 创建 HTTP 服务器
	a.router = r
	a.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func (a *App) setupRoutes(r *gin.Engine, articleHandler *handler.ArticleHandler,
	userHandler *handler.UserHandler) {
	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 分组路由
	publicGroup := r.Group("/api/v1")
	authGroup := r.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware())

	// 公开路由
	publicRouters := []router.Router{
		api.NewHealthRouter(),
		api.NewUserRouter(userHandler),
	}

	// 认证路由
	authRouters := []router.Router{
		api.NewArticleRouter(articleHandler),
	}

	// 注册路由
	for _, r := range publicRouters {
		r.Register(publicGroup)
	}
	for _, r := range authRouters {
		r.Register(authGroup)
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
