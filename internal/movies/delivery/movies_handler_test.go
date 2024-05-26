package httpMovies

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	mockDomain "github.com/themilchenko/vk-tech_internship-problem_2024/internal/mocks/domain"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreateMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithIDCast)

	tests := []struct {
		name                 string
		inputBody            string
		inputMovie           httpModels.MovieWithIDCast
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful movie creation",
			inputBody: `{"title":"The Godfather","description":"description","releaseDate":"1972-03-24","rating":9.2}`,
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "The Godfather",
				Description: "description",
				ReleaseDate: "1972-03-24",
				Rating:      9.2,
			},
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithIDCast) {
				m.EXPECT().
					CreateMovie(movie).
					Return(uint64(1), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Failed movie creation",
			inputBody: `{"title":"The Godfather","description":"description","releaseDate":"1972-03-24","rating":9.2}`,
			inputMovie: httpModels.MovieWithIDCast{
				Title:       "The Godfather",
				Description: "description",
				ReleaseDate: "1972-03-24",
				Rating:      9.2,
			},
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithIDCast) {
				m.EXPECT().
					CreateMovie(movie).
					Return(uint64(0), domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
		{
			name:                 "Bad request",
			inputBody:            `{"title":"The Godfather","description":description","release_date":"1972-03-24","rating":9.2}`,
			inputMovie:           httpModels.MovieWithIDCast{},
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithIDCast) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid character 'd' looking for beginning of value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			tt.mockBehavior(mockMoviesUsecase, tt.inputMovie)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /movies", handler.CreateMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/movies", strings.NewReader(tt.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_DeleteMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, id uint64)

	tests := []struct {
		name                 string
		inputID              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Successful movie deletion",
			inputID:              "1",
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, id uint64) { m.EXPECT().DeleteMovieByID(id).Return(nil) },
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{}",
		},
		{
			name:                 "Bad request",
			inputID:              "abs",
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, id uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"abs\": invalid syntax"}`,
		},
		{
			name:    "Err delete",
			inputID: "2",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, id uint64) {
				m.EXPECT().DeleteMovieByID(id).Return(domain.ErrNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"failed to find item"}`,
		},
		{
			name:    "Internal error",
			inputID: "2",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, id uint64) {
				m.EXPECT().DeleteMovieByID(id).Return(domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			id, _ := strconv.ParseUint(tt.inputID, 10, 64)
			tt.mockBehavior(mockMoviesUsecase, id)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /movies/{id}", handler.DeleteMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/movies/"+tt.inputID, nil)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_GetMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, id uint64)

	tests := []struct {
		name                 string
		inputID              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Successful movie get",
			inputID: "1",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, id uint64) {
				m.EXPECT().GetMovieByID(id).Return(httpModels.MovieResponse{
					ID:          1,
					Title:       "The Godfather",
					Description: "description",
					ReleaseDate: "1972-03-24",
					Rating:      9.2,
					CastList:    []httpModels.ActorResponse{},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"title\":\"The Godfather\",\"description\":\"description\",\"releaseDate\":\"1972-03-24\",\"rating\":9.2}",
		},
		{
			name:                 "Bad Request",
			inputID:              "a",
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, id uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"a\": invalid syntax"}`,
		},
		{
			name:    "Error Not Found",
			inputID: "1",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, id uint64) {
				m.EXPECT().GetMovieByID(id).Return(httpModels.MovieResponse{}, domain.ErrNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"failed to find item"}`,
		},
		{
			name:    "Intenal Error",
			inputID: "1",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, id uint64) {
				m.EXPECT().GetMovieByID(id).Return(httpModels.MovieResponse{}, domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			inputID, _ := strconv.ParseUint(tt.inputID, 10, 64)

			tt.mockBehavior(mockMoviesUsecase, inputID)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /movies/{id}", handler.GetMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/movies/"+tt.inputID, nil)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_GetMovies(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, title, actor string, filter httpModels.SortBy, isOrder bool)

	tests := []struct {
		name                 string
		title                string
		actor                string
		sortBy               string
		order                string
		inputPageNum         uint64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Successful movies get",
			inputPageNum: 1,
			title:        "godfather",
			actor:        "",
			sortBy:       "",
			order:        "",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, title, actor string, filter httpModels.SortBy, isOrder bool) {
				m.EXPECT().
					GetMovies(title, actor, filter, isOrder).
					Return([]httpModels.MovieResponse{
						{
							ID:          1,
							Title:       "The Godfather",
							Description: "description",
							ReleaseDate: "1972-03-24",
							Rating:      9.2,
							CastList:    []httpModels.ActorResponse{},
						},
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"title":"The Godfather","description":"description","releaseDate":"1972-03-24","rating":9.2}]`,
		},
		{
			name:                 "Error Bad Request Order",
			inputPageNum:         1,
			title:                "",
			actor:                "",
			sortBy:               "",
			order:                "a",
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, title, actor string, filter httpModels.SortBy, isOrder bool) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseBool: parsing \"a\": invalid syntax"}`,
		},
		{
			name:                 "Error Bad Request Filter",
			inputPageNum:         1,
			title:                "",
			actor:                "",
			sortBy:               "a",
			order:                "",
			mockBehavior:         func(m *mockDomain.MockMoviesUsecase, title, actor string, filter httpModels.SortBy, isOrder bool) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad request"}`,
		},
		{
			name:         "Internal Error",
			inputPageNum: 1,
			title:        "",
			actor:        "",
			sortBy:       "",
			order:        "",
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, title, actor string, filter httpModels.SortBy, isOrder bool) {
				m.EXPECT().
					GetMovies(title, actor, filter, isOrder).
					Return([]httpModels.MovieResponse{}, domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			order := true
			if len(tt.order) != 0 {
				order, _ = strconv.ParseBool(tt.order)
			}
			sortBy := "rating"
			if len(tt.sortBy) != 0 {
				sortBy = tt.sortBy
			}

			tt.mockBehavior(
				mockMoviesUsecase,
				tt.title,
				tt.actor,
				httpModels.SortBy(sortBy),
				order,
			)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /movies", handler.GetMovies)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/movies", nil)

			v := req.URL.Query()
			v.Add("title", fmt.Sprint(tt.title))
			v.Add("actor", fmt.Sprint(tt.actor))
			v.Add("order", fmt.Sprint(tt.order))
			v.Add("filter", fmt.Sprint(tt.sortBy))
			req.URL.RawQuery = v.Encode()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_UpdateMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithoutCastList, id uint64)

	tests := []struct {
		name                 string
		inputID              uint64
		inputBody            string
		inputMovie           httpModels.MovieWithoutCastList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful movie update",
			inputID:   1,
			inputBody: `{"title":"The Godfather","description":"description","releaseDate":"1972-03-24","rating":9.2}`,
			inputMovie: httpModels.MovieWithoutCastList{
				Title:       "The Godfather",
				Description: "description",
				ReleaseDate: "1972-03-24",
				Rating:      9.2,
			},
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, movie httpModels.MovieWithoutCastList, id uint64) {
				m.EXPECT().
					UpdateMovie(movie, id).
					Return(httpModels.MovieResponse{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			tt.mockBehavior(mockMoviesUsecase, tt.inputMovie, tt.inputID)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /movies/{id}", handler.UpdateMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/movies/1", strings.NewReader(tt.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_AddActorToMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, movieID, actorID uint64)

	tests := []struct {
		name                 string
		inputMovieID         uint64
		inputActorID         uint64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Successful actor addition to movie",
			inputMovieID: 1,
			inputActorID: 1,
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, movieID, actorID uint64) {
				m.EXPECT().AddActorFromMovie(movieID, actorID).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			tt.mockBehavior(mockMoviesUsecase, tt.inputMovieID, tt.inputActorID)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /movies/{movieID}/actors/{actorID}", handler.AddActorToMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/movies/1/actors/1", nil)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_DeleteActorFromMovie(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockMoviesUsecase, movieID, actorID uint64)

	tests := []struct {
		name                 string
		inputMovieID         uint64
		inputActorID         uint64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Successful actor deletion from movie",
			inputMovieID: 1,
			inputActorID: 1,
			mockBehavior: func(m *mockDomain.MockMoviesUsecase, movieID, actorID uint64) {
				m.EXPECT().DeleteActorFromMovie(movieID, actorID).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockMoviesUsecase := mockDomain.NewMockMoviesUsecase(cntx)
			handler := NewActorsUsecase(mockMoviesUsecase)

			tt.mockBehavior(mockMoviesUsecase, tt.inputMovieID, tt.inputActorID)

			mux := http.NewServeMux()
			mux.HandleFunc(
				"DELETE /movies/{movieID}/actors/{actorID}",
				handler.DeleteActorFromMoive,
			)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/movies/1/actors/1", nil)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
