package httpModels

type MovieWithIdCast struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
	CastIDList  []int32 `json:"castIDList"`
}
