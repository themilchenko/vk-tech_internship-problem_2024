package gormModels

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	ID          uint64
	Title       string `gorm:"unique"`
	Description string
	ReleaseDate time.Time
	Rating      float32
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
