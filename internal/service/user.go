package service

import (
	"errors"

	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/pkg/auth"
	"github.com/wuwen/hello-go/internal/repository"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidAuth  = errors.New("invalid username or password")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email format")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	// 检查用户是否已存在
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, ErrUserExists
	}

	// 检查邮箱是否已被使用
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, ErrUserExists
	}

	// 创建新用户
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   model.UserStatusActive,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidAuth
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidAuth
	}

	// 生成 JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}
