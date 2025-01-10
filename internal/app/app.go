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

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/wuwen/hello-go/internal/handler"
	"github.com/wuwen/hello-go/internal/middleware"
	"github.com/wuwen/hello-go/internal/pkg/config"
	"github.com/wuwen/hello-go/internal/repository"
	"github.com/wuwen/hello-go/internal/router"
	"github.com/wuwen/hello-go/internal/router/api"
	"github.com/wuwen/hello-go/internal/service"
)

type App struct {
	config   *config.Config
	router   *gin.Engine
	server   *http.Server
	enforcer *casbin.Enforcer
}

func New() *App {
	return &App{}
}

func (a *App) setupDependencies(db *gorm.DB) {
	// 创建 gin 引擎
	r := gin.Default()

	// 应用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LimiterMiddleware(time.Second, 100))

	// 初始化策略服务
	policyService := service.NewPolicyService(a.enforcer)

	// 初始化角色服务
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo, policyService)
	roleHandler := handler.NewRoleHandler(roleService)

	// 初始化文章服务
	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	// 初始化用户服务
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, roleRepo, policyService)
	userHandler := handler.NewUserHandler(userService)

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
	authGroup.Use(middleware.CasbinMiddleware(a.enforcer))

	// 路由注册
	routers := []router.Router{
		api.NewHealthRouter(),
		api.NewUserRouter(userHandler),
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
