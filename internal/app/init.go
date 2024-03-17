package app

import (
	"log"
	"net/http"
	"os"

	httpActors "github.com/themilchenko/vk-tech_internship-problem_2024/internal/actors/delivery"
	httpAuth "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/delivery"
	authMiddleware "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/delivery/middleware"
	authRepository "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/repository"
	authUsecase "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/usecase"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpMovies "github.com/themilchenko/vk-tech_internship-problem_2024/internal/movies/delivery"
	password "github.com/themilchenko/vk-tech_internship-problem_2024/internal/utils/hash"
	logger "github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

type Server struct {
	Server *http.Server
	Config *config.Config
	Router http.Handler

	authUsecase   domain.AuthUsecase
	actorsUsecase domain.ActorsUsecase
	moviesUsecase domain.MoviesUsecase

	authHandler   httpAuth.AuthHandler
	actorsHandler httpActors.ActorsHandler
	moviesHandler httpMovies.ActorsHandler

	authMiddleware *authMiddleware.Middleware
	routes         Routes
}

func NewServer(s *http.Server, c *config.Config) *Server {
	return &Server{
		Server: s,
		Config: c,
	}
}

func (s *Server) Start() error {
	s.init()
	return s.Server.ListenAndServe()
}

func (s *Server) init() {
	s.Server.ErrorLog = log.New(os.Stdout, "SERVERLOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	s.makeUsecases()
	s.makeHandlers()
	s.makeMiddlewares()
	s.makeRouter()
}

func (s *Server) makeRouter() {
	s.routes = s.getRoutes()
	s.Router = logger.Logger(NewRouter(s.authMiddleware.LoginRequired, s.routes))
	s.Server.Handler = s.Router
}

func (s *Server) makeHandlers() {
	s.authHandler = httpAuth.NewAuthHandler(s.authUsecase, s.Config.CookieSettings)
}

func (s *Server) makeUsecases() {
	pgParams := s.Config.FormatDbAddr()

	authDB, err := authRepository.NewPostgres(pgParams)
	if err != nil {
		s.Server.ErrorLog.Println(err)
	}

	// moviesDB, err := moviesRepository.NewPostgres(pgParams)
	// if err != nil {
	// 	s.Server.ErrorLog.Println(err)
	// }

	s.authUsecase = authUsecase.NewAuthUsecase(
		authDB,
		s.Config.CookieSettings,
		password.HashPassword,
	)
}

func (s *Server) makeMiddlewares() {
	s.authMiddleware = authMiddleware.NewMiddleware(s.authUsecase)
}
