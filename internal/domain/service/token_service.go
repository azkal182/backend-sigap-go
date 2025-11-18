package service

import (
	"time"

	"github.com/google/uuid"
)

// TokenService defines the interface for token operations
type TokenService interface {
	GenerateAccessToken(userID uuid.UUID, username string, roles []string) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   uuid.UUID
	Username string
	Roles    []string
	Exp      int64
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
