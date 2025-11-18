package service

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/testutil"
)

func TestJWTService_GenerateAccessToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()
	userID := uuid.New()
	username := "test"
	roles := []string{"admin", "user"}

	token, err := service.GenerateAccessToken(userID, username, roles)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()
	userID := uuid.New()

	token, err := service.GenerateRefreshToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTService_ValidateToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()
	userID := uuid.New()
	username := "test"
	roles := []string{"admin", "user"}

	// Generate token
	token, err := service.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Validate token
	claims, err := service.ValidateToken(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, roles, claims.Roles)
	assert.Greater(t, claims.Exp, time.Now().Unix())
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid token string",
			token: "invalid.token.string",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "malformed token",
			token: "not.a.valid.jwt.token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	// Set a very short expiry time
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRY", "1ms")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_ACCESS_TOKEN_EXPIRY")

	service := NewJWTService()
	userID := uuid.New()
	username := "test"
	roles := []string{"admin"}

	// Generate token
	token, err := service.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Validate token should fail
	claims, err := service.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret1")
	defer os.Unsetenv("JWT_SECRET")

	service1 := NewJWTService()
	userID := uuid.New()
	username := "test"
	roles := []string{"admin"}

	// Generate token with service1
	token, err := service1.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Try to validate with service2 using different secret
	os.Setenv("JWT_SECRET", "secret2")
	service2 := NewJWTService()

	claims, err := service2.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_RefreshAccessToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()
	userID := uuid.New()

	// Generate refresh token
	refreshToken, err := service.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Refresh access token
	accessToken, err := service.RefreshAccessToken(refreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}

func TestJWTService_RefreshAccessToken_InvalidToken(t *testing.T) {
	testutil.SetTestEnv()
	defer testutil.UnsetTestEnv()

	service := NewJWTService()

	// Try to refresh with access token (should fail)
	userID := uuid.New()
	username := "test"
	roles := []string{"admin"}

	accessToken, err := service.GenerateAccessToken(userID, username, roles)
	require.NoError(t, err)

	// Try to use access token as refresh token
	newAccessToken, err := service.RefreshAccessToken(accessToken)
	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
}

func TestJWTService_RefreshAccessToken_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRY", "1ms")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_REFRESH_TOKEN_EXPIRY")

	service := NewJWTService()
	userID := uuid.New()

	// Generate refresh token
	refreshToken, err := service.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to refresh should fail
	accessToken, err := service.RefreshAccessToken(refreshToken)
	assert.Error(t, err)
	assert.Empty(t, accessToken)
}
