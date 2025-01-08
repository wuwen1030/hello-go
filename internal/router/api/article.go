package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/handler"
)

type ArticleRouter struct {
	handler *handler.ArticleHandler
}

func NewArticleRouter(handler *handler.ArticleHandler) *ArticleRouter {
	return &ArticleRouter{
		handler: handler,
	}
}

func (r *ArticleRouter) Register(group *gin.RouterGroup) {
	articles := group.Group("/articles")
	{
		articles.POST("", r.handler.Create)
		articles.GET("/:id", r.handler.Get)
		articles.GET("", r.handler.List)
		articles.PUT("/:id", r.handler.Update)
		articles.DELETE("/:id", r.handler.Delete)
	}
}
