package httpModels

type Actor struct {
	Name      string `json:"name"`
	Gender    bool   `json:"gender"`
	BirthDate string `json:"birthDate"`
}
