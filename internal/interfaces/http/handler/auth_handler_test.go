package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/handler/mocks"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockAuthUseCase)
		expectedStatus int
	}{
		{
			name: "success - register new user",
			requestBody: dto.RegisterRequest{
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("Register", mock.Anything, mock.MatchedBy(func(req dto.RegisterRequest) bool {
					return req.Email == "newuser@example.com" && req.Name == "New User"
				})).Return(&dto.AuthResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					User: dto.UserDTO{
						ID:    "user-id",
						Email: "newuser@example.com",
						Name:  "New User",
					},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "failure - invalid request body",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				// No mock call expected for invalid request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failure - user already exists",
			requestBody: dto.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Existing User",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("Register", mock.Anything, mock.Anything).Return(nil, domainErrors.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(mocks.MockAuthUseCase)
			tt.setupMocks(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := setupRouter()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockAuthUseCase)
		expectedStatus int
	}{
		{
			name: "success - login with valid credentials",
			requestBody: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("Login", mock.Anything, mock.MatchedBy(func(req dto.LoginRequest) bool {
					return req.Email == "user@example.com"
				})).Return(&dto.AuthResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					User: dto.UserDTO{
						ID:    "user-id",
						Email: "user@example.com",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure - invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("Login", mock.Anything, mock.Anything).Return(nil, domainErrors.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - invalid request body",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				// No mock call expected
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(mocks.MockAuthUseCase)
			tt.setupMocks(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := setupRouter()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockAuthUseCase)
		expectedStatus int
	}{
		{
			name: "success - refresh token",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req dto.RefreshTokenRequest) bool {
					return req.RefreshToken == "valid_refresh_token"
				})).Return(&dto.AuthResponse{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
					User: dto.UserDTO{
						ID: "user-id",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure - invalid token",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "invalid_token",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				mockUseCase.On("RefreshToken", mock.Anything, mock.Anything).Return(nil, domainErrors.ErrInvalidToken)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - missing refresh token",
			requestBody: map[string]interface{}{
				"refresh_token": "",
			},
			setupMocks: func(mockUseCase *mocks.MockAuthUseCase) {
				// No mock call expected
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(mocks.MockAuthUseCase)
			tt.setupMocks(mockUseCase)

			handler := NewAuthHandler(mockUseCase)

			router := setupRouter()
			router.POST("/refresh", handler.RefreshToken)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}
