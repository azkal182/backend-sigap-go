package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/service"
)

func TestAuthUseCase_Register(t *testing.T) {
	tests := []struct {
		name          string
		req           dto.RegisterRequest
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockTokenService)
		expectedError error
	}{
		{
			name: "success - register new user",
			req: dto.RegisterRequest{
				Username: "newuser",
				Password: "password123",
				Name:     "New User",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				// User doesn't exist
				userRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, domainErrors.ErrUserNotFound)

				// Create user
				userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username == "newuser" && u.Name == "New User"
				})).Return(nil)

				// Get user with roles (user ID akan di-generate di dalam use case)
				userRepo.On("GetWithRoles", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       uuid.New(),
					Username: "newuser",
					Name:     "New User",
					IsActive: true,
					Roles:    []entity.Role{},
				}, nil)

				// Generate tokens (tidak mengikat ke UUID tertentu)
				tokenService.On("GenerateAccessToken", mock.Anything, "newuser", []string{}).Return("access_token", nil)
				tokenService.On("GenerateRefreshToken", mock.Anything).Return("refresh_token", nil)
			},
			expectedError: nil,
		},
		{
			name: "failure - user already exists",
			req: dto.RegisterRequest{
				Username: "existing",
				Password: "password123",
				Name:     "Existing User",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				userRepo.On("GetByUsername", mock.Anything, "existing").Return(&entity.User{
					ID:       uuid.New(),
					Username: "existing",
				}, nil)
			},
			expectedError: domainErrors.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenService := new(mocks.MockTokenService)
			tt.setupMocks(userRepo, tokenService)

			authUseCase := NewAuthUseCase(userRepo, tokenService)
			resp, err := authUseCase.Register(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, tt.req.Username, resp.User.Username)
			}

			userRepo.AssertExpectations(t)
			tokenService.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	userID := uuid.New()

	// Create a user with properly hashed password for testing
	testUser := &entity.User{
		ID:       userID,
		Username: "user",
		Password: "password123",
		Name:     "Test User",
		IsActive: true,
	}
	require.NoError(t, testUser.HashPassword())
	hashedPassword := testUser.Password

	tests := []struct {
		name          string
		req           dto.LoginRequest
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockTokenService)
		expectedError error
	}{
		{
			name: "success - login with correct credentials",
			req: dto.LoginRequest{
				Username: "user",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				user := &entity.User{
					ID:       userID,
					Username: "user",
					Password: hashedPassword,
					Name:     "Test User",
					IsActive: true,
				}
				userRepo.On("GetByUsername", mock.Anything, "user").Return(user, nil)

				userWithRoles := &entity.User{
					ID:       userID,
					Username: "user",
					Password: hashedPassword,
					Name:     "Test User",
					IsActive: true,
					Roles:    []entity.Role{{ID: uuid.New(), Name: "user"}},
				}
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(userWithRoles, nil)

				tokenService.On("GenerateAccessToken", userID, "user", []string{"user"}).Return("access_token", nil)
				tokenService.On("GenerateRefreshToken", userID).Return("refresh_token", nil)
			},
			expectedError: nil,
		},
		{
			name: "failure - user not found",
			req: dto.LoginRequest{
				Username: "notfound",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				userRepo.On("GetByUsername", mock.Anything, "notfound").Return(nil, domainErrors.ErrUserNotFound)
			},
			expectedError: domainErrors.ErrInvalidCredentials,
		},
		{
			name: "failure - user inactive",
			req: dto.LoginRequest{
				Username: "inactive",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				user := &entity.User{
					ID:       userID,
					Username: "inactive",
					Password: hashedPassword,
					IsActive: false,
				}
				userRepo.On("GetByUsername", mock.Anything, "inactive").Return(user, nil)
			},
			expectedError: domainErrors.ErrUserInactive,
		},
		{
			name: "failure - wrong password",
			req: dto.LoginRequest{
				Username: "user",
				Password: "wrongpassword",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				user := &entity.User{
					ID:       userID,
					Username: "user",
					Password: hashedPassword,
					IsActive: true,
				}
				userRepo.On("GetByUsername", mock.Anything, "user").Return(user, nil)
			},
			expectedError: domainErrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenService := new(mocks.MockTokenService)
			tt.setupMocks(userRepo, tokenService)

			authUseCase := NewAuthUseCase(userRepo, tokenService)
			resp, err := authUseCase.Login(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
			}

			userRepo.AssertExpectations(t)
			tokenService.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_RefreshToken(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		req           dto.RefreshTokenRequest
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockTokenService)
		expectedError error
	}{
		{
			name: "success - refresh token",
			req: dto.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				claims := &service.TokenClaims{
					UserID: userID,
					Exp:    time.Now().Add(time.Hour).Unix(),
				}
				tokenService.On("ValidateToken", "valid_refresh_token").Return(claims, nil)

				userWithRoles := &entity.User{
					ID:       userID,
					Username: "user",
					Name:     "Test User",
					IsActive: true,
					Roles:    []entity.Role{{ID: uuid.New(), Name: "user"}},
				}
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(userWithRoles, nil)

				tokenService.On("GenerateAccessToken", userID, "user", []string{"user"}).Return("new_access_token", nil)
				tokenService.On("GenerateRefreshToken", userID).Return("new_refresh_token", nil)
			},
			expectedError: nil,
		},
		{
			name: "failure - invalid token",
			req: dto.RefreshTokenRequest{
				RefreshToken: "invalid_token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				tokenService.On("ValidateToken", "invalid_token").Return(nil, domainErrors.ErrInvalidToken)
			},
			expectedError: domainErrors.ErrInvalidToken,
		},
		{
			name: "failure - user not found",
			req: dto.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				claims := &service.TokenClaims{
					UserID: userID,
					Exp:    time.Now().Add(time.Hour).Unix(),
				}
				tokenService.On("ValidateToken", "valid_refresh_token").Return(claims, nil)
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(nil, domainErrors.ErrUserNotFound)
			},
			expectedError: domainErrors.ErrUserNotFound,
		},
		{
			name: "failure - user inactive",
			req: dto.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenService *mocks.MockTokenService) {
				claims := &service.TokenClaims{
					UserID: userID,
					Exp:    time.Now().Add(time.Hour).Unix(),
				}
				tokenService.On("ValidateToken", "valid_refresh_token").Return(claims, nil)

				userWithRoles := &entity.User{
					ID:       userID,
					Username: "user",
					IsActive: false,
				}
				userRepo.On("GetWithRoles", mock.Anything, userID).Return(userWithRoles, nil)
			},
			expectedError: domainErrors.ErrUserInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenService := new(mocks.MockTokenService)
			tt.setupMocks(userRepo, tokenService)

			authUseCase := NewAuthUseCase(userRepo, tokenService)
			resp, err := authUseCase.RefreshToken(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
			}

			userRepo.AssertExpectations(t)
			tokenService.AssertExpectations(t)
		})
	}
}
