package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/handler"
)

type RoleRouter struct {
	handler *handler.RoleHandler
}

func NewRoleRouter(handler *handler.RoleHandler) *RoleRouter {
	return &RoleRouter{
		handler: handler,
	}
}

func (r *RoleRouter) Register(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	authRoles := privateGroup.Group("/roles")
	{
		authRoles.POST("", r.handler.Create)
		authRoles.GET("/:id", r.handler.Get)
		authRoles.PUT("/:id", r.handler.Update)
		authRoles.DELETE("/:id", r.handler.Delete)
	}
}
