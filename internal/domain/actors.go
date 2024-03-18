package domain

import (
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
)

type ActorsUsecase interface {
	CreateActor(actor httpModels.Actor) (uint64, error)
	GetActorByID(actorID uint64) (httpModels.ActorResponse, error)
	UpdateActor(actor httpModels.Actor, actorID uint64) (httpModels.ActorResponse, error)
	DeleteActorByID(actorID uint64) error
	GetActors(pageNum uint64) ([]httpModels.GetActorsResponse, error)
}

type ActorsRepository interface {
	CreateActor(actor gormModels.Actor) (uint64, error)
	UpdateActor(actor gormModels.Actor) (gormModels.Actor, error)
	DeleteActorByID(actorID uint64) error
	GetActorByID(actorID uint64) (gormModels.Actor, error)
	GetActorsFromMovie(movieID uint64) ([]gormModels.Actor, error)
	GetActors(pageNum uint64) ([]gormModels.Actor, error)
}
