// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/movies.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/movies.go -destination=internal/mocks/domain/movies.go
//

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	reflect "reflect"

	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	gomock "go.uber.org/mock/gomock"
)

// MockMoviesUsecase is a mock of MoviesUsecase interface.
type MockMoviesUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockMoviesUsecaseMockRecorder
}

// MockMoviesUsecaseMockRecorder is the mock recorder for MockMoviesUsecase.
type MockMoviesUsecaseMockRecorder struct {
	mock *MockMoviesUsecase
}

// NewMockMoviesUsecase creates a new mock instance.
func NewMockMoviesUsecase(ctrl *gomock.Controller) *MockMoviesUsecase {
	mock := &MockMoviesUsecase{ctrl: ctrl}
	mock.recorder = &MockMoviesUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMoviesUsecase) EXPECT() *MockMoviesUsecaseMockRecorder {
	return m.recorder
}

// AddActorFromMovie mocks base method.
func (m *MockMoviesUsecase) AddActorFromMovie(movieID, actorID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddActorFromMovie", movieID, actorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddActorFromMovie indicates an expected call of AddActorFromMovie.
func (mr *MockMoviesUsecaseMockRecorder) AddActorFromMovie(movieID, actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddActorFromMovie", reflect.TypeOf((*MockMoviesUsecase)(nil).AddActorFromMovie), movieID, actorID)
}

// CreateMovie mocks base method.
func (m *MockMoviesUsecase) CreateMovie(movie httpModels.MovieWithIDCast) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovie", movie)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMovie indicates an expected call of CreateMovie.
func (mr *MockMoviesUsecaseMockRecorder) CreateMovie(movie any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovie", reflect.TypeOf((*MockMoviesUsecase)(nil).CreateMovie), movie)
}

// DeleteActorFromMovie mocks base method.
func (m *MockMoviesUsecase) DeleteActorFromMovie(movieID, actorID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActorFromMovie", movieID, actorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActorFromMovie indicates an expected call of DeleteActorFromMovie.
func (mr *MockMoviesUsecaseMockRecorder) DeleteActorFromMovie(movieID, actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActorFromMovie", reflect.TypeOf((*MockMoviesUsecase)(nil).DeleteActorFromMovie), movieID, actorID)
}

// DeleteMovieByID mocks base method.
func (m *MockMoviesUsecase) DeleteMovieByID(movieID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMovieByID", movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMovieByID indicates an expected call of DeleteMovieByID.
func (mr *MockMoviesUsecaseMockRecorder) DeleteMovieByID(movieID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMovieByID", reflect.TypeOf((*MockMoviesUsecase)(nil).DeleteMovieByID), movieID)
}

// GetMovieByID mocks base method.
func (m *MockMoviesUsecase) GetMovieByID(movieID uint64) (httpModels.MovieResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieByID", movieID)
	ret0, _ := ret[0].(httpModels.MovieResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieByID indicates an expected call of GetMovieByID.
func (mr *MockMoviesUsecaseMockRecorder) GetMovieByID(movieID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieByID", reflect.TypeOf((*MockMoviesUsecase)(nil).GetMovieByID), movieID)
}

// GetMovies mocks base method.
func (m *MockMoviesUsecase) GetMovies(title, actorName string, sortBy httpModels.SortBy, order bool) ([]httpModels.MovieResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", title, actorName, sortBy, order)
	ret0, _ := ret[0].([]httpModels.MovieResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockMoviesUsecaseMockRecorder) GetMovies(title, actorName, sortBy, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockMoviesUsecase)(nil).GetMovies), title, actorName, sortBy, order)
}

// UpdateMovie mocks base method.
func (m *MockMoviesUsecase) UpdateMovie(movie httpModels.MovieWithoutCastList, movieID uint64) (httpModels.MovieResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMovie", movie, movieID)
	ret0, _ := ret[0].(httpModels.MovieResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMovie indicates an expected call of UpdateMovie.
func (mr *MockMoviesUsecaseMockRecorder) UpdateMovie(movie, movieID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMovie", reflect.TypeOf((*MockMoviesUsecase)(nil).UpdateMovie), movie, movieID)
}

// MockMoviesRepository is a mock of MoviesRepository interface.
type MockMoviesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMoviesRepositoryMockRecorder
}

// MockMoviesRepositoryMockRecorder is the mock recorder for MockMoviesRepository.
type MockMoviesRepositoryMockRecorder struct {
	mock *MockMoviesRepository
}

// NewMockMoviesRepository creates a new mock instance.
func NewMockMoviesRepository(ctrl *gomock.Controller) *MockMoviesRepository {
	mock := &MockMoviesRepository{ctrl: ctrl}
	mock.recorder = &MockMoviesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMoviesRepository) EXPECT() *MockMoviesRepositoryMockRecorder {
	return m.recorder
}

// AddActorToMovie mocks base method.
func (m *MockMoviesRepository) AddActorToMovie(movieID, actorID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddActorToMovie", movieID, actorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddActorToMovie indicates an expected call of AddActorToMovie.
func (mr *MockMoviesRepositoryMockRecorder) AddActorToMovie(movieID, actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddActorToMovie", reflect.TypeOf((*MockMoviesRepository)(nil).AddActorToMovie), movieID, actorID)
}

// CreateMovieWithCastList mocks base method.
func (m *MockMoviesRepository) CreateMovieWithCastList(movie gormModels.Movie, castList []uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovieWithCastList", movie, castList)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMovieWithCastList indicates an expected call of CreateMovieWithCastList.
func (mr *MockMoviesRepositoryMockRecorder) CreateMovieWithCastList(movie, castList any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovieWithCastList", reflect.TypeOf((*MockMoviesRepository)(nil).CreateMovieWithCastList), movie, castList)
}

// CreateMovieWithoutCastList mocks base method.
func (m *MockMoviesRepository) CreateMovieWithoutCastList(movie gormModels.Movie) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovieWithoutCastList", movie)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMovieWithoutCastList indicates an expected call of CreateMovieWithoutCastList.
func (mr *MockMoviesRepositoryMockRecorder) CreateMovieWithoutCastList(movie any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovieWithoutCastList", reflect.TypeOf((*MockMoviesRepository)(nil).CreateMovieWithoutCastList), movie)
}

// DeleteActorFromMovie mocks base method.
func (m *MockMoviesRepository) DeleteActorFromMovie(movieID, actorID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActorFromMovie", movieID, actorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActorFromMovie indicates an expected call of DeleteActorFromMovie.
func (mr *MockMoviesRepositoryMockRecorder) DeleteActorFromMovie(movieID, actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActorFromMovie", reflect.TypeOf((*MockMoviesRepository)(nil).DeleteActorFromMovie), movieID, actorID)
}

// DeleteMovieByID mocks base method.
func (m *MockMoviesRepository) DeleteMovieByID(movieID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMovieByID", movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMovieByID indicates an expected call of DeleteMovieByID.
func (mr *MockMoviesRepositoryMockRecorder) DeleteMovieByID(movieID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMovieByID", reflect.TypeOf((*MockMoviesRepository)(nil).DeleteMovieByID), movieID)
}

// GetMovieByID mocks base method.
func (m *MockMoviesRepository) GetMovieByID(movieID uint64) (gormModels.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieByID", movieID)
	ret0, _ := ret[0].(gormModels.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieByID indicates an expected call of GetMovieByID.
func (mr *MockMoviesRepositoryMockRecorder) GetMovieByID(movieID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieByID", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovieByID), movieID)
}

// GetMovies mocks base method.
func (m *MockMoviesRepository) GetMovies(title, actorName string, sortBy httpModels.SortBy, order bool) ([]gormModels.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", title, actorName, sortBy, order)
	ret0, _ := ret[0].([]gormModels.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockMoviesRepositoryMockRecorder) GetMovies(title, actorName, sortBy, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovies), title, actorName, sortBy, order)
}

// GetMoviesOfActor mocks base method.
func (m *MockMoviesRepository) GetMoviesOfActor(actorID uint64) ([]gormModels.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesOfActor", actorID)
	ret0, _ := ret[0].([]gormModels.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesOfActor indicates an expected call of GetMoviesOfActor.
func (mr *MockMoviesRepositoryMockRecorder) GetMoviesOfActor(actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesOfActor", reflect.TypeOf((*MockMoviesRepository)(nil).GetMoviesOfActor), actorID)
}

// UpdateMovie mocks base method.
func (m *MockMoviesRepository) UpdateMovie(movie gormModels.Movie) (gormModels.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMovie", movie)
	ret0, _ := ret[0].(gormModels.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMovie indicates an expected call of UpdateMovie.
func (mr *MockMoviesRepositoryMockRecorder) UpdateMovie(movie any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMovie", reflect.TypeOf((*MockMoviesRepository)(nil).UpdateMovie), movie)
}
