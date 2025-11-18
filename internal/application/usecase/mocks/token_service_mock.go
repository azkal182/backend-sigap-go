package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/service"
)

// MockTokenService is a mock implementation of TokenService
type MockTokenService struct {
	mock.Mock
}

// Ensure MockTokenService implements service.TokenService
var _ service.TokenService = (*MockTokenService)(nil)

func (m *MockTokenService) GenerateAccessToken(userID uuid.UUID, username string, roles []string) (string, error) {
	args := m.Called(userID, username, roles)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(tokenString string) (*service.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenClaims), args.Error(1)
}

func (m *MockTokenService) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}
