package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/handler"
)

type UserRouter struct {
	userHandler *handler.UserHandler
}

func NewUserRouter(userHandler *handler.UserHandler) *UserRouter {
	return &UserRouter{
		userHandler: userHandler,
	}
}

func (r *UserRouter) Register(public, auth *gin.RouterGroup) {
	// 公开路由组
	users := public.Group("/users")
	{
		users.POST("/register", r.userHandler.Register)
		users.POST("/login", r.userHandler.Login)
	}

	// 需要认证的路由组
	authUsers := auth.Group("/users")
	{
		authUsers.PUT("/:id", r.userHandler.Update)
		authUsers.PUT("/:id/role", r.userHandler.UpdateUserRole)
	}
}
