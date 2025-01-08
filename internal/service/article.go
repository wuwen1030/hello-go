package service

import (
	"errors"

	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/repository"
)

var (
	ErrTitleRequired   = errors.New("title is required")
	ErrContentRequired = errors.New("content is required")
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleService struct {
	repo *repository.ArticleRepository
}

func NewArticleService(repo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

func (s *ArticleService) Create(req *CreateArticleRequest) (*model.Article, error) {
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Content == "" {
		return nil, ErrContentRequired
	}

	article := &model.Article{
		Title:   req.Title,
		Content: req.Content,
		Status:  1, // 默认为草稿状态
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) Get(id uint) (*model.Article, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrArticleNotFound
	}
	return article, nil
}

func (s *ArticleService) List(page, pageSize int) ([]*model.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

func (s *ArticleService) Update(id uint, req *UpdateArticleRequest) (*model.Article, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Status != 0 {
		article.Status = req.Status
	}

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return ErrArticleNotFound
	}
	return s.repo.Delete(id)
}
