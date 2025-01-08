package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/pkg/response"
	"github.com/wuwen/hello-go/internal/service"
)

type ArticleHandler struct {
	svc *service.ArticleService
}

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// @Summary     Create article
// @Description Create a new article
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       article body     service.CreateArticleRequest true "Article info"
// @Success     200    {object} response.Response{data=model.Article}
// @Failure     400    {object} response.Response
// @Failure     500    {object} response.Response
// @Router      /articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	article, err := h.svc.Create(&req)
	if err != nil {
		switch err {
		case service.ErrTitleRequired, service.ErrContentRequired:
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, article)
}

// @Summary     Get article
// @Description Get article by ID
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       id  path     int true "Article ID"
// @Success     200 {object} response.Response{data=model.Article}
// @Failure     404 {object} response.Response
// @Failure     500 {object} response.Response
// @Router      /articles/{id} [get]
func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid article id")
		return
	}

	article, err := h.svc.Get(uint(id))
	if err != nil {
		switch err {
		case service.ErrArticleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, article)
}

// @Summary     List articles
// @Description Get articles with pagination
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       page      query    int false "Page number"
// @Param       page_size query    int false "Page size"
// @Success     200      {object} response.Response{data=response.ListResponse{items=[]model.Article}}
// @Failure     500      {object} response.Response
// @Router      /articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(c, gin.H{
		"items": articles,
		"total": total,
	})
}

// @Summary     Update article
// @Description Update article by ID
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       id      path     int                      true "Article ID"
// @Param       article body     service.UpdateArticleRequest true "Article info"
// @Success     200    {object} response.Response{data=model.Article}
// @Failure     400    {object} response.Response
// @Failure     404    {object} response.Response
// @Failure     500    {object} response.Response
// @Router      /articles/{id} [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid article id")
		return
	}

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	article, err := h.svc.Update(uint(id), &req)
	if err != nil {
		switch err {
		case service.ErrArticleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, article)
}

// @Summary     Delete article
// @Description Delete article by ID
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       id  path     int true "Article ID"
// @Success     200 {object} response.Response
// @Failure     404 {object} response.Response
// @Failure     500 {object} response.Response
// @Router      /articles/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid article id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		switch err {
		case service.ErrArticleNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.Success(c, nil)
}
