package actorsUsecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mockDomain "github.com/themilchenko/vk-tech_internship-problem_2024/internal/mocks/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"go.uber.org/mock/gomock"
)

func TestUsecase_CreateActor(t *testing.T) {
	type mockBehaviorCreateActor func(r *mockDomain.MockActorsRepository, actor gormModels.Actor)

	tests := []struct {
		name                    string
		inputActor              httpModels.Actor
		mockBehaviorCreateActor mockBehaviorCreateActor
		expectedActorID         uint64
		expectedError           error
	}{
		{
			name: "CreateActor success",
			inputActor: httpModels.Actor{
				Name:      "Name",
				BirthDate: "2006-01-02",
				Gender:    true,
			},
			mockBehaviorCreateActor: func(m *mockDomain.MockActorsRepository, actor gormModels.Actor) {
				m.EXPECT().
					CreateActor(actor).
					Return(uint64(1), nil)
			},
			expectedActorID: uint64(1),
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockActorsRepository(ctrl)
			mockMovieRepo := mockDomain.NewMockMoviesRepository(ctrl)

			u := NewActorsUsecase(mockRepo, mockMovieRepo)

			tm, _ := time.Parse(time.DateOnly, tt.inputActor.BirthDate)

			actor := gormModels.Actor{
				Name:      tt.inputActor.Name,
				BirthDate: tm,
				Gender:    tt.inputActor.Gender,
			}

			tt.mockBehaviorCreateActor(mockRepo, actor)

			actorID, err := u.CreateActor(tt.inputActor)
			assert.Equal(t, tt.expectedActorID, actorID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_GetActorByID(t *testing.T) {
	type mockBehaviorGetActorByID func(r *mockDomain.MockActorsRepository, actorID uint64)

	tests := []struct {
		name                     string
		inputActorID             uint64
		mockBehaviorGetActorByID mockBehaviorGetActorByID
		expectedActorResponse    httpModels.ActorResponse
		expectedError            error
	}{
		{
			name:         "GetActorByID success",
			inputActorID: uint64(1),
			mockBehaviorGetActorByID: func(m *mockDomain.MockActorsRepository, actorID uint64) {
				m.EXPECT().
					GetActorByID(actorID).
					Return(gormModels.Actor{
						ID:        1,
						Name:      "Name",
						BirthDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
						Gender:    true,
					}, nil)
			},
			expectedActorResponse: httpModels.ActorResponse{
				ID:        1,
				Name:      "Name",
				BirthDate: "2006-01-02",
				Gender:    true,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockActorsRepository(ctrl)
			mockMovieRepo := mockDomain.NewMockMoviesRepository(ctrl)

			u := NewActorsUsecase(mockRepo, mockMovieRepo)

			tt.mockBehaviorGetActorByID(mockRepo, tt.inputActorID)

			actor, err := u.GetActorByID(tt.inputActorID)
			assert.Equal(t, tt.expectedActorResponse, actor)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_UpdateActor(t *testing.T) {
	type mockBehaviorUpdateActor func(r *mockDomain.MockActorsRepository, actor gormModels.Actor, actorID uint64)

	tests := []struct {
		name                    string
		inputActor              httpModels.Actor
		inputActorID            uint64
		mockBehaviorUpdateActor mockBehaviorUpdateActor
		expectedActorResponse   httpModels.ActorResponse
		expectedError           error
	}{
		{
			name: "UpdateActor success",
			inputActor: httpModels.Actor{
				Name:      "Name",
				BirthDate: "2006-01-02",
				Gender:    true,
			},
			inputActorID: uint64(1),
			mockBehaviorUpdateActor: func(m *mockDomain.MockActorsRepository, actor gormModels.Actor, actorID uint64) {
				m.EXPECT().
					UpdateActor(actor).
					Return(gormModels.Actor{
						ID:        1,
						Name:      "Name",
						BirthDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
						Gender:    true,
					}, nil)
			},
			expectedActorResponse: httpModels.ActorResponse{
				ID:        1,
				Name:      "Name",
				BirthDate: "2006-01-02",
				Gender:    true,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockActorsRepository(ctrl)
			mockMovieRepo := mockDomain.NewMockMoviesRepository(ctrl)

			u := NewActorsUsecase(mockRepo, mockMovieRepo)

			tm, _ := time.Parse(time.DateOnly, tt.inputActor.BirthDate)

			actor := gormModels.Actor{
				ID:        tt.inputActorID,
				Name:      tt.inputActor.Name,
				BirthDate: tm,
				Gender:    tt.inputActor.Gender,
			}

			tt.mockBehaviorUpdateActor(mockRepo, actor, tt.inputActorID)

			actorResponse, err := u.UpdateActor(tt.inputActor, tt.inputActorID)
			assert.Equal(t, tt.expectedActorResponse, actorResponse)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_GetActors(t *testing.T) {
	type mockBehaviorGetActors func(r *mockDomain.MockActorsRepository, pageNum uint64)
	type mockBehaviorGetMoviesByActorID func(r *mockDomain.MockMoviesRepository, actorID uint64)

	tests := []struct {
		name                         string
		inputPageNum                 uint64
		mockBehaviorGetActors        mockBehaviorGetActors
		mockBehaviorGetMoviesByActor mockBehaviorGetMoviesByActorID
		expectedActorResponse        []httpModels.GetActorsResponse
		expectedError                error
	}{
		{
			name:         "GetActors success",
			inputPageNum: uint64(1),
			mockBehaviorGetActors: func(m *mockDomain.MockActorsRepository, pageNum uint64) {
				m.EXPECT().
					GetActors(pageNum).
					Return([]gormModels.Actor{
						{
							ID:        1,
							Name:      "Name",
							BirthDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
							Gender:    true,
						},
					}, nil)
			},
			mockBehaviorGetMoviesByActor: func(m *mockDomain.MockMoviesRepository, actorID uint64) {
				m.EXPECT().
					GetMoviesOfActor(actorID).
					Return([]gormModels.Movie{
						{
							ID:          1,
							Title:       "Title",
							Description: "Description",
							ReleaseDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
						},
					}, nil)
			},
			expectedActorResponse: []httpModels.GetActorsResponse{
				{
					httpModels.ActorResponse{
						ID:        1,
						Name:      "Name",
						BirthDate: "2006-01-02",
						Gender:    true,
					},
					[]httpModels.MovieWithoutCastList{
						{
							ID:          1,
							Title:       "Title",
							ReleaseDate: "2006-01-02",
							Description: "Description",
						},
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockActorsRepository(ctrl)
			mockMovieRepo := mockDomain.NewMockMoviesRepository(ctrl)

			u := NewActorsUsecase(mockRepo, mockMovieRepo)

			tt.mockBehaviorGetActors(mockRepo, tt.inputPageNum)
			tt.mockBehaviorGetMoviesByActor(mockMovieRepo, uint64(1))

			actors, err := u.GetActors(tt.inputPageNum)
			assert.Equal(t, tt.expectedActorResponse, actors)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
