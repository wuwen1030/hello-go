package router

import "github.com/gin-gonic/gin"

// Router 接口定义了路由注册的行为
type Router interface {
	Register(r *gin.RouterGroup)
}
