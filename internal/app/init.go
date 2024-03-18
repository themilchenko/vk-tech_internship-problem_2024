package app

import (
	"log"
	"net/http"
	"os"

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

func (s *Server) init() {
	s.Server.ErrorLog = log.New(os.Stdout, "SERVERLOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	s.makeUsecases()
	s.makeHandlers()
	s.makeMiddlewares()
	s.makeRouter()
}

func (s *Server) makeRouter() {
	s.Router = http.NewServeMux()
	http.Handle("/", logger.Middleware(s.Router))

	// authorization
	s.Router.HandleFunc("GET /auth", s.authHandler.Auth)
	s.Router.HandleFunc("POST /signup", s.authHandler.Signup)
	s.Router.HandleFunc("POST /login", s.authHandler.Login)
	s.Router.HandleFunc("DELETE /logout", s.authMiddleware.LoginRequired(s.authHandler.Logout))

	// actors
	s.Router.HandleFunc(
		"POST /actors",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.CreateActor),
		),
	)
	s.Router.HandleFunc("GET /actors", s.authMiddleware.LoginRequired(s.actorsHandler.GetActors))
	s.Router.HandleFunc(
		"PUT /actors/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.UpdateActor),
		),
	)
	s.Router.HandleFunc(
		"GET /actors/{id}",
		s.authMiddleware.LoginRequired(s.actorsHandler.GetActor),
	)
	s.Router.HandleFunc(
		"DELETE /actors/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.actorsHandler.DeleteActor),
		),
	)

	// movies
	s.Router.HandleFunc(
		"POST /movies",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.CreateMovie),
		),
	)
	s.Router.HandleFunc(
		"GET /movies",
		s.authMiddleware.LoginRequired(s.moviesHandler.GetMovies),
	)
	s.Router.HandleFunc(
		"PUT /movies/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.UpdateMovie),
		),
	)
	s.Router.HandleFunc(
		"DELETE /movies/{id}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.DeleteMovie),
		),
	)
	s.Router.HandleFunc(
		"GET /movies/{id}",
		s.authMiddleware.LoginRequired(s.moviesHandler.GetMovie),
	)
	s.Router.HandleFunc(
		"POST /movies/{movieID}/actors/{actorID}",
		s.authMiddleware.LoginRequired(
			s.authMiddleware.AccessRestriction(s.moviesHandler.AddActorToMovie),
		),
	)
	s.Router.HandleFunc(
		"DELETE /movies/{movieID}/actors/{actorID}",
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

func (s *Server) makeUsecases() {
	pgParams := s.Config.FormatDbAddr()

	authDB, err := authRepository.NewPostgres(pgParams)
	if err != nil {
		s.Server.ErrorLog.Println(err)
	}

	actorsDB, err := actorsRepository.NewPostgres(pgParams, s.Config.PageSize)
	if err != nil {
		s.Server.ErrorLog.Println(err)
	}

	moviesDB, err := moviesRepository.NewPostgres(pgParams)
	if err != nil {
		s.Server.ErrorLog.Println(err)
	}

	s.authUsecase = authUsecase.NewAuthUsecase(
		authDB,
		s.Config.CookieSettings,
		password.HashPassword,
	)
	s.actorsUsecase = actorsUsecase.NewActorsUsecase(actorsDB, moviesDB)
	s.moviesUsecase = moviesUsecase.NewMoviesUsecase(moviesDB, actorsDB)
}

func (s *Server) makeMiddlewares() {
	s.authMiddleware = authMiddleware.NewMiddleware(s.authUsecase)
}
