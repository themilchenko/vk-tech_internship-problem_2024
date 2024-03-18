package httpModels

type GetActorsResponse struct {
	Actor        ActorResponse          `json:"actor,omitempty"`
	ActedInFilms []MovieWithoutCastList `json:"actedInFilms,omitempty"`
}

type Actor struct {
	Name      string `json:"name"`
	Gender    bool   `json:"gender"`
	BirthDate string `json:"birthDate"`
}

type ActorResponse struct {
	ID        uint64 `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Gender    bool   `json:"gender"`
	BirthDate string `json:"birthDate,omitempty"`
}

type ActorID struct {
	ID uint64 `json:"id,omitempty"`
}
