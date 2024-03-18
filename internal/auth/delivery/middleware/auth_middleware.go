package authMiddleware

import (
	"net/http"

	httpAuth "github.com/themilchenko/vk-tech_internship-problem_2024/internal/auth/delivery"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	"github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

type Middleware struct {
	authUsecase domain.AuthUsecase
}

func NewMiddleware(a domain.AuthUsecase) *Middleware {
	return &Middleware{
		authUsecase: a,
	}
}

func (m Middleware) LoginRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(httpAuth.CookieName)
		if err != nil {
			pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if _, err = m.authUsecase.Auth(cookie.Value); err != nil {
			pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (m Middleware) AccessRestriction(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(httpAuth.CookieName)
		if err != nil {
			pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := m.authUsecase.GetUserBySessionID(cookie.Value)
		if err != nil {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
			return
		}
		if user.Role != "admin" {
			pkg.HandleError(w, domain.ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
