package httpModels

import "time"

type Movie struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"releaseDate"`
	Rating      float32   `json:"rating"`
	CastList    []Actor   `json:"castList"`
}
