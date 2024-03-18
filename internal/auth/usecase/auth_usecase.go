package authUsecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	password "github.com/themilchenko/vk-tech_internship-problem_2024/internal/utils/hash"
	"gorm.io/gorm"
)

type (
	hashCreator   func(password string) (string, error)
	cookieCreator func(userID uint64, c config.CookieSettings) gormModels.Session
)

type AuthUsecase struct {
	authRepository domain.AuthRepository

	cookieSettings config.CookieSettings
	cookieCreator  cookieCreator
	hashCreator    hashCreator
}

func NewAuthUsecase(a domain.AuthRepository, c config.CookieSettings, h hashCreator) AuthUsecase {
	return AuthUsecase{
		authRepository: a,
		cookieSettings: c,
		cookieCreator:  generateCookie,
		hashCreator:    h,
	}
}

func NewCustomAuthUsecase(a domain.AuthRepository, c config.CookieSettings, h hashCreator, cc cookieCreator) AuthUsecase {
	return AuthUsecase{
		authRepository: a,
		cookieSettings: c,
		cookieCreator:  cc,
		hashCreator:    h,
	}
}

func generateCookie(userID uint64, c config.CookieSettings) gormModels.Session {
	return gormModels.Session{
		UserID:    userID,
		SessionID: uuid.New().String(),
		ExpireDate: time.Now().AddDate(
			int(c.ExpireDate.Years),
			int(c.ExpireDate.Months),
			int(c.ExpireDate.Days),
		),
	}
}

func (u AuthUsecase) SignUp(user httpModels.AuthUser) (string, uint64, error) {
	hash, err := u.hashCreator(user.Password)
	if err != nil {
		return "", 0, err
	}

	userID, err := u.authRepository.CreateUser(gormModels.User{
		Username: user.Username,
		Password: hash,
		Role:     user.Role,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", 0, domain.ErrUserAlreadyExist
		} else {
			return "", 0, err
		}
	}

	sessionID, err := u.authRepository.CreateSession(u.cookieCreator(userID, u.cookieSettings))
	if err != nil {
		return "", 0, err
	}

	return sessionID, userID, nil
}

func (u AuthUsecase) Login(user httpModels.AuthUser) (string, uint64, error) {
	recUser, err := u.authRepository.GetUserByUsername(user.Username)
	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			return "", 0, domain.ErrNotFound
		default:
			return "", 0, domain.ErrInternal
		}
	}

	matchPassword := password.CheckHashPassword(user.Password, recUser.Password)

	if !matchPassword {
		return "", 0, domain.ErrPasswordsNotEqual
	}

	sessionID, err := u.authRepository.CreateSession(u.cookieCreator(recUser.ID, u.cookieSettings))
	if err != nil {
		return "", 0, domain.ErrInternal
	}
	return sessionID, recUser.ID, nil
}

func (u AuthUsecase) Logout(sessionID string) error {
	return u.authRepository.DeleteBySessionID(sessionID)
}

func (u AuthUsecase) Auth(sessionID string) (uint64, error) {
	user, err := u.authRepository.GetUserBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.ErrNotFound
		}
		return 0, err
	}
	return user.ID, nil
}

func (u AuthUsecase) GetUserBySessionID(sessionID string) (httpModels.AuthUser, error) {
	recievedUser, err := u.authRepository.GetUserBySessionID(sessionID)
	if err != nil {
		return httpModels.AuthUser{}, err
	}
	return recievedUser.ToHTTPModel(), nil
}
