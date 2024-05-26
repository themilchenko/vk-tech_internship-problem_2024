package gormModels

import (
	"time"

	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type Actor struct {
	gorm.Model
	ID        uint64
	Name      string `gorm:"unique"`
	Gender    bool
	BirthDate time.Time
}

func (a Actor) ToHTTPModel() httpModels.ActorResponse {
	return httpModels.ActorResponse{
		ID:        a.ID,
		Name:      a.Name,
		Gender:    a.Gender,
		BirthDate: a.BirthDate.Format(time.DateOnly),
	}
}

type ActorMovieRelation struct {
	gorm.Model
	MovieID uint64 `gorm:"uniqueIndex:idx_movie_actor"`
	ActorID uint64 `gorm:"uniqueIndex:idx_movie_actor"`
}

func (a *Actor) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("actor_id = ?", a.ID).Unscoped().Delete(&ActorMovieRelation{})
	return
}
