package model

import (
	"time"
)

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusInactive UserStatus = iota // 0: 未激活
	UserStatusActive                     // 1: 已激活
	UserStatusBanned                     // 2: 已封禁
)

// String 实现 Stringer 接口
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "inactive"
	case UserStatusActive:
		return "active"
	case UserStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

// User 用户模型
type User struct {
	ID        uint       `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2024-07-20T10:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2024-07-20T10:00:00Z"`
	Username  string     `gorm:"size:50;not null;uniqueIndex" json:"username" example:"testuser"`
	Password  string     `gorm:"size:100;not null" json:"-"`
	Email     string     `gorm:"size:100;uniqueIndex" json:"email" example:"test@example.com"`
	Status    UserStatus `gorm:"default:1" json:"status" example:"1"`
}
