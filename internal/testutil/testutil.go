package testutil

import (
	"os"
	"testing"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.Dormitory{},
		&entity.Student{},
		&entity.StudentDormitoryHistory{},
		&entity.UserRole{},
		&entity.RolePermission{},
		&entity.UserDormitory{},
		&entity.AuditLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestDB closes the test database connection
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Error getting database instance: %v", err)
		return
	}
	sqlDB.Close()
}

// SetTestEnv sets test environment variables
func SetTestEnv() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRY", "15m")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRY", "168h")
}

// UnsetTestEnv unsets test environment variables
func UnsetTestEnv() {
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_ACCESS_TOKEN_EXPIRY")
	os.Unsetenv("JWT_REFRESH_TOKEN_EXPIRY")
}
