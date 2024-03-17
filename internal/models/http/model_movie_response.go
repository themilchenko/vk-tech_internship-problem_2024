package httpModels

type MovieResponse struct {
	Id          int32   `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	ReleaseDate string  `json:"releaseDate,omitempty"`
	Rating      float32 `json:"rating,omitempty"`
	CastList    []Actor `json:"castList,omitempty"`
}
