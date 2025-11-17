package database

import (
	"fmt"
	"log"
	"os"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect initializes the database connection
func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate runs database migrations (legacy - uses AutoMigrate)
// For versioned migrations, use MigrateUp() instead
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	err := DB.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.Dormitory{},
		&entity.UserRole{},
		&entity.RolePermission{},
		&entity.UserDormitory{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// MigrateUpVersioned runs versioned migrations
// This is a wrapper that ensures migrations are registered
func MigrateUpVersioned() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	return MigrateUp(DB)
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
