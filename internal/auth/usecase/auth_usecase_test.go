package authUsecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/domain"
	mockRepository "github.com/themilchenko/vk-tech_internship-problem_2024/internal/mocks/domain"
	gormModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	password "github.com/themilchenko/vk-tech_internship-problem_2024/internal/utils/hash"
	"go.uber.org/mock/gomock"
)

func TestUsecase_SignUp(t *testing.T) {
	type mockBehavior func(r *mockRepository.MockAuthRepository, user gormModels.User)
	type mockBehaviorSession func(r *mockRepository.MockAuthRepository, session gormModels.Session)

	tests := []struct {
		name                string
		inputUser           httpModels.AuthUser
		mockBehavior        mockBehavior
		mockBehaviorSession mockBehaviorSession
		expectedSessionID   string
		expectedUserID      uint64
		expectedError       error
	}{
		{
			name: "Successful signup",
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
				Role:     "admin",
			},
			mockBehavior: func(m *mockRepository.MockAuthRepository, user gormModels.User) {
				m.EXPECT().
					CreateUser(user).
					Return(uint64(1), nil)
			},
			mockBehaviorSession: func(m *mockRepository.MockAuthRepository, session gormModels.Session) {
				m.EXPECT().
					CreateSession(session).
					Return("cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2", nil)
			},
			expectedSessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			expectedUserID:    uint64(1),
			expectedError:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockAuthRepository(ctrl)

			u := NewCustomAuthUsecase(mockRepo, config.CookieSettings{},
				func(password string) (string, error) {
					if len(password) == 0 {
						return "", domain.ErrInternal
					}
					return password, nil
				},
				func(userID uint64, c config.CookieSettings) gormModels.Session {
					return gormModels.Session{
						UserID:    userID,
						SessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
					}
				},
			)

			user := gormModels.User{
				Username: tt.inputUser.Username,
				Password: tt.inputUser.Password,
				Role:     tt.inputUser.Role,
			}

			tt.mockBehavior(mockRepo, user)
			tt.mockBehaviorSession(mockRepo, u.cookieCreator(1, config.CookieSettings{}))

			sessionID, userID, err := u.SignUp(tt.inputUser)
			assert.Equal(t, tt.expectedSessionID, sessionID)
			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	type mockBehaviorGetUser func(r *mockRepository.MockAuthRepository, username string)
	type mockBehaviorSession func(r *mockRepository.MockAuthRepository, session gormModels.Session)

	tests := []struct {
		name                string
		inputUser           httpModels.AuthUser
		mockBehaviorGetUser mockBehaviorGetUser
		mockBehaviorSession mockBehaviorSession
		expectedSessionID   string
		expectedUserID      uint64
		expectedError       error
	}{
		{
			name: "Successful login",
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
			},
			mockBehaviorGetUser: func(m *mockRepository.MockAuthRepository, username string) {
				hashedPassword, _ := password.HashPassword("123") // Хешируем пароль
				m.EXPECT().
					GetUserByUsername(username).
					Return(gormModels.User{
						ID:       1,
						Username: "Jane",
						Password: hashedPassword, // Используем хеш пароля
						Role:     "admin",
					}, nil)
			},
			mockBehaviorSession: func(m *mockRepository.MockAuthRepository, session gormModels.Session) {
				m.EXPECT().
					CreateSession(session).
					Return("cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2", nil)
			},
			expectedSessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			expectedUserID:    uint64(1),
			expectedError:     nil,
		},
		{
			name: "Failed login",
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
			},
			mockBehaviorGetUser: func(m *mockRepository.MockAuthRepository, username string) {
				m.EXPECT().
					GetUserByUsername(username).
					Return(gormModels.User{}, errors.New("empty password"))
			},
			mockBehaviorSession: func(m *mockRepository.MockAuthRepository, session gormModels.Session) {},
			expectedSessionID:   "",
			expectedUserID:      uint64(0),
			expectedError:       errors.New("server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockAuthRepository(ctrl)

			u := NewCustomAuthUsecase(mockRepo, config.CookieSettings{},
				func(password string) (string, error) {
					if len(password) == 0 {
						return "", domain.ErrInternal
					}
					return password, nil
				},
				func(userID uint64, c config.CookieSettings) gormModels.Session {
					return gormModels.Session{
						UserID:    userID,
						SessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
					}
				},
			)

			tt.mockBehaviorGetUser(mockRepo, tt.inputUser.Username)
			tt.mockBehaviorSession(mockRepo, u.cookieCreator(1, config.CookieSettings{}))

			sessionID, userID, err := u.Login(tt.inputUser)
			assert.Equal(t, tt.expectedSessionID, sessionID)
			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUsecase_Auth(t *testing.T) {
	type mockBehaviorSession func(r *mockRepository.MockAuthRepository, session gormModels.Session)

	tests := []struct {
		name                string
		sessionID           string
		mockBehaviorSession mockBehaviorSession
		expectedUserID      uint64
		expectedError       error
	}{
		{
			name:      "Successful auth",
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehaviorSession: func(m *mockRepository.MockAuthRepository, session gormModels.Session) {
				m.EXPECT().
					GetUserBySessionID(session.SessionID).
					Return(gormModels.User{
						ID:       1,
						Username: "Jane",
						Password: "123",
						Role:     "admin",
					}, nil)
			},
			expectedUserID: uint64(1),
			expectedError:  nil,
		},
		{
			name:                "Failed auth",
			sessionID:           "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehaviorSession: func(m *mockRepository.MockAuthRepository, session gormModels.Session) {
				m.EXPECT().
					GetUserBySessionID(session.SessionID).
					Return(gormModels.User{}, errors.New("session not found"))
			},
			expectedUserID:      uint64(0),
			expectedError:       errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockAuthRepository(ctrl)

			u := NewCustomAuthUsecase(mockRepo, config.CookieSettings{},
				func(password string) (string, error) {
					if len(password) == 0 {
						return "", domain.ErrInternal
					}
					return password, nil
				},
				func(userID uint64, c config.CookieSettings) gormModels.Session {
					return gormModels.Session{
						UserID:    userID,
						SessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
					}
				},
			)

			tt.mockBehaviorSession(mockRepo, gormModels.Session{
				UserID:    1,
				SessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			})

			userID, err := u.Auth(tt.sessionID)
			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
