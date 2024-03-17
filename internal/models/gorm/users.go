package gormModels

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint64
	Username string `gorm:"unique"`
	Password string
	Role     string
}

type Session struct {
	gorm.Model
	UserID     uint64
	SessionID  string
	ExpireDate time.Time
}
