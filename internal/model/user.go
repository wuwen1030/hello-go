package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password string `gorm:"size:100;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`
	Status   int    `gorm:"default:1" json:"status"`
}
