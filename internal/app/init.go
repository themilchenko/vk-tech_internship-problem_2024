package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
	httpActors "github.com/themilchenko/vk-tech_internship-problem_2024/internal/actors/delivery"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpMovies "github.com/themilchenko/vk-tech_internship-problem_2024/internal/movies/delivery"
)

type Server struct {
	Server *http.Server
	Config *config.Config
	Router *http.ServeMux
	Logger *logrus.Logger

	actorsUsecase domain.ActorsUsecase
	moviesUsecase domain.MoviesUsecase

	actorsHandler httpActors.ActorsHandler
	moviesHandler httpMovies.ActorsHandler
}

func NewServer(s *http.Server, c *config.Config) *Server {
	return &Server{
		Server: s,
		Config: c,
		Router: http.NewServeMux(),
	}
}

func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}
	return s.Server.ListenAndServe()
}

func (s *Server) init() error {
	if err := s.makeLogger(); err != nil {
		return err
	}

	return nil
}

func (s *Server) makeHandlers() {
}

func (s *Server) makeUsecases() {
}

func (s *Server) makeRouter() {
	s.Router.HandleFunc("GET /hello")
}

func (s *Server) makeMiddlewares() {
}

func (s *Server) makeLogger() error {
	lvl, err := logrus.ParseLevel(s.Config.LoggerLvl)
	if err != nil {
		return err
	}
	s.Logger.SetLevel(lvl)
	return nil
}

func (s *Server) makeCORS() {
}
