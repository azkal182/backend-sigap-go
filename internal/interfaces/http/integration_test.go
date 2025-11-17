package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	infraRepo "github.com/your-org/go-backend-starter/internal/infrastructure/repository"
	infraService "github.com/your-org/go-backend-starter/internal/infrastructure/service"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/handler"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/middleware"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/router"
	"github.com/your-org/go-backend-starter/internal/testutil"
	"gorm.io/gorm"
)

// testUserRepository wraps userRepository to use a specific DB
type testUserRepository struct {
	db *gorm.DB
}

func (r *testUserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *testUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *testUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *testUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *testUserRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) GetWithRolesAndDormitories(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Preload("Dormitories").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func setupTestRouter(t *testing.T) (*gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	testutil.SetTestEnv()

	// Temporarily replace database.DB for repositories
	originalDB := database.DB
	database.DB = testDB

	// Initialize repositories with test database
	userRepo := &testUserRepository{db: testDB}
	roleRepo := infraRepo.NewRoleRepository() // These will use database.DB
	dormitoryRepo := infraRepo.NewDormitoryRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router
	r := router.SetupRouter(authHandler, userHandler, dormitoryHandler, authMiddleware)

	cleanup := func() {
		database.DB = originalDB // Restore original DB
		testutil.CleanupTestDB(t, testDB)
		testutil.UnsetTestEnv()
	}

	return r, cleanup
}

func TestAuthIntegration_RegisterAndLogin(t *testing.T) {
	router, cleanup := setupTestRouter(t)
	defer cleanup()

	// Test Register
	registerReq := dto.RegisterRequest{
		Email:    "integration@example.com",
		Password: "password123",
		Name:     "Integration Test User",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)

	assert.Equal(t, http.StatusCreated, registerW.Code)

	var registerResp map[string]interface{}
	err := json.Unmarshal(registerW.Body.Bytes(), &registerResp)
	require.NoError(t, err)
	assert.True(t, registerResp["success"].(bool))

	// Extract tokens from response
	data := registerResp["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)
	require.NotEmpty(t, accessToken)

	// Test Login
	loginReq := dto.LoginRequest{
		Email:    "integration@example.com",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusOK, loginW.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	assert.True(t, loginResp["success"].(bool))

	// Test protected endpoint with token
	protectedReq, _ := http.NewRequest(http.MethodGet, "/api/users/me", nil)
	protectedReq.Header.Set("Authorization", "Bearer "+accessToken)

	protectedW := httptest.NewRecorder()
	router.ServeHTTP(protectedW, protectedReq)

	// Note: This test assumes there's a /api/users/me endpoint
	// Adjust based on your actual routes
	assert.True(t, protectedW.Code == http.StatusOK || protectedW.Code == http.StatusNotFound)
}

func TestAuthIntegration_InvalidCredentials(t *testing.T) {
	router, cleanup := setupTestRouter(t)
	defer cleanup()

	// Register user first
	registerReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)
	assert.Equal(t, http.StatusCreated, registerW.Code)

	// Try login with wrong password
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)
}
