package moviesUsecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	mockDomain "github.com/themilchenko/vk-tech_internship-problem_2024/internal/mocks/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"go.uber.org/mock/gomock"
)

func TestUsecase_CreateMovie(t *testing.T) {
	type mockBehaviorCreateMovie func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64)
	type mockBehaviorCreateMovieWithoutCastList func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie)
	type mockBehaviorGetActor func(m *mockDomain.MockActorsRepository, actorID uint64)

	tests := []struct {
		name                                   string
		inputMovie                             httpModels.MovieWithIDCast
		mockBehaviorCreateMovie                mockBehaviorCreateMovie
		mockBehaviorCreateMovieWithoutCastList mockBehaviorCreateMovieWithoutCastList
		mockBehaviorGetActor                   mockBehaviorGetActor
		expectedMovieID                        uint64
		expectedError                          error
	}{
		{
			name: "CreateMovie success with cast list",
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
				CastIDList:  []uint64{1},
			},
			mockBehaviorCreateMovie: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64) {
				m.EXPECT().
					CreateMovieWithCastList(movie, castIDList).
					Return(uint64(1), nil)
			},
			mockBehaviorGetActor: func(m *mockDomain.MockActorsRepository, actorID uint64) {
				m.EXPECT().GetActorByID(actorID).Return(gormModels.Actor{}, nil)
			},
			mockBehaviorCreateMovieWithoutCastList: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie) {},
			expectedMovieID:                        uint64(1),
			expectedError:                          nil,
		},
		{
			name: "CreateMovie success without cast list",
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
				CastIDList:  []uint64{},
			},
			mockBehaviorCreateMovie: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64) {},
			mockBehaviorCreateMovieWithoutCastList: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie) {
				m.EXPECT().CreateMovieWithoutCastList(movie).Return(uint64(1), nil)
			},
			mockBehaviorGetActor: func(m *mockDomain.MockActorsRepository, actorID uint64) {},
			expectedMovieID:      uint64(1),
			expectedError:        nil,
		},
		{
			name: "CreateMovie error with cast list",
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
				CastIDList:  []uint64{1},
			},
			mockBehaviorCreateMovie: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64) {
				m.EXPECT().
					CreateMovieWithCastList(movie, castIDList).
					Return(uint64(0), errors.New("failed to create item"))
			},
			mockBehaviorCreateMovieWithoutCastList: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie) {},
			mockBehaviorGetActor: func(m *mockDomain.MockActorsRepository, actorID uint64) {
				m.EXPECT().GetActorByID(actorID).Return(gormModels.Actor{}, nil)
			},
			expectedMovieID: uint64(0),
			expectedError:   domain.ErrCreate,
		},
		{
			name: "CreateMovie error with cast list",
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
				CastIDList:  []uint64{},
			},
			mockBehaviorCreateMovie: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64) {},
			mockBehaviorCreateMovieWithoutCastList: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie) {
				m.EXPECT().
					CreateMovieWithoutCastList(movie).
					Return(uint64(0), errors.New("failed to create item"))
			},
			mockBehaviorGetActor: func(m *mockDomain.MockActorsRepository, actorID uint64) {},
			expectedMovieID:      uint64(0),
			expectedError:        domain.ErrCreate,
		},
		{
			name: "CreateMovie error parse date",
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006=01-02",
				Rating:      5.0,
				CastIDList:  []uint64{},
			},
			mockBehaviorCreateMovie:                func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, castIDList []uint64) {},
			mockBehaviorCreateMovieWithoutCastList: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie) {},
			mockBehaviorGetActor:                   func(m *mockDomain.MockActorsRepository, actorID uint64) {},
			expectedMovieID:                        uint64(0),
			expectedError: &time.ParseError{
				Layout:     "2006-01-02",
				Value:      "2006=01-02",
				LayoutElem: "-",
				ValueElem:  "=01-02",
				Message:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockMoviesRepository(ctrl)
			mockActorRepo := mockDomain.NewMockActorsRepository(ctrl)

			u := NewMoviesUsecase(mockRepo, mockActorRepo)

			tm, _ := time.Parse(time.DateOnly, tt.inputMovie.ReleaseDate)

			movie := gormModels.Movie{
				Title:       tt.inputMovie.Title,
				Description: tt.inputMovie.Description,
				ReleaseDate: tm,
				Rating:      tt.inputMovie.Rating,
			}

			tt.mockBehaviorCreateMovie(mockRepo, movie, tt.inputMovie.CastIDList)
			tt.mockBehaviorCreateMovieWithoutCastList(mockRepo, movie)
			if len(tt.inputMovie.CastIDList) > 0 {
				tt.mockBehaviorGetActor(mockActorRepo, tt.inputMovie.CastIDList[0])
			}

			movieID, err := u.CreateMovie(tt.inputMovie)
			assert.Equal(t, tt.expectedMovieID, movieID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_UpdateMovie(t *testing.T) {
	type mockBehaviorUpdateMovie func(r *mockDomain.MockMoviesRepository, movie gormModels.Movie, movieID uint64)

	tests := []struct {
		name                    string
		inputMovie              httpModels.MovieWithoutCastList
		inputMovieID            uint64
		mockBehaviorUpdateMovie mockBehaviorUpdateMovie
		expectedMovieResponse   httpModels.MovieResponse
		expectedError           error
	}{
		{
			name: "UpdateMovie success",
			inputMovie: httpModels.MovieWithoutCastList{
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
			},
			inputMovieID: uint64(1),
			mockBehaviorUpdateMovie: func(m *mockDomain.MockMoviesRepository, movie gormModels.Movie, movieID uint64) {
				m.EXPECT().
					UpdateMovie(movie).
					Return(gormModels.Movie{
						ID:          1,
						Title:       "Title",
						Description: "Description",
						ReleaseDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
						Rating:      5.0,
					}, nil)
			},
			expectedMovieResponse: httpModels.MovieResponse{
				ID:          1,
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockDomain.NewMockMoviesRepository(ctrl)
			mockActorRepo := mockDomain.NewMockActorsRepository(ctrl)

			u := NewMoviesUsecase(mockRepo, mockActorRepo)

			tm, _ := time.Parse(time.DateOnly, tt.inputMovie.ReleaseDate)

			movie := gormModels.Movie{
				ID:          tt.inputMovieID,
				Title:       tt.inputMovie.Title,
				Description: tt.inputMovie.Description,
				ReleaseDate: tm,
				Rating:      tt.inputMovie.Rating,
			}

			tt.mockBehaviorUpdateMovie(mockRepo, movie, tt.inputMovieID)

			movieResponse, err := u.UpdateMovie(tt.inputMovie, tt.inputMovieID)
			assert.Equal(t, tt.expectedMovieResponse, movieResponse)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_GetMovieByID(t *testing.T) {
	type mockBehaviorGetMovieByID func(r *mockDomain.MockMoviesRepository, movieID uint64)
	type mockBehaviorGetActorsFromMovie func(r *mockDomain.MockActorsRepository, movieID uint64)

	tests := []struct {
		name                           string
		inputMovieID                   uint64
		mockBehaviorGetMovieByID       mockBehaviorGetMovieByID
		mockBehaviorGetActorsFromMovie mockBehaviorGetActorsFromMovie
		expectedMovieResponse          httpModels.MovieResponse
		expectedError                  error
	}{
		{
			name:         "GetMovieByID success",
			inputMovieID: uint64(1),
			mockBehaviorGetMovieByID: func(m *mockDomain.MockMoviesRepository, movieID uint64) {
				m.EXPECT().
					GetMovieByID(movieID).
					Return(gormModels.Movie{
						ID:          1,
						Title:       "Title",
						Description: "Description",
						ReleaseDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
						Rating:      5.0,
					}, nil)
			},
			mockBehaviorGetActorsFromMovie: func(m *mockDomain.MockActorsRepository, movieID uint64) {
				m.EXPECT().
					GetActorsFromMovie(movieID).
					Return([]gormModels.Actor{
						{
							ID:        1,
							Name:      "Name",
							BirthDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
							Gender:    true,
						},
					}, nil)
			},
			expectedMovieResponse: httpModels.MovieResponse{
				ID:          1,
				Title:       "Title",
				Description: "Description",
				ReleaseDate: "2006-01-02",
				Rating:      5.0,
				CastList: []httpModels.ActorResponse{
					{
						ID:        1,
						Name:      "Name",
						Gender:    true,
						BirthDate: "2006-01-02",
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

			mockRepo := mockDomain.NewMockMoviesRepository(ctrl)
			mockActorRepo := mockDomain.NewMockActorsRepository(ctrl)

			u := NewMoviesUsecase(mockRepo, mockActorRepo)

			tt.mockBehaviorGetMovieByID(mockRepo, tt.inputMovieID)
			tt.mockBehaviorGetActorsFromMovie(mockActorRepo, tt.inputMovieID)

			movieResponse, err := u.GetMovieByID(tt.inputMovieID)
			assert.Equal(t, tt.expectedMovieResponse, movieResponse)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_GetMovies(t *testing.T) {
	type mockBehaviorGetMovies func(r *mockDomain.MockMoviesRepository, title, actor string, sortBy httpModels.SortBy, order bool)
	type mockBehaviorGetActorsFromMovie func(r *mockDomain.MockActorsRepository, movieID uint64)

	tests := []struct {
		name                           string
		title                          string
		actorName                      string
		sortBy                         httpModels.SortBy
		order                          bool
		mockBehaviorGetMovies          mockBehaviorGetMovies
		mockBehaviorGetActorsFromMovie mockBehaviorGetActorsFromMovie
		expectedMovieList              []httpModels.MovieResponse
		expectedError                  error
	}{
		{
			name:      "GetMovies success",
			title:     "Title",
			actorName: "Name",
			sortBy:    "title",
			order:     true,
			mockBehaviorGetMovies: func(m *mockDomain.MockMoviesRepository, title, actor string, sortBy httpModels.SortBy, order bool) {
				m.EXPECT().
					GetMovies(title, actor, sortBy, order).
					Return([]gormModels.Movie{
						{
							ID:          1,
							Title:       "Title",
							Description: "Description",
							ReleaseDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
							Rating:      5.0,
						},
					}, nil)
			},
			mockBehaviorGetActorsFromMovie: func(m *mockDomain.MockActorsRepository, movieID uint64) {
				m.EXPECT().
					GetActorsFromMovie(movieID).
					Return([]gormModels.Actor{
						{
							ID:        1,
							Name:      "Name",
							BirthDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
							Gender:    true,
						},
					}, nil)
			},
			expectedMovieList: []httpModels.MovieResponse{
				{
					ID:          1,
					Title:       "Title",
					Description: "Description",
					ReleaseDate: "2006-01-02",
					Rating:      5.0,
					CastList: []httpModels.ActorResponse{
						{
							ID:        1,
							Name:      "Name",
							Gender:    true,
							BirthDate: "2006-01-02",
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

			mockRepo := mockDomain.NewMockMoviesRepository(ctrl)
			mockActorRepo := mockDomain.NewMockActorsRepository(ctrl)

			u := NewMoviesUsecase(mockRepo, mockActorRepo)

			tt.mockBehaviorGetMovies(mockRepo, tt.title, tt.actorName, tt.sortBy, tt.order)
			tt.mockBehaviorGetActorsFromMovie(mockActorRepo, uint64(1))

			movieList, err := u.GetMovies(tt.title, tt.actorName, tt.sortBy, tt.order)
			assert.Equal(t, tt.expectedMovieList, movieList)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
