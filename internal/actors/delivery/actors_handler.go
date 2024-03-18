package httpActors

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

const (
	defaultPage = 1
)

type ActorsHandler struct {
	actorsUsecase domain.ActorsUsecase
}

func NewActorsUsecase(a domain.ActorsUsecase) ActorsHandler {
	return ActorsHandler{
		actorsUsecase: a,
	}
}

func (h ActorsHandler) CreateActor(w http.ResponseWriter, r *http.Request) {
	var receivedActor httpModels.Actor
	if err := json.NewDecoder(r.Body).Decode(&receivedActor); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	actorID, err := h.actorsUsecase.CreateActor(receivedActor)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
	}

	responseData, err := json.Marshal(httpModels.ID{ID: actorID})
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	actorID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.actorsUsecase.DeleteActorByID(actorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(httpModels.EmptyModel)
}

func (h ActorsHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	actorID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
	}

	actor, err := h.actorsUsecase.GetActorByID(actorID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	responseData, err := json.Marshal(actor)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h ActorsHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	actorID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var receivedActor httpModels.Actor
	if err := json.NewDecoder(r.Body).Decode(&receivedActor); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedActor, err := h.actorsUsecase.UpdateActor(receivedActor, actorID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	responseData, err := json.Marshal(updatedActor)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h ActorsHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	var pageNum uint64
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageNum = defaultPage
	}

	pageNum, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, err := h.actorsUsecase.GetActors(pageNum)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(actors)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
