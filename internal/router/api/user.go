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

func (r *UserRouter) Register(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	authUsers := privateGroup.Group("/users")
	{
		authUsers.PUT("/:id", r.handler.Update)
	}
	publicUsers := publicGroup.Group("/users")
	{
		publicUsers.POST("/register", r.handler.Register)
		publicUsers.POST("/login", r.handler.Login)
	}
}
