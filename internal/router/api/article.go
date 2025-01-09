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

func (r *ArticleRouter) Register(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	authArticles := privateGroup.Group("/articles")
	{
		// 所有接口都需要认证，因为使用的是 authGroup
		authArticles.POST("", r.handler.Create)
		authArticles.PUT("/:id", r.handler.Update)
		authArticles.DELETE("/:id", r.handler.Delete)
	}
	publicArticles := publicGroup.Group("/articles")
	{
		publicArticles.GET("/:id", r.handler.Get)
		publicArticles.GET("", r.handler.List)
	}
}
