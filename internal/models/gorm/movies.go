package gormModels

import (
	"time"

	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	ID          uint64    `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(150);unique;not null"`
	Description string    `gorm:"type:text;not null"`
	ReleaseDate time.Time `gorm:"not null"`
	Rating      float32   `gorm:"check:rating >= 0 and rating <= 10;not null"`
}

func (m Movie) ToHTTPResponse() httpModels.MovieResponse {
	return httpModels.MovieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: m.ReleaseDate.Format(time.DateOnly),
		Rating:      m.Rating,
	}
}

func (m Movie) ToHTTPMovies() httpModels.MovieWithoutCastList {
	return httpModels.MovieWithoutCastList{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: m.ReleaseDate.Format(time.DateOnly),
		Rating:      m.Rating,
	}
}

// поучаем фильмы
// select * from movies
// where movies.id=1
// только потом список актеров
// select actor.id, actor.name, actor.gender, actor.birth_date
// from actors
// join actor_movie_relations on actor.id=actor_movie_relations.id
// where movie_id=?
// group by actors.id
