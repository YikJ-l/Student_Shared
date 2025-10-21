package model

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 基础模型，包含共有字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password    string    `gorm:"size:100;not null" json:"-"`
	Email       string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Nickname    string    `gorm:"size:50" json:"nickname"`
	Avatar      string    `gorm:"size:255" json:"avatar"`
	School      string    `gorm:"size:100" json:"school"`
	Department  string    `gorm:"size:100" json:"department"`
	Major       string    `gorm:"size:100" json:"major"`
	Introduction string    `gorm:"size:500" json:"introduction"`
	LastLogin   time.Time `json:"last_login"`
	Role        string    `gorm:"size:20;default:'student'" json:"role"` // student, teacher, admin
}