package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Verify user was created
	var foundUser entity.User
	err = db.Where("id = ?", user.ID).First(&foundUser).Error
	require.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Name, foundUser.Name)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create a user first
	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Get user by ID
	foundUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)

	// Test not found
	_, err = repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create a user first
	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Get user by username
	foundUser, err := repo.GetByUsername(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)

	// Test not found
	_, err = repo.GetByUsername(ctx, "notfound")
	assert.Error(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create a user first
	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Update user
	user.Name = "Updated Name"
	user.UpdatedAt = time.Now()
	err := repo.Update(ctx, user)
	require.NoError(t, err)

	// Verify update
	var foundUser entity.User
	require.NoError(t, db.Where("id = ?", user.ID).First(&foundUser).Error)
	assert.Equal(t, "Updated Name", foundUser.Name)
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create a user first
	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(user).Error)

	// Delete user
	err := repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Verify deletion
	var foundUser entity.User
	err = db.Where("id = ?", user.ID).First(&foundUser).Error
	assert.Error(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &entity.User{
			ID:        uuid.New(),
			Username:  fmt.Sprintf("user%d", i),
			Password:  "hashedpassword",
			Name:      fmt.Sprintf("User %d", i),
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		require.NoError(t, db.Create(user).Error)
	}

	// List users
	users, total, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 5)

	// Test pagination
	users, total, err = repo.List(ctx, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 2)
}

func TestUserRepository_GetWithRoles(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := &userRepository{db: db}
	ctx := context.Background()

	// Create role
	role := &entity.Role{
		ID:        uuid.New(),
		Name:      "admin",
		Slug:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(role).Error)

	// Create user with role
	user := &entity.User{
		ID:        uuid.New(),
		Username:  "test",
		Password:  "hashedpassword",
		Name:      "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     []entity.Role{*role},
	}
	require.NoError(t, db.Create(user).Error)

	// Get user with roles
	userWithRoles, err := repo.GetWithRoles(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, userWithRoles.ID)
	assert.Len(t, userWithRoles.Roles, 1)
	assert.Equal(t, "admin", userWithRoles.Roles[0].Name)
}
