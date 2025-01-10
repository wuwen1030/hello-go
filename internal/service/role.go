package service

import (
	"errors"
	"time"

	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/repository"
)

var (
	ErrRoleExist = errors.New("role already exists")
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

func (s *RoleService) Create(req *CreateRoleRequest) (*model.Role, error) {
	// check if role exists
	_, err := s.repo.FindByName(req.Name)
	if err == nil {
		return nil, ErrRoleExist
	}

	role := &model.Role{
		Name:        req.Name,
		Permissions: req.Permissions,
	}

	return s.repo.Create(role)
}

func (s *RoleService) Get(id uint) (*model.Role, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

type UpdateRoleRequest struct {
	Name        string   `json:"name" binding:"omitempty"`
	Permissions []string `json:"permissions" binding:"omitempty"`
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) (*model.Role, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	if req.Permissions != nil {
		role.Permissions = req.Permissions
	}

	role.UpdatedAt = time.Now()

	return s.repo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrRoleNotFound
	}
	return s.repo.Delete(id)
}
