package moviesRepository

import (
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		gormModels.Movie{},
		gormModels.ActorMovieRelation{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db Postgres) CreateMovie(movie gormModels.Movie) (uint64, error) {
	var recievedMovie gormModels.Movie
	if err := db.DB.Create(&movie).Scan(&recievedMovie).Error; err != nil {
		return 0, err
	}
	return recievedMovie.ID, nil
}

func (db Postgres) UpdateMovie(movie gormModels.Movie) (gormModels.Movie, error) {
	var recievedMovie gormModels.Movie
	if err := db.DB.Model(&gormModels.Movie{ID: movie.ID}).
		Updates(movie).
		Scan(&recievedMovie).
		Error; err != nil {
		return gormModels.Movie{}, err
	}
	return recievedMovie, nil
}

func (db Postgres) GetMovieByID(movieID uint64) (gormModels.Movie, error) {
	var recievedMovie gormModels.Movie
	if err := db.DB.First(&gormModels.Movie{ID: movieID}).
		Scan(&recievedMovie).
		Error; err != nil {
		return gormModels.Movie{}, err
	}
	return recievedMovie, nil
}

func (db Postgres) DeleteMovieByID(movieID uint64) error {
	if err := db.DB.Unscoped().
		Delete(&gormModels.Movie{}, "id = ?", movieID).
		Error; err != nil {
		return err
	}
	return nil
}

func (db Postgres) AddActorToMovie(movieID, actorID uint64) error {
	if err := db.DB.Create(&gormModels.ActorMovieRelation{
		MovieID: movieID,
		ActorID: actorID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (db Postgres) DeleteActorFromMovie(movieID, actorID uint64) error {
	if err := db.DB.Unscoped().Delete(&gormModels.ActorMovieRelation{
		MovieID: movieID,
		ActorID: actorID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (db Postgres) GetMoviesOfActor(actorID uint64) ([]gormModels.Movie, error) {
	var recievedMovies []gormModels.Movie
	if err := db.DB.Model(&gormModels.Movie{}).
		Joins("JOIN actor_movie_relations ON actor_movie_relations.movie_id=movies.id").
		Where("actor_id = ?", actorID).
		Select("movies.id, movies.title, movies.description, movies.release_date, movies.rating").
		Find(&recievedMovies).
		Error; err != nil {
		return recievedMovies, err
	}
	return recievedMovies, nil
}

func (db Postgres) GetMovies(
	title, actorName string,
	sortBy httpModels.SortBy,
	order bool,
) ([]gormModels.Movie, error) {
	var movies []gormModels.Movie

	query := db.DB.Model(&gormModels.Movie{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if actorName != "" {
		query = query.Joins("JOIN actor_movie_relations ON movies.id = actor_movie_relations.movie_id").
			Joins("JOIN actors ON actor_movie_relations.actor_id = actors.id").
			Where("actors.name LIKE ?", "%"+actorName+"%")
	}

	switch sortBy {
	case "title":
		query = query.Order("title")
	case "rating":
		query = query.Order("rating")
	case "release_date":
		query = query.Order("release_date")
	}

	if !order {
		query = query.Order(sortBy + " DESC")
	}

	if err := query.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}
