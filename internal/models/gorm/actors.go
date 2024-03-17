package gormModels

import (
	"time"

	"gorm.io/gorm"
)

type Actor struct {
	gorm.Model
	ID        uint64
	Name      string `gorm:"unique"`
	Gender    bool
	BirthDate time.Time
}

type ActorMovieRelation struct {
	gorm.Model
	MovieID uint64
	ActorID uint64
}
