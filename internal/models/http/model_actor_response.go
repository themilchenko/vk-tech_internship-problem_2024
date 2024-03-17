package httpModels

type ActorResponse struct {
	Id        int32  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Gender    bool   `json:"gender,omitempty"`
	BirthDate string `json:"birthDate,omitempty"`
}
