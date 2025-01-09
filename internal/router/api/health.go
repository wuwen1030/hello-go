package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
)

type HealthRouter struct{}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) Register(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	publicGroup.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}
