package model

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

// Role 角色模型
type Role struct {
	ID          uint        `gorm:"primarykey;autoIncrement" json:"id" example:"1"`
	CreatedAt   time.Time   `json:"created_at" example:"2024-07-20T10:00:00Z"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2024-07-20T10:00:00Z"`
	Name        string      `gorm:"size:50;not null;uniqueIndex" json:"name" example:"admin"`
	Permissions Permissions `gorm:"type:text;not null" json:"permissions" example:"article.create,article.update,article.delete,article.read"`
}

// Permissions 自定义类型，用于处理权限的序列化和反序列化
type Permissions []string

// Value 实现 driver.Valuer 接口，用于将 Permissions 转换为数据库值
func (p Permissions) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, errors.New("permissions cannot be empty")
	}
	return strings.Join(p, ","), nil
}

// Scan 实现 sql.Scanner 接口，用于将数据库值转换为 Permissions
func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		*p = make(Permissions, 0)
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New("failed to convert permissions value to string")
		}
		str = string(bytes)
	}

	if str == "" {
		*p = make(Permissions, 0)
		return nil
	}

	*p = strings.Split(str, ",")
	return nil
}
