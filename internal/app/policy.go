package app

import "fmt"

func (a *App) initializeDefaultPolicies() error {
	// 如果策略存在，直接返回
	if p, err := a.enforcer.GetPolicy(); err != nil || len(p) > 0 {
		return nil
	}

	policies := map[string][][]string{
		"role:user": {
			{"/api/v1/articles", "GET"},
			{"/api/v1/articles", "POST"},
			{"/api/v1/articles/*", "PUT"},
			{"/api/v1/articles/*", "DELETE"},
			{"/api/v1/users", "GET"},
			{"/api/v1/users", "POST"},
			{"/api/v1/users/*", "PUT"},
			{"/api/v1/users/*", "DELETE"},
		},
		"role:admin": {
			{"/api/v1/articles", "GET"},
			{"/api/v1/articles", "POST"},
			{"/api/v1/articles/*", "PUT"},
			{"/api/v1/articles/*", "DELETE"},
			{"/api/v1/roles", "GET"},
			{"/api/v1/roles", "POST"},
			{"/api/v1/roles/*", "PUT"},
			{"/api/v1/roles/*", "DELETE"},
		},
	}

	for role, rules := range policies {
		for _, rule := range rules {
			if _, err := a.enforcer.AddPolicy(role, rule[0], rule[1]); err != nil {
				return fmt.Errorf("failed to add %s policy: %v", role, err)
			}
		}
	}

	if err := a.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save casbin policy: %v", err)
	}

	// 为管理员用户分配管理员角色
	if _, err := a.enforcer.AddGroupingPolicy("user:admin", "role:admin"); err != nil {
		return fmt.Errorf("failed to assign admin role to admin user: %v", err)
	}

	return nil
}
