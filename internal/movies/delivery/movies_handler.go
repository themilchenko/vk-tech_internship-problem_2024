package httpMovies

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

type ActorsHandler struct {
	moviesUsecase domain.MoviesUsecase
}

func NewActorsUsecase(a domain.MoviesUsecase) ActorsHandler {
	return ActorsHandler{
		moviesUsecase: a,
	}
}

func (h ActorsHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var receivedMovie httpModels.MovieWithIDCast
	if err := json.NewDecoder(r.Body).Decode(&receivedMovie); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieID, err := h.moviesUsecase.CreateMovie(receivedMovie)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseDate, err := json.Marshal(httpModels.ID{ID: movieID})
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseDate)
}

func (h ActorsHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.moviesUsecase.DeleteMovieByID(movieID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(httpModels.EmptyModel)
}

func (h ActorsHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	movie, err := h.moviesUsecase.GetMovieByID(movieID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	responseData, err := json.Marshal(movie)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h ActorsHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	actor := r.URL.Query().Get("actor")

	order := r.URL.Query().Get("ordrer")
	isOrder := true
	if len(order) != 0 {
		var err error
		isOrder, err = strconv.ParseBool(order)
		if err != nil {
			pkg.HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	filter := httpModels.SortBy(r.URL.Query().Get("filter"))
	if len(filter) == 0 {
		filter = "rating"
	} else {
		if filter != "title" && filter != "rating" && filter != "releaseDate" {
			pkg.HandleError(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
	}

	movies, err := h.moviesUsecase.GetMovies(title, actor, filter, isOrder)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(movies)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h ActorsHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var receivedMovie httpModels.MovieWithoutCastList
	if err := json.NewDecoder(r.Body).Decode(&receivedMovie); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	movie, err := h.moviesUsecase.UpdateMovie(receivedMovie, movieID)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseDate, err := json.Marshal(&movie)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseDate)
}

func (h ActorsHandler) AddActorToMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseUint(r.PathValue("movieID"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}
	actorID, err := strconv.ParseUint(r.PathValue("actorID"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.moviesUsecase.AddActorFromMovie(movieID, actorID); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(httpModels.EmptyModel)
}

func (h ActorsHandler) DeleteActorFromMoive(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseUint(r.PathValue("movieID"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}
	actorID, err := strconv.ParseUint(r.PathValue("actorID"), 10, 64)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.moviesUsecase.DeleteActorFromMovie(movieID, actorID); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(httpModels.EmptyModel)
}
