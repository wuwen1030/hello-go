package model

import "gorm.io/gorm"

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

type User struct {
	gorm.Model
	Username string     `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password string     `gorm:"size:100;not null" json:"-"`
	Email    string     `gorm:"size:100;uniqueIndex" json:"email"`
	Status   UserStatus `gorm:"default:1" json:"status"`
}
