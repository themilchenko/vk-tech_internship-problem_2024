package actorsRepository

import (
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB

	pageSize uint64
}

func NewPostgres(url string, ps uint64) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		gormModels.Actor{},
	)

	return &Postgres{
		DB:       db,
		pageSize: ps,
	}, nil
}

func (db Postgres) CreateActor(actor gormModels.Actor) (uint64, error) {
	var recievedActor gormModels.Actor
	if err := db.DB.Create(&actor).Scan(&recievedActor).Error; err != nil {
		return 0, err
	}
	return recievedActor.ID, nil
}

func (db Postgres) UpdateActor(actor gormModels.Actor) (gormModels.Actor, error) {
	var recievedActor gormModels.Actor
	if err := db.DB.Model(&gormModels.Actor{ID: actor.ID}).
		Updates(actor).
		Scan(&recievedActor).
		Error; err != nil {
		return gormModels.Actor{}, err
	}
	return gormModels.Actor{}, nil
}

func (db Postgres) DeleteActorByID(actorID uint64) error {
	if err := db.DB.Unscoped().
		Delete(&gormModels.Actor{}, "id = ?", actorID).
		Error; err != nil {
		return err
	}
	return nil
}

func (db Postgres) GetActorByID(actorID uint64) (gormModels.Actor, error) {
	var recievedActor gormModels.Actor
	if err := db.DB.First(&gormModels.Actor{ID: actorID}).
		Scan(&recievedActor).
		Error; err != nil {
		return gormModels.Actor{}, nil
	}
	return recievedActor, nil
}

func (db Postgres) GetActorsFromMovie(movieID uint64) ([]gormModels.Actor, error) {
	var recievedActors []gormModels.Actor
	if err := db.DB.Model(&gormModels.ActorMovieRelation{}).
		Joins("JOIN actors ON actors.id=actor_movie_relations.actor_id").
		Where("movie_id = ?", movieID).
		Select("actors.id, actors.name, actors.birth_date").
		Find(&recievedActors).
		Error; err != nil {
		return recievedActors, err
	}
	return recievedActors, nil
}

func (db Postgres) GetActors(pageNum uint64) ([]gormModels.Actor, error) {
	offset := db.pageSize * (pageNum - 1)
	var recievedActors []gormModels.Actor

	if err := db.DB.
		Offset(int(offset)).
		Limit(int(db.pageSize)).
		Order("name").
		Find(&recievedActors).
		Error; err != nil {
		return nil, err
	}
	return recievedActors, nil
}
