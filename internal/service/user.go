package service

import (
	"errors"
	"time"

	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/pkg/auth"
	"github.com/wuwen/hello-go/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidAuth  = errors.New("invalid username or password")
	ErrUserExist    = errors.New("user already exists")
	ErrRoleNotFound = errors.New("role not found")
)

type UserService struct {
	repo     *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(repo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{repo: repo, roleRepo: roleRepo}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

type UpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty"`
}

func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	// check if user exists
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, ErrUserExist
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Status:   model.UserStatusActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidAuth
	}

	if !user.ValidatePassword(req.Password) {
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

func (s *UserService) UpdateUser(id uint, req *UpdateRequest) (*model.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			return nil, err
		}
	}

	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

func (s *UserService) UpdateUserRole(userID uint, roleID uint) (*model.User, error) {
	user, err := s.repo.FindById(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	user.RoleID = role.ID
	user.UpdatedAt = time.Now()

	return s.repo.UpdateRole(user)
}
