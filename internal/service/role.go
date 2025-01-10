package service

import (
	"errors"
	"time"

	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/repository"
)

var (
	ErrRoleExist    = errors.New("role already exists")
	ErrRoleNotFound = errors.New("role not found")
)

type RoleService struct {
	roleRepo      *repository.RoleRepository
	policyService *PolicyService
}

func NewRoleService(roleRepo *repository.RoleRepository, policyService *PolicyService) *RoleService {
	return &RoleService{
		roleRepo:      roleRepo,
		policyService: policyService,
	}
}

type PolicyRule struct {
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

type CreateRoleRequest struct {
	Name     string       `json:"name" binding:"required"`
	Policies []PolicyRule `json:"policies" binding:"required"`
}

func (s *RoleService) Create(req *CreateRoleRequest) (*model.Role, error) {
	// 检查角色是否存在
	_, err := s.roleRepo.FindByName(req.Name)
	if err == nil {
		return nil, ErrRoleExist
	}

	// 创建角色
	role := &model.Role{
		Name: req.Name,
	}

	role, err = s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	// 添加权限策略
	for _, policy := range req.Policies {
		if err := s.policyService.AddPolicy(role.Name, policy.Path, policy.Method); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *RoleService) Get(id uint) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

type UpdateRoleRequest struct {
	Name     string       `json:"name" binding:"omitempty"`
	Policies []PolicyRule `json:"policies" binding:"omitempty"`
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if req.Name != "" && req.Name != role.Name {
		// 如果角色名变更，需要更新所有相关的策略
		oldName := role.Name
		role.Name = req.Name
		role.UpdatedAt = time.Now()

		role, err = s.roleRepo.Update(role)
		if err != nil {
			return nil, err
		}

		// 更新策略中的角色名
		if err := s.policyService.UpdateRoleName(oldName, role.Name); err != nil {
			return nil, err
		}
	}

	if req.Policies != nil {
		if err := s.UpdatePermissions(id, req.Policies); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return ErrRoleNotFound
	}

	// 先删除角色相关的所有策略
	permissions, err := s.policyService.GetPermissionsForRole(role.Name)
	if err != nil {
		return err
	}

	for _, p := range permissions {
		if err := s.policyService.RemovePolicy(role.Name, p[1], p[2]); err != nil {
			return err
		}
	}

	return s.roleRepo.Delete(id)
}

func (s *RoleService) UpdatePermissions(roleID uint, policies []PolicyRule) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return err
	}

	// 移除旧的权限
	oldPermissions, err := s.policyService.GetPermissionsForRole(role.Name)
	if err != nil {
		return err
	}
	for _, p := range oldPermissions {
		if err := s.policyService.RemovePolicy(role.Name, p[1], p[2]); err != nil {
			return err
		}
	}

	// 添加新的权限
	for _, policy := range policies {
		if err := s.policyService.AddPolicy(role.Name, policy.Path, policy.Method); err != nil {
			return err
		}
	}

	return nil
}
