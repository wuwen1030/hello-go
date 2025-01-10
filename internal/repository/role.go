package repository

import (
	"github.com/wuwen/hello-go/internal/model"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *model.Role) (*model.Role, error) {
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByName(name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Update(role *model.Role) (*model.Role, error) {
	if err := r.db.Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}
