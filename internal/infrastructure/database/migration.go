package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// TableName specifies the table name for migrations
func (Migration) TableName() string {
	return "schema_migrations"
}

// MigrationFunc represents a migration function
type MigrationFunc func(*gorm.DB) error

// MigrationStep represents a single migration step
type MigrationStep struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

var migrations []MigrationStep

// RegisterMigration registers a new migration
func RegisterMigration(version, name string, up, down MigrationFunc) {
	migrations = append(migrations, MigrationStep{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// GetMigrations returns all registered migrations
func GetMigrations() []MigrationStep {
	return migrations
}

// EnsureMigrationTable ensures the migration tracking table exists
func EnsureMigrationTable(db *gorm.DB) error {
	return db.AutoMigrate(&Migration{})
}

// GetAppliedMigrations returns all applied migrations
func GetAppliedMigrations(db *gorm.DB) (map[string]bool, error) {
	var applied []Migration
	if err := db.Find(&applied).Error; err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	return appliedMap, nil
}

// MarkMigrationApplied marks a migration as applied
func MarkMigrationApplied(db *gorm.DB, version, name string) error {
	migration := Migration{
		Version:   version,
		Name:      name,
		AppliedAt: time.Now(),
	}
	return db.Create(&migration).Error
}

// MarkMigrationRolledBack removes a migration from the applied list
func MarkMigrationRolledBack(db *gorm.DB, version string) error {
	return db.Where("version = ?", version).Delete(&Migration{}).Error
}

// MigrateUp runs all pending migrations
func MigrateUp(db *gorm.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		if err := migration.Up(db); err != nil {
			return fmt.Errorf("failed to apply migration %s (%s): %w", migration.Version, migration.Name, err)
		}

		if err := MarkMigrationApplied(db, migration.Version, migration.Name); err != nil {
			return fmt.Errorf("failed to mark migration as applied: %w", err)
		}

		fmt.Printf("Applied migration: %s - %s\n", migration.Version, migration.Name)
	}

	return nil
}

// MigrateDown rolls back the last migration
func MigrateDown(db *gorm.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find the last applied migration
	var lastMigration *MigrationStep
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].Version] {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		return fmt.Errorf("no migrations to rollback")
	}

	if lastMigration.Down == nil {
		return fmt.Errorf("migration %s (%s) does not have a rollback function", lastMigration.Version, lastMigration.Name)
	}

	if err := lastMigration.Down(db); err != nil {
		return fmt.Errorf("failed to rollback migration %s (%s): %w", lastMigration.Version, lastMigration.Name, err)
	}

	if err := MarkMigrationRolledBack(db, lastMigration.Version); err != nil {
		return fmt.Errorf("failed to mark migration as rolled back: %w", err)
	}

	fmt.Printf("Rolled back migration: %s - %s\n", lastMigration.Version, lastMigration.Name)
	return nil
}

// MigrateToVersion migrates to a specific version
func MigrateToVersion(db *gorm.DB, targetVersion string) error {
	if err := EnsureMigrationTable(db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find target migration index
	targetIndex := -1
	for i, m := range migrations {
		if m.Version == targetVersion {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return fmt.Errorf("migration version %s not found", targetVersion)
	}

	// Apply or rollback migrations
	for i, migration := range migrations {
		isApplied := applied[migration.Version]

		if i <= targetIndex && !isApplied {
			// Need to apply
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			if err := MarkMigrationApplied(db, migration.Version, migration.Name); err != nil {
				return fmt.Errorf("failed to mark migration as applied: %w", err)
			}
			fmt.Printf("Applied migration: %s - %s\n", migration.Version, migration.Name)
		} else if i > targetIndex && isApplied {
			// Need to rollback
			if migration.Down == nil {
				return fmt.Errorf("migration %s does not have rollback function", migration.Version)
			}
			if err := migration.Down(db); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
			}
			if err := MarkMigrationRolledBack(db, migration.Version); err != nil {
				return fmt.Errorf("failed to mark migration as rolled back: %w", err)
			}
			fmt.Printf("Rolled back migration: %s - %s\n", migration.Version, migration.Name)
		}
	}

	return nil
}

// GetMigrationStatus returns the status of all migrations
func GetMigrationStatus(db *gorm.DB) ([]map[string]interface{}, error) {
	if err := EnsureMigrationTable(db); err != nil {
		return nil, fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	var status []map[string]interface{}
	for _, migration := range migrations {
		status = append(status, map[string]interface{}{
			"version":   migration.Version,
			"name":      migration.Name,
			"applied":   applied[migration.Version],
		})
	}

	return status, nil
}
