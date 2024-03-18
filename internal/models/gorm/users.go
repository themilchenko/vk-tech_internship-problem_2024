package gormModels

import (
	"time"

	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint64
	Username string `gorm:"unique"`
	Password string
	Role     string
}

func (u User) ToHTTPModel() httpModels.AuthUser {
	return httpModels.AuthUser{
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
	}
}

type Session struct {
	gorm.Model
	UserID     uint64
	SessionID  string
	ExpireDate time.Time
}
