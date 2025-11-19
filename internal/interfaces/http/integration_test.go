package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/service"
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

func (r *testUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
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

func (r *testUserRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Create(&entity.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *testUserRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entity.UserRole{}).Error
}

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, service.TokenService, func()) {
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
	studentRepo := infraRepo.NewStudentRepository()
	permissionRepo := infraRepo.NewPermissionRepository()
	auditLogRepo := infraRepo.NewAuditLogRepository()
	provinceRepo := infraRepo.NewProvinceRepository()
	regencyRepo := infraRepo.NewRegencyRepository()
	districtRepo := infraRepo.NewDistrictRepository()
	villageRepo := infraRepo.NewVillageRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()
	auditLogger := appService.NewAuditLogger(auditLogRepo)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, auditLogger)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo, auditLogger)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, dormitoryRepo, auditLogger)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, auditLogger)
	locationUseCase := usecase.NewLocationUseCase(provinceRepo, regencyRepo, districtRepo, villageRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	auditLogUseCase := usecase.NewAuditLogUseCase(auditLogRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)
	studentHandler := handler.NewStudentHandler(studentUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	permissionHandler := handler.NewPermissionHandler(permissionUseCase)
	auditLogHandler := handler.NewAuditLogHandler(auditLogUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router
	r := router.SetupRouter(authHandler, userHandler, dormitoryHandler, studentHandler, roleHandler, locationHandler, permissionHandler, auditLogHandler, authMiddleware)

	cleanup := func() {
		database.DB = originalDB // Restore original DB
		testutil.CleanupTestDB(t, testDB)
		testutil.UnsetTestEnv()
	}

	return r, testDB, tokenService, cleanup
}

func seedDormitory(t *testing.T, db *gorm.DB, name string) entity.Dormitory {
	dorm := entity.Dormitory{
		ID:        uuid.New(),
		Name:      name,
		Gender:    "male",
		Level:     "senior",
		Code:      fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&dorm).Error; err != nil {
		t.Fatalf("failed to seed dormitory: %v", err)
	}
	return dorm
}

// createTestUser creates an active user and returns user + bearer token
func createTestUser(t *testing.T, db *gorm.DB, username string, tokenService service.TokenService, permissions ...string) (entity.User, string) {
	user := entity.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  "$2a$10$Q9f1iG.zRV/X9sYt8GvGle6hwzEwA9H9n1tFoZT3zh0TBTtPlqHcC", // bcrypt hash for "password123"
		Name:      "Dorm Staff",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if len(permissions) == 0 {
		permissions = []string{"dorm:update"}
	}

	token, err := tokenService.GenerateAccessToken(user.ID, user.Username, permissions)
	if err != nil {
		t.Fatalf("failed generating access token: %v", err)
	}

	return user, token
}

func assignDormitoryAdminRole(t *testing.T, db *gorm.DB, userID uuid.UUID) {
	var perm entity.Permission
	if err := db.Where("name = ?", "dorm:update").First(&perm).Error; err != nil {
		perm = entity.Permission{
			ID:        uuid.New(),
			Name:      "dorm:update",
			Slug:      fmt.Sprintf("dorm-update-%d", time.Now().UnixNano()),
			Resource:  "dorm",
			Action:    "update",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&perm).Error; err != nil {
			t.Fatalf("failed to create permission: %v", err)
		}
	}

	role := entity.Role{
		ID:        uuid.New(),
		Name:      "admin",
		Slug:      fmt.Sprintf("admin-%d", time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to create role: %v", err)
	}

	rolePerm := entity.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
	if err := db.Create(&rolePerm).Error; err != nil {
		t.Fatalf("failed to link role permission: %v", err)
	}

	userRole := entity.UserRole{UserID: userID, RoleID: role.ID}
	if err := db.Create(&userRole).Error; err != nil {
		t.Fatalf("failed to assign role to user: %v", err)
	}
}

func assignStudentAdminRole(t *testing.T, db *gorm.DB, userID uuid.UUID) {
	permNames := []string{"student:read", "student:create", "student:update"}
	permissions := make([]entity.Permission, 0, len(permNames))
	for _, name := range permNames {
		var perm entity.Permission
		if err := db.Where("name = ?", name).First(&perm).Error; err != nil {
			perm = entity.Permission{
				ID:        uuid.New(),
				Name:      name,
				Slug:      fmt.Sprintf("%s-%d", strings.ReplaceAll(name, ":", "-"), time.Now().UnixNano()),
				Resource:  "student",
				Action:    strings.Split(name, ":")[1],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&perm).Error; err != nil {
				t.Fatalf("failed to create permission %s: %v", name, err)
			}
		}
		permissions = append(permissions, perm)
	}

	role := entity.Role{
		ID:        uuid.New(),
		Name:      "student-admin",
		Slug:      fmt.Sprintf("student-admin-%d", time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to create student role: %v", err)
	}

	for _, perm := range permissions {
		rp := entity.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
		if err := db.Create(&rp).Error; err != nil {
			t.Fatalf("failed to assign permission %s: %v", perm.Name, err)
		}
	}

	userRole := entity.UserRole{UserID: userID, RoleID: role.ID}
	if err := db.Create(&userRole).Error; err != nil {
		t.Fatalf("failed to assign student role: %v", err)
	}
}

func TestAuthIntegration_RegisterAndLogin(t *testing.T) {
	router, _, _, cleanup := setupTestRouter(t)
	defer cleanup()

	// Test Register
	registerReq := dto.RegisterRequest{
		Username: "integration",
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
		Username: "integration",
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
	// protectedReq, _ := http.NewRequest(http.MethodGet, "/api/users/me", nil)
	// protectedReq.Header.Set("Authorization", "Bearer "+accessToken)

	// protectedW := httptest.NewRecorder()
	// router.ServeHTTP(protectedW, protectedReq)

	// Note: This test assumes there's a /api/users/me endpoint
	// Adjust based on your actual routes
	// assert.True(t, protectedW.Code == http.StatusOK || protectedW.Code == http.StatusNotFound)
}

func TestAuthIntegration_InvalidCredentials(t *testing.T) {
	router, _, _, cleanup := setupTestRouter(t)
	defer cleanup()

	// Register user first
	registerReq := dto.RegisterRequest{
		Username: "test",
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
		Username: "test",
		Password: "wrongpassword",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)
}

func TestDormitoryIntegration_AssignAndRemoveUser(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Integration Dorm")
	user, token := createTestUser(t, db, "dorm-staff", tokenService)
	assignDormitoryAdminRole(t, db, user.ID)

	assignPayload := map[string]string{"user_id": user.ID.String()}
	body, _ := json.Marshal(assignPayload)
	assignReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/dormitories/%s/users", dorm.ID.String()), bytes.NewBuffer(body))
	assignReq.Header.Set("Content-Type", "application/json")
	assignReq.Header.Set("Authorization", "Bearer "+token)

	assignW := httptest.NewRecorder()
	router.ServeHTTP(assignW, assignReq)
	assert.Equal(t, http.StatusNoContent, assignW.Code)

	removeReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/dormitories/%s/users/%s", dorm.ID.String(), user.ID.String()), nil)
	removeReq.Header.Set("Authorization", "Bearer "+token)

	removeW := httptest.NewRecorder()
	router.ServeHTTP(removeW, removeReq)
	assert.Equal(t, http.StatusNoContent, removeW.Code)
}

func TestStudentIntegration_CreateStatusMutate(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Student Dorm")
	studentNumber := fmt.Sprintf("STD%d", time.Now().UnixNano())
	user, token := createTestUser(t, db, "student-admin", tokenService, "student:read", "student:create", "student:update")
	assignStudentAdminRole(t, db, user.ID)

	birthDate := time.Now().AddDate(-15, 0, 0).UTC().Format(time.RFC3339)
	createPayload := map[string]interface{}{
		"student_number": studentNumber,
		"full_name":      "Integration Student",
		"birth_date":     birthDate,
		"gender":         "male",
		"parent_name":    "Integration Parent",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/api/students", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	assert.Equal(t, http.StatusCreated, createW.Code)

	type studentResp struct {
		Success bool                `json:"success"`
		Data    dto.StudentResponse `json:"data"`
	}

	var created studentResp
	require.NoError(t, json.Unmarshal(createW.Body.Bytes(), &created))
	require.True(t, created.Success)
	require.NotEmpty(t, created.Data.ID)

	// Update status
	statusPayload := map[string]string{"status": entity.StudentStatusLeave}
	statusBody, _ := json.Marshal(statusPayload)
	statusReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/students/%s/status", created.Data.ID), bytes.NewBuffer(statusBody))
	statusReq.Header.Set("Content-Type", "application/json")
	statusReq.Header.Set("Authorization", "Bearer "+token)

	statusW := httptest.NewRecorder()
	router.ServeHTTP(statusW, statusReq)
	assert.Equal(t, http.StatusOK, statusW.Code)

	var statusResp studentResp
	require.NoError(t, json.Unmarshal(statusW.Body.Bytes(), &statusResp))
	assert.Equal(t, entity.StudentStatusLeave, statusResp.Data.Status)

	// Mutate dormitory
	mutatePayload := map[string]interface{}{
		"dormitory_id": dorm.ID.String(),
		"start_date":   time.Now().UTC().Format(time.RFC3339),
	}
	mutateBody, _ := json.Marshal(mutatePayload)
	mutateReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/students/%s/mutate-dormitory", created.Data.ID), bytes.NewBuffer(mutateBody))
	mutateReq.Header.Set("Content-Type", "application/json")
	mutateReq.Header.Set("Authorization", "Bearer "+token)

	mutateW := httptest.NewRecorder()
	router.ServeHTTP(mutateW, mutateReq)
	assert.Equal(t, http.StatusOK, mutateW.Code)

	var mutateResp studentResp
	require.NoError(t, json.Unmarshal(mutateW.Body.Bytes(), &mutateResp))
	require.Greater(t, len(mutateResp.Data.DormitoryHistory), 0)
	assert.Equal(t, dorm.ID.String(), mutateResp.Data.DormitoryHistory[0].DormitoryID)
}
