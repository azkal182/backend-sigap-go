package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
)

// MockAuthUseCase is a mock implementation of AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *MockAuthUseCase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}
