package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"gorm.io/gorm"
)

// init registers all migrations
func init() {
	// Migration 001: Initial schema
	RegisterMigration(
		"001_initial_schema",
		"Create initial database schema",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.User{},
				&entity.Role{},
				&entity.Permission{},
				&entity.Dormitory{},
				&entity.UserRole{},
				&entity.RolePermission{},
				&entity.UserDormitory{},
			)
		},
		func(db *gorm.DB) error {
			// Rollback: Drop all tables in reverse order
			if err := db.Migrator().DropTable(&entity.UserDormitory{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.RolePermission{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.UserRole{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Dormitory{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Permission{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Role{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.User{}); err != nil {
				return err
			}
			return nil
		},
	)

	// Migration 002: Add IsProtected field to roles and seed default roles
	RegisterMigration(
		"002_add_role_protection_and_seed",
		"Add IsProtected field to roles table and seed default roles",
		func(db *gorm.DB) error {
			// Add IsProtected column if it doesn't exist
			if !db.Migrator().HasColumn(&entity.Role{}, "is_protected") {
				if err := db.Migrator().AddColumn(&entity.Role{}, "is_protected"); err != nil {
					return err
				}
			}

			// Seed default roles if they don't exist
			now := time.Now()
			defaultRoles := []entity.Role{
				{
					ID:          uuid.New(),
					Name:        "User",
					Slug:        "user",
					IsActive:    true,
					IsProtected: false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Admin",
					Slug:        "admin",
					IsActive:    true,
					IsProtected: true, // Protected role
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          uuid.New(),
					Name:        "Super Admin",
					Slug:        "super_admin",
					IsActive:    true,
					IsProtected: true, // Protected role
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}

			for _, role := range defaultRoles {
				var existingRole entity.Role
				result := db.Where("slug = ?", role.Slug).First(&existingRole)
				if result.Error == gorm.ErrRecordNotFound {
					// Role doesn't exist, create it
					if err := db.Create(&role).Error; err != nil {
						return err
					}
				} else if result.Error != nil {
					return result.Error
				}
				// If role exists, skip it
			}

			return nil
		},
		func(db *gorm.DB) error {
			// Rollback: Remove IsProtected column
			if db.Migrator().HasColumn(&entity.Role{}, "is_protected") {
				return db.Migrator().DropColumn(&entity.Role{}, "is_protected")
			}
			return nil
		},
	)

	// Migration 003: Remove address and capacity from dormitories
	RegisterMigration(
		"003_remove_dormitory_address_and_capacity",
		"Remove address and capacity columns from dormitories table",
		func(db *gorm.DB) error {
			// Drop address column if it exists
			if db.Migrator().HasColumn(&entity.Dormitory{}, "address") {
				if err := db.Migrator().DropColumn(&entity.Dormitory{}, "address"); err != nil {
					return err
				}
			}

			// Drop capacity column if it exists
			if db.Migrator().HasColumn(&entity.Dormitory{}, "capacity") {
				if err := db.Migrator().DropColumn(&entity.Dormitory{}, "capacity"); err != nil {
					return err
				}
			}

			return nil
		},
		func(db *gorm.DB) error {
			// Rollback: re-add address column if it does not exist
			if !db.Migrator().HasColumn(&entity.Dormitory{}, "address") {
				if err := db.Migrator().AddColumn(&entity.Dormitory{}, "address"); err != nil {
					return err
				}
			}

			// Rollback: re-add capacity column if it does not exist
			if !db.Migrator().HasColumn(&entity.Dormitory{}, "capacity") {
				if err := db.Migrator().AddColumn(&entity.Dormitory{}, "capacity"); err != nil {
					return err
				}
			}

			return nil
		},
	)
}
