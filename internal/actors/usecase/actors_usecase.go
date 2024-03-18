package actorsUsecase

import (
	"errors"
	"time"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type ActorsUsecase struct {
	actorsRepository domain.ActorsRepository
	moviesRepository domain.MoviesRepository
}

func NewActorsUsecase(a domain.ActorsRepository, m domain.MoviesRepository) ActorsUsecase {
	return ActorsUsecase{
		actorsRepository: a,
		moviesRepository: m,
	}
}

func (u ActorsUsecase) CreateActor(actor httpModels.Actor) (uint64, error) {
	t, err := time.Parse(time.DateOnly, actor.BirthDate)
	if err != nil {
		return 0, err
	}

	actorID, err := u.actorsRepository.CreateActor(gormModels.Actor{
		Name:      actor.Name,
		BirthDate: t,
		Gender:    actor.Gender,
	})
	if err != nil {
		return 0, err
	}

	return actorID, nil
}

func (u ActorsUsecase) GetActorByID(actorID uint64) (httpModels.ActorResponse, error) {
	actor, err := u.actorsRepository.GetActorByID(actorID)
	if err != nil {
		return httpModels.ActorResponse{}, err
	}
	return actor.ToHTTPModel(), nil
}

func (u ActorsUsecase) UpdateActor(
	actor httpModels.Actor,
	actorID uint64,
) (httpModels.ActorResponse, error) {
	t, err := time.Parse(time.DateOnly, actor.BirthDate)
	if err != nil {
		return httpModels.ActorResponse{}, err
	}

	updatedActor, err := u.actorsRepository.UpdateActor(gormModels.Actor{
		ID:        actorID,
		Name:      actor.Name,
		Gender:    actor.Gender,
		BirthDate: t,
	})
	if err != nil {
		return httpModels.ActorResponse{}, err
	}
	return updatedActor.ToHTTPModel(), nil
}

func (u ActorsUsecase) DeleteActorByID(actorID uint64) error {
	return u.actorsRepository.DeleteActorByID(actorID)
}

func (u ActorsUsecase) GetActors(pageNum uint64) ([]httpModels.GetActorsResponse, error) {
	actors, err := u.actorsRepository.GetActors(pageNum)
	if err != nil {
		return []httpModels.GetActorsResponse{}, err
	}

	responseActors := make([]httpModels.GetActorsResponse, len(actors))
	for i, v := range actors {
		responseActors[i].Actor = v.ToHTTPModel()
		movies, err := u.moviesRepository.GetMoviesOfActor(v.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return []httpModels.GetActorsResponse{}, err
		}

		httpMovies := make([]httpModels.MovieWithoutCastList, len(movies))
		for i, m := range movies {
			httpMovies[i] = m.ToHTTPMovies()
		}

		responseActors[i].ActedInFilms = httpMovies
	}

	return responseActors, nil
}
