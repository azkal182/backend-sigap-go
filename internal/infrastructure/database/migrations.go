package database

import (
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

	// Add more migrations here as needed
	// Example:
	// RegisterMigration(
	// 	"002_add_user_phone",
	// 	"Add phone field to users table",
	// 	func(db *gorm.DB) error {
	// 		return db.Migrator().AddColumn(&entity.User{}, "phone")
	// 	},
	// 	func(db *gorm.DB) error {
	// 		return db.Migrator().DropColumn(&entity.User{}, "phone")
	// 	},
	// )
}
