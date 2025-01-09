package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/handler"
)

type UserRouter struct {
	handler *handler.UserHandler
}

func NewUserRouter(handler *handler.UserHandler) *UserRouter {
	return &UserRouter{
		handler: handler,
	}
}

func (r *UserRouter) Register(group *gin.RouterGroup) {
	users := group.Group("/users")
	{
		users.POST("/register", r.handler.Register)
		users.POST("/login", r.handler.Login)
	}
}
