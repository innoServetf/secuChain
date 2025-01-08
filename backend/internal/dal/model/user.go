package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"size:100;not null" json:"-"`
	Email       string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// UserVerification 邮箱验证码模型
type UserVerification struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	Code      string    `gorm:"size:6;not null"`
	Type      string    `gorm:"size:20;not null"` // register/reset_password
	ExpiredAt time.Time `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt time.Time
}
