package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/handler"
)

type UserRouter struct {
	handler     *handler.UserHandler
	roleHandler *handler.RoleHandler
}

func NewUserRouter(handler *handler.UserHandler, roleHandler *handler.RoleHandler) *UserRouter {
	return &UserRouter{
		handler:     handler,
		roleHandler: roleHandler,
	}
}

func (r *UserRouter) Register(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	authUsers := privateGroup.Group("/users")
	{
		authUsers.PUT("/:id", r.handler.Update)
		authUsers.PUT("/:id/roles/:role_id", r.handler.UpdateRole)
	}
	publicUsers := publicGroup.Group("/users")
	{
		publicUsers.POST("/register", r.handler.Register)
		publicUsers.POST("/login", r.handler.Login)
	}
}
