package domain

import (
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
)

type MoviesUsecase interface {
	CreateMovie(movie httpModels.MovieWithIDCast) (uint64, error)
	UpdateMovie(
		movie httpModels.MovieWithoutCastList,
		movieID uint64,
	) (httpModels.MovieResponse, error)
	GetMovieByID(movieID uint64) (httpModels.MovieResponse, error)
	DeleteMovieByID(movieID uint64) error
	DeleteActorFromMovie(movieID, actorID uint64) error
	AddActorFromMovie(movieID, actorID uint64) error
	GetMovies(
		title, actorName string,
		sortBy httpModels.SortBy,
		order bool,
	) ([]httpModels.MovieResponse, error)
}

type MoviesRepository interface {
	CreateMovieWithoutCastList(movie gormModels.Movie) (uint64, error)
	CreateMovieWithCastList(movie gormModels.Movie, castList []uint64) (uint64, error)
	UpdateMovie(movie gormModels.Movie) (gormModels.Movie, error)
	GetMovieByID(movieID uint64) (gormModels.Movie, error)
	DeleteMovieByID(movieID uint64) error
	DeleteActorFromMovie(movieID, actorID uint64) error
	AddActorToMovie(movieID, actorID uint64) error
	GetMoviesOfActor(actorID uint64) ([]gormModels.Movie, error)
	GetMovies(
		title, actorName string,
		sortBy httpModels.SortBy,
		order bool,
	) ([]gormModels.Movie, error)
}
