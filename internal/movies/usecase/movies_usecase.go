package moviesUsecase

import (
	"time"

	"github.com/pkg/errors"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type MoviesUsecase struct {
	moviesRepository domain.MoviesRepository
	actorsRepository domain.ActorsRepository
}

func NewMoviesUsecase(m domain.MoviesRepository, a domain.ActorsRepository) MoviesUsecase {
	return MoviesUsecase{
		moviesRepository: m,
		actorsRepository: a,
	}
}

func (u MoviesUsecase) CreateMovie(movie httpModels.MovieWithIDCast) (uint64, error) {
	t, err := time.Parse(time.DateOnly, movie.ReleaseDate)
	if err != nil {
		return 0, err
	}

	gormMovie := gormModels.Movie{
		Title:       movie.Title,
		Description: movie.Description,
		ReleaseDate: t,
		Rating:      movie.Rating,
	}
	var movieID uint64

	if len(movie.CastIDList) == 0 {
		movieID, err = u.moviesRepository.CreateMovieWithoutCastList(gormMovie)
		if err != nil {
			return 0, domain.ErrCreate
		}
	} else {
		movieID, err = u.moviesRepository.CreateMovieWithCastList(gormMovie, movie.CastIDList)
		if err != nil {
			return 0, domain.ErrCreate
		}
	}

	return movieID, nil
}

func (u MoviesUsecase) UpdateMovie(
	movie httpModels.MovieWithoutCastList,
	movieID uint64,
) (httpModels.MovieResponse, error) {
	t, err := time.Parse(time.DateOnly, movie.ReleaseDate)
	if err != nil {
		return httpModels.MovieResponse{}, err
	}

	updatedMovie, err := u.moviesRepository.UpdateMovie(gormModels.Movie{
		ID:          movieID,
		Title:       movie.Title,
		Description: movie.Description,
		ReleaseDate: t,
		Rating:      movie.Rating,
	})
	if err != nil {
		return httpModels.MovieResponse{}, err
	}
	return updatedMovie.ToHTTPResponse(), nil
}

func (u MoviesUsecase) GetMovieByID(movieID uint64) (httpModels.MovieResponse, error) {
	movie, err := u.moviesRepository.GetMovieByID(movieID)
	if err != nil {
		return httpModels.MovieResponse{}, err
	}

	castList, err := u.actorsRepository.GetActorsFromMovie(movieID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return httpModels.MovieResponse{}, err
	}
	httpCastList := make([]httpModels.ActorResponse, len(castList))
	for i, a := range castList {
		httpCastList[i] = a.ToHTTPModel()
	}

	httpMovie := movie.ToHTTPResponse()
	httpMovie.CastList = httpCastList
	return httpMovie, nil
}

func (u MoviesUsecase) DeleteMovieByID(movieID uint64) error {
	return u.moviesRepository.DeleteMovieByID(movieID)
}

func (u MoviesUsecase) DeleteActorFromMovie(movieID, actorID uint64) error {
	return u.moviesRepository.DeleteActorFromMovie(movieID, actorID)
}

func (u MoviesUsecase) AddActorFromMovie(movieID, actorID uint64) error {
	return u.moviesRepository.AddActorToMovie(movieID, actorID)
}

func (u MoviesUsecase) GetMovies(
	title, actorName string,
	sortBy httpModels.SortBy,
	order bool,
) ([]httpModels.MovieResponse, error) {
	movies, err := u.moviesRepository.GetMovies(title, actorName, sortBy, order)
	if err != nil {
		return []httpModels.MovieResponse{}, err
	}

	httpMovies := make([]httpModels.MovieResponse, len(movies))
	for i, v := range movies {
		httpMovies[i] = v.ToHTTPResponse()
		actors, err := u.actorsRepository.GetActorsFromMovie(v.ID)
		if err != nil {
			return []httpModels.MovieResponse{}, err
		}

		httpActors := make([]httpModels.ActorResponse, len(actors))
		for j, actor := range actors {
			httpActors[j] = actor.ToHTTPModel()
		}

		httpMovies[i].CastList = httpActors
	}

	return httpMovies, nil
}
