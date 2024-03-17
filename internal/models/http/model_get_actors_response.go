package httpModels

type GetActorsResponse struct {
	Actor        *ActorResponse         `json:"actor,omitempty"`
	ActedInFilms []MovieWithoutCastList `json:"actedInFilms,omitempty"`
}
