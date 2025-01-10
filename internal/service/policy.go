package service

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

type PolicyService struct {
	enforcer *casbin.Enforcer
}

var (
	RolePrefix = "role:"
	UserPrefix = "user:"
)

func NewPolicyService(enforcer *casbin.Enforcer) *PolicyService {
	return &PolicyService{
		enforcer: enforcer,
	}
}

// AddRoleForUser 为用户分配角色
func (s *PolicyService) AddRoleForUser(username, role string) error {
	_, err := s.enforcer.AddGroupingPolicy(UserPrefix+username, RolePrefix+role)
	if err != nil {
		return fmt.Errorf("failed to add role for user: %v", err)
	}
	return nil
}

// RemoveRoleForUser 移除用户的角色
func (s *PolicyService) RemoveRoleForUser(username, role string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(UserPrefix+username, RolePrefix+role)
	if err != nil {
		return fmt.Errorf("failed to remove role for user: %v", err)
	}
	return nil
}

// AddPolicy 添加权限策略
func (s *PolicyService) AddPolicy(role, path, method string) error {
	_, err := s.enforcer.AddPolicy(RolePrefix+role, path, method)
	if err != nil {
		return fmt.Errorf("failed to add policy: %v", err)
	}
	return nil
}

// RemovePolicy 移除权限策略
func (s *PolicyService) RemovePolicy(role, path, method string) error {
	_, err := s.enforcer.RemovePolicy(RolePrefix+role, path, method)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %v", err)
	}
	return nil
}

// GetRolesForUser 获取用户的所有角色
func (s *PolicyService) GetRolesForUser(username string) ([]string, error) {
	return s.enforcer.GetRolesForUser(username)
}

// GetPermissionsForRole 获取角色的所有权限
func (s *PolicyService) GetPermissionsForRole(role string) ([][]string, error) {
	return s.enforcer.GetPermissionsForUser(RolePrefix + role)
}

// UpdateRoleName 更新角色名称
func (s *PolicyService) UpdateRoleName(oldName, newName string) error {
	// 更新策略规则中的角色名
	oldPolicies, err := s.enforcer.GetFilteredPolicy(0, RolePrefix+oldName)
	if err != nil {
		return err
	}
	for _, p := range oldPolicies {
		_, err := s.enforcer.UpdatePolicy(
			[]string{RolePrefix + oldName, p[1], p[2]},
			[]string{RolePrefix + newName, p[1], p[2]},
		)
		if err != nil {
			return err
		}
	}

	// 更新角色关系中的角色名
	oldRoles, err := s.enforcer.GetFilteredGroupingPolicy(1, RolePrefix+oldName)
	if err != nil {
		return err
	}
	for _, g := range oldRoles {
		_, err := s.enforcer.UpdateGroupingPolicy(
			[]string{g[0], RolePrefix + oldName},
			[]string{g[0], RolePrefix + newName},
		)
		if err != nil {
			return err
		}
	}

	return s.enforcer.SavePolicy()
}
