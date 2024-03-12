package httpMovies

import (
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
