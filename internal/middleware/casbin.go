package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的用户
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		username := user.(string)

		// 获取请求的路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 移除路径中的 /api/v1 前缀
		obj = strings.TrimPrefix(obj, "/api/v1")

		// 执行 Casbin 鉴权
		ok, err := enforcer.Enforce(username, obj, act)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		if !ok {
			response.Error(c, http.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
