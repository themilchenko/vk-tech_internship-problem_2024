package httpActors

import (
	"bytes"
	"errors"
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

func TestHandler_CreateActor(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockActorsUsecase, actor httpModels.Actor)

	tests := []struct {
		name                 string
		inputBody            string
		inputActor           httpModels.Actor
		actorID              uint64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful actor creation",
			inputBody: `{"name":"John","birthDate":"2000-01-01","gender":true}`,
			inputActor: httpModels.Actor{
				Name:      "John",
				BirthDate: "2000-01-01",
				Gender:    true,
			},
			actorID: uint64(1),
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor) {
				m.EXPECT().
					CreateActor(actor).
					Return(uint64(1), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Bad request",
			inputBody:            `sa;ldfkj`,
			inputActor:           httpModels.Actor{},
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid character 's' looking for beginning of value"}`,
		},
		{
			name:       "Internal server error",
			inputBody:  `{}`,
			inputActor: httpModels.Actor{},
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor) {
				m.EXPECT().
					CreateActor(actor).
					Return(uint64(0), errors.New("empty name"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"empty name"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockActorsUsecase := mockDomain.NewMockActorsUsecase(cntx)

			tt.mockBehavior(mockActorsUsecase, tt.inputActor)

			handler := NewActorsUsecase(mockActorsUsecase)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /actors", handler.CreateActor)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/actors",
				bytes.NewBufferString(tt.inputBody),
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_DeleteActor(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockActorsUsecase, actorID uint64)

	tests := []struct {
		name                 string
		actorID              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Successful actor deletion",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					DeleteActorByID(actorID).
					Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{}`,
		},
		{
			name:    "Internal server error",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					DeleteActorByID(actorID).
					Return(errors.New("empty actor"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"empty actor"}`,
		},
		{
			name:    "Status Not Found",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					DeleteActorByID(actorID).
					Return(domain.ErrNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"failed to find item"}`,
		},
		{
			name:                 "Bad request",
			actorID:              "a",
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, actorID uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"a\": invalid syntax"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockActorsUsecase := mockDomain.NewMockActorsUsecase(cntx)

			actorID, _ := strconv.ParseUint(tt.actorID, 10, 64)

			tt.mockBehavior(mockActorsUsecase, actorID)

			handler := NewActorsUsecase(mockActorsUsecase)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /actors/{id}", handler.DeleteActor)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/actors/"+tt.actorID,
				nil,
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_GetActor(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockActorsUsecase, actorID uint64)

	tests := []struct {
		name                 string
		actorID              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Successful actor getting",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					GetActorByID(actorID).
					Return(httpModels.ActorResponse{
						ID:        1,
						Name:      "John",
						Gender:    true,
						BirthDate: "2000-01-01",
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"name":"John","gender":true,"birthDate":"2000-01-01"}`,
		},
		{
			name:                 "Bad request",
			actorID:              "a",
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, actorID uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"a\": invalid syntax"}`,
		},
		{
			name:    "Error Not Found",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					GetActorByID(actorID).
					Return(httpModels.ActorResponse{}, domain.ErrNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"failed to find item"}`,
		},
		{
			name:    "Error Not Found",
			actorID: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actorID uint64) {
				m.EXPECT().
					GetActorByID(actorID).
					Return(httpModels.ActorResponse{}, domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockActorsUsecase := mockDomain.NewMockActorsUsecase(cntx)

			actorID, _ := strconv.ParseUint(tt.actorID, 10, 64)

			tt.mockBehavior(mockActorsUsecase, actorID)

			handler := NewActorsUsecase(mockActorsUsecase)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /actors/{id}", handler.GetActor)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/actors/"+tt.actorID,
				nil,
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_UpdateActor(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64)

	tests := []struct {
		name                 string
		actorID              string
		inputBody            string
		inputActor           httpModels.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful actor update",
			actorID:   "1",
			inputBody: `{"name":"John","birthDate":"2000-01-01","gender":true}`,
			inputActor: httpModels.Actor{
				Name:      "John",
				BirthDate: "2000-01-01",
				Gender:    true,
			},
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64) {
				m.EXPECT().
					UpdateActor(actor, actorID).
					Return(httpModels.ActorResponse{
						ID:        1,
						Name:      "John",
						Gender:    true,
						BirthDate: "2000-01-01",
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"name":"John","gender":true,"birthDate":"2000-01-01"}`,
		},
		{
			name:                 "Bad request",
			actorID:              "1",
			inputBody:            `sa;ldfkj`,
			inputActor:           httpModels.Actor{},
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid character 's' looking for beginning of value"}`,
		},
		{
			name:      "Error not found",
			actorID:   "1",
			inputBody: `{"name":"John","birthDate":"2000-01-01","gender":true}`,
			inputActor: httpModels.Actor{
				Name:      "John",
				BirthDate: "2000-01-01",
				Gender:    true,
			},
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64) {
				m.EXPECT().
					UpdateActor(actor, actorID).
					Return(httpModels.ActorResponse{}, domain.ErrNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"failed to find item"}`,
		},
		{
			name:                 "Bad Request",
			actorID:              "a",
			inputBody:            `{}`,
			inputActor:           httpModels.Actor{},
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"a\": invalid syntax"}`,
		},
		{
			name:      "Internal Error",
			actorID:   "1",
			inputBody: `{"name":"John","birthDate":"2000-01-01","gender":true}`,
			inputActor: httpModels.Actor{
				Name:      "John",
				BirthDate: "2000-01-01",
				Gender:    true,
			},
			mockBehavior: func(m *mockDomain.MockActorsUsecase, actor httpModels.Actor, actorID uint64) {
				m.EXPECT().
					UpdateActor(actor, actorID).
					Return(httpModels.ActorResponse{}, domain.ErrInternal)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockActorsUsecase := mockDomain.NewMockActorsUsecase(cntx)

			actorID, _ := strconv.ParseUint(tt.actorID, 10, 64)

			tt.mockBehavior(mockActorsUsecase, tt.inputActor, actorID)

			handler := NewActorsUsecase(mockActorsUsecase)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /actors/{id}", handler.UpdateActor)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut,
				"/actors/"+tt.actorID,
				bytes.NewBufferString(tt.inputBody),
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_GetActors(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockActorsUsecase, pageNum uint64)

	tests := []struct {
		name                 string
		pageNum              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Successful actors getting",
			pageNum: "",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, pageNum uint64) {
				m.EXPECT().
					GetActors(pageNum).
					Return([]httpModels.GetActorsResponse{
						{
							httpModels.ActorResponse{
								ID:        1,
								Name:      "John",
								Gender:    true,
								BirthDate: "2000-01-01",
							},
							[]httpModels.MovieWithoutCastList{},
						},
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"actor":{"id":1,"name":"John","gender":true,"birthDate":"2000-01-01"}}]`,
		},
		{
			name:                 "Error Bad Request",
			pageNum:              "a",
			mockBehavior:         func(m *mockDomain.MockActorsUsecase, pageNum uint64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.ParseUint: parsing \"a\": invalid syntax"}`,
		},
		{
			name:    "Internal server error",
			pageNum: "1",
			mockBehavior: func(m *mockDomain.MockActorsUsecase, pageNum uint64) {
				m.EXPECT().
					GetActors(pageNum).
					Return([]httpModels.GetActorsResponse{}, errors.New("empty actors"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"empty actors"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockActorsUsecase := mockDomain.NewMockActorsUsecase(cntx)

			pageNum := uint64(1)
			if len(tt.pageNum) != 0 {
				pageNum, _ = strconv.ParseUint(tt.pageNum, 10, 64)
			}

			tt.mockBehavior(mockActorsUsecase, pageNum)

			handler := NewActorsUsecase(mockActorsUsecase)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /actors", handler.GetActors)

			w := httptest.NewRecorder()
			var req *http.Request
			if tt.pageNum == "" {
				req = httptest.NewRequest(
					http.MethodGet,
					"/actors",
					nil,
				)
			} else {
				req = httptest.NewRequest(
					http.MethodGet,
					"/actors?page="+tt.pageNum,
					nil,
				)
			}

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
