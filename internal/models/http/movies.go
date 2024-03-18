package httpModels

type Movie struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
	CastList    []Actor `json:"castList"`
}

func (m Movie) ToMovieWithoutCastList() MovieWithoutCastList {
	return MovieWithoutCastList{
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: m.ReleaseDate,
		Rating:      m.Rating,
	}
}

type MovieWithoutCastList struct {
	ID          uint64  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
}

type MovieWithIDCast struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseDate string   `json:"releaseDate"`
	Rating      float32  `json:"rating"`
	CastIDList  []uint64 `json:"castIDList"`
}

type MovieResponse struct {
	ID          uint64  `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	ReleaseDate string  `json:"releaseDate,omitempty"`
	Rating      float32 `json:"rating,omitempty"`
	CastList    []Actor `json:"castList,omitempty"`
}

type MovieID struct {
	ID uint64 `json:"id,omitempty"`
}
