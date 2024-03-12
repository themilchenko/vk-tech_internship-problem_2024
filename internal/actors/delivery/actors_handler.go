package httpActors

import (
	"fmt"
	"html"
	"net/http"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
)

type ActorsHandler struct {
	actorsUsecase domain.ActorsUsecase
}

func NewActorsUsecase(a domain.ActorsUsecase) ActorsHandler {
	return ActorsHandler{
		actorsUsecase: a,
	}
}

func (h *ActorsHandler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
