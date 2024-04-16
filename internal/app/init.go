package app

import (
	"net/http"

	httpActors "github.com/themilchenko/vk-tech_internship-problem_2024/internal/actors/delivery"
	actorsRepository "github.com/themilchenko/vk-tech_internship-problem_2024/internal/actors/repository"
	actorsUsecase "github.com/themilchenko/vk-tech_internship-problem_2024/internal/actors/usecase"
	httpAuth "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/delivery"
	authMiddleware "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/delivery/middleware"
	authRepository "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/repository"
	authUsecase "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/usecase"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpMovies "github.com/themilchenko/vk-tech_internship-problem_2024/internal/movies/delivery"
	moviesRepository "github.com/themilchenko/vk-tech_internship-problem_2024/internal/movies/repository"
	moviesUsecase "github.com/themilchenko/vk-tech_internship-problem_2024/internal/movies/usecase"
	password "github.com/themilchenko/vk-tech_internship-problem_2024/internal/utils/hash"
	logger "github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

const (
	baseURLPath = "/api/v1"
)

type Server struct {
	Server *http.Server
	Config *config.Config
	Router *http.ServeMux

	authUsecase   domain.AuthUsecase
	actorsUsecase domain.ActorsUsecase
	moviesUsecase domain.MoviesUsecase

	authHandler   httpAuth.AuthHandler
	actorsHandler httpActors.ActorsHandler
	moviesHandler httpMovies.ActorsHandler

	authMiddleware *authMiddleware.Middleware
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

func (s *Server) init() error {
	s.makeUsecases()
	s.makeHandlers()
	s.makeMiddlewares()
	s.makeRouter()

	return nil
}

func (s *Server) makeRouter() {
	s.Router = http.NewServeMux()
	http.Handle("/", logger.Middleware(s.Router))

	// authorization
	s.Router.HandleFunc("GET "+baseURLPath+"/auth", s.authHandler.Auth)
	s.Router.HandleFunc("POST "+baseURLPath+"/signup", s.authHandler.Signup)
	s.Router.HandleFunc("POST "+baseURLPath+"/login", s.authHandler.Login)
	s.Router.HandleFunc(
		"DELETE "+baseURLPath+"/logout",
		s.authMiddleware.LoginRequired(s.authHandler.Logout),
	)

	// actors
	s.Router.HandleFunc(
		"POST "+baseURLPath+"/actors",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.CreateActor),
		),
	)
	s.Router.HandleFunc(
		"GET "+baseURLPath+"/actors",
		s.authMiddleware.LoginRequired(s.actorsHandler.GetActors),
	)
	s.Router.HandleFunc(
		"PUT "+baseURLPath+"/actors/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.UpdateActor),
		),
	)
	s.Router.HandleFunc(
		"GET "+baseURLPath+"/actors/{id}",
		s.authMiddleware.LoginRequired(s.actorsHandler.GetActor),
	)
	s.Router.HandleFunc(
		"DELETE "+baseURLPath+"/actors/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.DeleteActor),
		),
	)

	// movies
	s.Router.HandleFunc(
		"POST "+baseURLPath+"/movies",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.CreateMovie),
		),
	)
	s.Router.HandleFunc(
		"GET "+baseURLPath+"/movies",
		s.authMiddleware.LoginRequired(s.moviesHandler.GetMovies),
	)
	s.Router.HandleFunc(
		"PUT "+baseURLPath+"/movies/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.UpdateMovie),
		),
	)
	s.Router.HandleFunc(
		"DELETE "+baseURLPath+"/movies/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.DeleteMovie),
		),
	)
	s.Router.HandleFunc(
		"GET "+baseURLPath+"/movies/{id}",
		s.authMiddleware.LoginRequired(s.moviesHandler.GetMovie),
	)
	s.Router.HandleFunc(
		"POST "+baseURLPath+"/movies/{movieID}/actors/{actorID}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.AddActorToMovie),
		),
	)
	s.Router.HandleFunc(
		"DELETE "+baseURLPath+"/movies/{movieID}/actors/{actorID}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.DeleteActorFromMoive),
		),
	)
}

func (s *Server) makeHandlers() {
	s.authHandler = httpAuth.NewAuthHandler(s.authUsecase, s.Config.CookieSettings)
	s.actorsHandler = httpActors.NewActorsUsecase(s.actorsUsecase)
	s.moviesHandler = httpMovies.NewActorsUsecase(s.moviesUsecase)
}

func (s *Server) makeUsecases() error {
	pgParams := s.Config.FormatDbAddr()

	authDB, err := authRepository.NewPostgres(pgParams)
	if err != nil {
		return err
	}

	actorsDB, err := actorsRepository.NewPostgres(pgParams, s.Config.PageSize)
	if err != nil {
		return err
	}

	moviesDB, err := moviesRepository.NewPostgres(pgParams)
	if err != nil {
		return err
	}

	s.authUsecase = authUsecase.NewAuthUsecase(
		authDB,
		s.Config.CookieSettings,
		password.HashPassword,
	)
	s.actorsUsecase = actorsUsecase.NewActorsUsecase(actorsDB, moviesDB)
	s.moviesUsecase = moviesUsecase.NewMoviesUsecase(moviesDB, actorsDB)

	return nil
}

func (s *Server) makeMiddlewares() {
	s.authMiddleware = authMiddleware.NewMiddleware(s.authUsecase)
}
