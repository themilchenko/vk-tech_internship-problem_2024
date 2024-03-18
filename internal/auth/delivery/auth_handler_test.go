package httpAuth

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
	mockDomain "github.com/themilchenko/vk-tech_internship-problem_2024/internal/mocks/domain"
	httpModels "github.com/themilchenko/vk-tech_internship-problem_2024/internal/models/http"
	"go.uber.org/mock/gomock"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockAuthUsecase, user httpModels.AuthUser)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            httpModels.AuthUser
		sessionID            string
		mockBehavior         mockBehavior
		cookieSettings       config.CookieSettings
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful signup",
			inputBody: `{"username":"Jane","password":"123","role":"admin"}`,
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
				Role:     "admin",
			},
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {
				m.EXPECT().
					SignUp(user).
					Return("cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2", uint64(1), nil)
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Failed signup",
			inputBody: `{}`,
			inputUser: httpModels.AuthUser{},
			mockBehavior: func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {
				m.EXPECT().
					SignUp(user).
					Return("", uint64(0), errors.New("empty password"))
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"empty password"}`,
		},
		{
			name:                 "Bad request",
			inputBody:            `sa;ldfkj`,
			inputUser:            httpModels.AuthUser{},
			mockBehavior:         func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid character 's' looking for beginning of value"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockAuthUsecase := mockDomain.NewMockAuthUsecase(cntx)

			tc.mockBehavior(mockAuthUsecase, tc.inputUser)

			handler := NewAuthHandler(mockAuthUsecase, tc.cookieSettings)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /signup", handler.Signup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/signup",
				bytes.NewBufferString(tc.inputBody),
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_Login(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockAuthUsecase, user httpModels.AuthUser)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            httpModels.AuthUser
		sessionID            string
		mockBehavior         mockBehavior
		cookieSettings       config.CookieSettings
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful login",
			inputBody: `{"username":"Jane","password":"123"}`,
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
			},
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {
				m.EXPECT().
					Login(user).
					Return("cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2", uint64(1), nil)
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Failed login",
			inputBody: `{"username":"Jane","password":"123"}`,
			inputUser: httpModels.AuthUser{
				Username: "Jane",
				Password: "123",
			},
			mockBehavior: func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {
				m.EXPECT().
					Login(user).
					Return("", uint64(0), errors.New("empty password"))
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"empty password"}`,
		},
		{
			name:                 "Bad request",
			inputBody:            `sa;ldfkj`,
			inputUser:            httpModels.AuthUser{},
			mockBehavior:         func(m *mockDomain.MockAuthUsecase, user httpModels.AuthUser) {},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid character 's' looking for beginning of value"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockAuthUsecase := mockDomain.NewMockAuthUsecase(cntx)

			tc.mockBehavior(mockAuthUsecase, tc.inputUser)

			handler := NewAuthHandler(mockAuthUsecase, tc.cookieSettings)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBufferString(tc.inputBody),
			)

			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockAuthUsecase, sessionID string)

	tests := []struct {
		name                 string
		sessionID            string
		mockBehavior         mockBehavior
		cookieSettings       config.CookieSettings
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful logout",
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, sessionID string) {
				m.EXPECT().
					Logout(sessionID).
					Return(nil)
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{}`,
		},
		{
			name:      "Unathorized",
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, sessionID string) {
				m.EXPECT().
					Logout(sessionID).
					Return(errors.New("session not found"))
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"session not found"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockAuthUsecase := mockDomain.NewMockAuthUsecase(cntx)

			tc.mockBehavior(mockAuthUsecase, tc.sessionID)

			handler := NewAuthHandler(mockAuthUsecase, tc.cookieSettings)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /logout", handler.Logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/logout",
				nil,
			)
			if len(tc.sessionID) > 0 {
				req.Header.Add("Cookie", "session_id="+tc.sessionID)
			}

			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandler_Auth(t *testing.T) {
	type mockBehavior func(r *mockDomain.MockAuthUsecase, sessionID string)

	tests := []struct {
		name                 string
		sessionID            string
		mockBehavior         mockBehavior
		cookieSettings       config.CookieSettings
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful auth",
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, sessionID string) {
				m.EXPECT().
					Auth(sessionID).
					Return(uint64(1), nil)
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Unathorized",
			sessionID: "cc982d23-e3e0-439e-a8b7-7fff4eeeb2f2",
			mockBehavior: func(m *mockDomain.MockAuthUsecase, sessionID string) {
				m.EXPECT().
					Auth(sessionID).
					Return(uint64(0), errors.New("session not found"))
			},
			cookieSettings:       config.CookieSettings{},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"session not found"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cntx := gomock.NewController(t)
			defer cntx.Finish()

			mockAuthUsecase := mockDomain.NewMockAuthUsecase(cntx)

			tc.mockBehavior(mockAuthUsecase, tc.sessionID)

			handler := NewAuthHandler(mockAuthUsecase, tc.cookieSettings)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /auth", handler.Auth)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/auth",
				nil,
			)
			if len(tc.sessionID) > 0 {
				req.Header.Add("Cookie", "session_id="+tc.sessionID)
			}

			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
