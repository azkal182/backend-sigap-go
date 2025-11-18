package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/service"
)

type jwtService struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService() service.TokenService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	accessExpiry := 15 * time.Minute
	if expiryStr := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"); expiryStr != "" {
		if parsed, err := time.ParseDuration(expiryStr); err == nil {
			accessExpiry = parsed
		}
	}

	refreshExpiry := 168 * time.Hour // 7 days
	if expiryStr := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"); expiryStr != "" {
		if parsed, err := time.ParseDuration(expiryStr); err == nil {
			refreshExpiry = parsed
		}
	}

	return &jwtService{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

// GenerateAccessToken generates a new access token
func (s *jwtService) GenerateAccessToken(userID uuid.UUID, username string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"roles":    roles,
		"type":     "access",
		"exp":      time.Now().Add(s.accessTokenExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken generates a new refresh token
func (s *jwtService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(s.refreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates and parses a JWT token
func (s *jwtService) ValidateToken(tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domainErrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domainErrors.ErrInvalidToken
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, domainErrors.ErrInvalidToken
	}

	if time.Now().Unix() > int64(exp) {
		return nil, domainErrors.ErrTokenExpired
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, domainErrors.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}

	tokenClaims := &service.TokenClaims{
		UserID: userID,
		Exp:    int64(exp),
	}

	// Extract username and roles for access tokens
	if username, ok := claims["username"].(string); ok {
		tokenClaims.Username = username
	}

	if roles, ok := claims["roles"].([]interface{}); ok {
		tokenClaims.Roles = make([]string, 0, len(roles))
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				tokenClaims.Roles = append(tokenClaims.Roles, roleStr)
			}
		}
	}

	return tokenClaims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *jwtService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Check if it's a refresh token
	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
		if tokenType, ok := tokenClaims["type"].(string); !ok || tokenType != "refresh" {
			return "", domainErrors.ErrInvalidToken
		}
	}

	// Generate new access token
	// Note: We need username and roles, but refresh token doesn't have them
	// So we'll need to fetch from database in the use case
	return s.GenerateAccessToken(claims.UserID, "", []string{})
}
