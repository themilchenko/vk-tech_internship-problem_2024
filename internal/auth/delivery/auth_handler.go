package httpAuth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"github.com/themilchenko/vk-tech_internship-problem_2024/pkg"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase

	cookieSettings config.CookieSettings
}

func NewAuthHandler(a domain.AuthUsecase, c config.CookieSettings) AuthHandler {
	return AuthHandler{
		authUsecase:    a,
		cookieSettings: c,
	}
}

func (h AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := h.authUsecase.Auth(cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}

	responseData, err := json.Marshal(httpModels.ID{ID: userID})
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var authUser httpModels.AuthUser
	err := json.NewDecoder(r.Body).Decode(&authUser)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, userID, err := h.authUsecase.Login(authUser)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			pkg.HandleError(w, err.Error(), http.StatusNotFound)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := h.makeHTTPCookie(session)
	http.SetCookie(w, cookie)

	responseData, err := json.Marshal(httpModels.ID{ID: userID})
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err = h.authUsecase.Logout(cookie.Value); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cookie.Expires = time.Now().AddDate(
		httpModels.DeleteExpire["year"],
		httpModels.DeleteExpire["month"],
		httpModels.DeleteExpire["day"],
	)
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(httpModels.EmptyModel)
}

func (h AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var receivedUser httpModels.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&receivedUser); err != nil {
		pkg.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, userID, err := h.authUsecase.SignUp(receivedUser)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExist) {
			pkg.HandleError(w, err.Error(), http.StatusConflict)
		} else {
			pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := h.makeHTTPCookie(sessionID)
	http.SetCookie(w, cookie)

	responseData, err := json.Marshal(httpModels.ID{ID: userID})
	if err != nil {
		pkg.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
