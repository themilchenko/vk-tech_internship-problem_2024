package httpModels

type MovieWithoutCastList struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
}
