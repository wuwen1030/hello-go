package model

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-07-20T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-07-20T10:00:00Z"`
	Name      string    `gorm:"size:50;not null;uniqueIndex" json:"name" example:"admin"`
}
