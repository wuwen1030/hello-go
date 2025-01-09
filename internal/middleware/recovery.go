package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				response.Error(c, http.StatusInternalServerError, "internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
