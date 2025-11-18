package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/domain/service"
)

// AuthUseCase handles authentication use cases
type AuthUseCase struct {
	userRepo     repository.UserRepository
	tokenService service.TokenService
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register handles user registration
func (uc *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, domainErrors.ErrUserAlreadyExists
	}

	// Create new user
	user := &entity.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Password:  req.Password,
		Name:      req.Name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Get user with roles for token generation
	userWithRoles, err := uc.userRepo.GetWithRoles(ctx, user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Generate tokens
	roles := make([]string, 0)
	for _, role := range userWithRoles.Roles {
		roles = append(roles, role.Name)
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID, user.Username, roles)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Username: user.Username,
			Name:     user.Name,
			Roles:    roles,
		},
	}, nil
}

// Login handles user login
func (uc *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by Username
	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// Check if user exists
	if user == nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domainErrors.ErrUserInactive
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// Get user with roles
	userWithRoles, err := uc.userRepo.GetWithRoles(ctx, user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Generate tokens
	roles := make([]string, 0)
	for _, role := range userWithRoles.Roles {
		roles = append(roles, role.Name)
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID, user.Username, roles)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Username: user.Username,
			Name:     user.Name,
			Roles:    roles,
		},
	}, nil
}

// RefreshToken handles token refresh
func (uc *AuthUseCase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}

	// Get user
	user, err := uc.userRepo.GetWithRoles(ctx, claims.UserID)
	if err != nil {
		return nil, domainErrors.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domainErrors.ErrUserInactive
	}

	// Generate new access token
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID, user.Username, roles)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Generate new refresh token
	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Username: user.Username,
			Name:     user.Name,
			Roles:    roles,
		},
	}, nil
}
