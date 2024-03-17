package httpActors

import (
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

func CreateActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetActors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func UpdateActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
