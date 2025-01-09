package model

import (
	"time"
)

type Article struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-07-20T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-07-20T10:00:00Z"`
	Title     string    `gorm:"size:200;not null" json:"title" example:"文章标题"`
	Content   string    `gorm:"type:text" json:"content" example:"文章内容"`
	Status    int       `gorm:"default:1" json:"status" example:"1"` // 1:draft 2:published
}
