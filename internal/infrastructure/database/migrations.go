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

	// Migration 004: Create location tables (provinces, regencies, districts, villages)
	RegisterMigration(
		"004_create_location_tables",
		"Create provinces, regencies, districts, and villages tables",
		func(db *gorm.DB) error {
			// Use AutoMigrate for simplicity; IDs are ints and relations use FK columns
			return db.AutoMigrate(
				&entity.Province{},
				&entity.Regency{},
				&entity.District{},
				&entity.Village{},
			)
		},
		func(db *gorm.DB) error {
			// Rollback: drop tables in reverse dependency order
			if err := db.Migrator().DropTable(&entity.Village{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.District{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Regency{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Province{}); err != nil {
				return err
			}
			return nil
		},
	)

	// Migration 005: Add indexes for location tables to optimize search
	RegisterMigration(
		"005_add_location_indexes",
		"Add indexes on location tables for name search and parent filters",
		func(db *gorm.DB) error {
			// Provinces: index on name
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_provinces_name ON provinces (name)").Error; err != nil {
				return err
			}

			// Regencies: index on (province_id, name)
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_regencies_province_id_name ON regencies (province_id, name)").Error; err != nil {
				return err
			}

			// Districts: index on (regency_id, name)
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_districts_regency_id_name ON districts (regency_id, name)").Error; err != nil {
				return err
			}

			// Villages: index on (district_id, name)
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_villages_district_id_name ON villages (district_id, name)").Error; err != nil {
				return err
			}

			return nil
		},
		func(db *gorm.DB) error {
			// Drop indexes if they exist
			if err := db.Exec("DROP INDEX IF EXISTS idx_villages_district_id_name").Error; err != nil {
				return err
			}
			if err := db.Exec("DROP INDEX IF EXISTS idx_districts_regency_id_name").Error; err != nil {
				return err
			}
			if err := db.Exec("DROP INDEX IF EXISTS idx_regencies_province_id_name").Error; err != nil {
				return err
			}
			if err := db.Exec("DROP INDEX IF EXISTS idx_provinces_name").Error; err != nil {
				return err
			}
			return nil
		},
	)

	// Migration 006: Create audit_logs table
	RegisterMigration(
		"006_create_audit_logs",
		"Create audit_logs table for audit logging",
		func(db *gorm.DB) error {
			return db.AutoMigrate(&entity.AuditLog{})
		},
		func(db *gorm.DB) error {
			return db.Migrator().DropTable(&entity.AuditLog{})
		},
	)

	// Migration 007: Create students and student_dormitory_history tables
	RegisterMigration(
		"007_create_students",
		"Create students and student dormitory history tables",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.Student{},
				&entity.StudentDormitoryHistory{},
			)
		},
		func(db *gorm.DB) error {
			if err := db.Migrator().DropTable(&entity.StudentDormitoryHistory{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&entity.Student{})
		},
	)

	// Migration 008: Create fans and classes related tables
	RegisterMigration(
		"008_create_fans_and_classes",
		"Create fans, classes, student_class_enrollments, and class_staff tables",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.Fan{},
				&entity.Class{},
				&entity.StudentClassEnrollment{},
				&entity.ClassStaff{},
			)
		},
		func(db *gorm.DB) error {
			// Drop in reverse dependency order
			if err := db.Migrator().DropTable(&entity.ClassStaff{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.StudentClassEnrollment{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.Class{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&entity.Fan{})
		},
	)

	// Migration 009: Create teachers table
	RegisterMigration(
		"009_create_teachers",
		"Create teachers table for instructor records",
		func(db *gorm.DB) error {
			return db.AutoMigrate(&entity.Teacher{})
		},
		func(db *gorm.DB) error {
			return db.Migrator().DropTable(&entity.Teacher{})
		},
	)

	// Migration 010: Create schedule_slots table
	RegisterMigration(
		"010_create_schedule_slots",
		"Create schedule_slots table for dormitory time slots",
		func(db *gorm.DB) error {
			return db.AutoMigrate(&entity.ScheduleSlot{})
		},
		func(db *gorm.DB) error {
			return db.Migrator().DropTable(&entity.ScheduleSlot{})
		},
	)

	RegisterMigration(
		"011_create_subjects_class_sks",
		"Create subjects, class_schedules, sks_definitions, and sks_exam_schedules tables",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.Subject{},
				&entity.ClassSchedule{},
				&entity.SKSDefinition{},
				&entity.SKSExamSchedule{},
			)
		},
		func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&entity.SKSExamSchedule{},
				&entity.SKSDefinition{},
				&entity.ClassSchedule{},
				&entity.Subject{},
			)
		},
	)

	RegisterMigration(
		"012_create_student_sks_results",
		"Create student_sks_results and fan_completion_status tables",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.StudentSKSResult{},
				&entity.FanCompletionStatus{},
			)
		},
		func(db *gorm.DB) error {
			if err := db.Migrator().DropTable(&entity.FanCompletionStatus{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&entity.StudentSKSResult{})
		},
	)

	RegisterMigration(
		"013_create_attendance_tables",
		"Create attendance_sessions, student_attendances, and teacher_attendances",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.AttendanceSession{},
				&entity.StudentAttendance{},
				&entity.TeacherAttendance{},
			)
		},
		func(db *gorm.DB) error {
			if err := db.Migrator().DropTable(&entity.TeacherAttendance{}); err != nil {
				return err
			}
			if err := db.Migrator().DropTable(&entity.StudentAttendance{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&entity.AttendanceSession{})
		},
	)

	RegisterMigration(
		"014_create_leave_health_tables",
		"Create leave_permits and health_statuses tables",
		func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.LeavePermit{},
				&entity.HealthStatus{},
			)
		},
		func(db *gorm.DB) error {
			if err := db.Migrator().DropTable(&entity.HealthStatus{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&entity.LeavePermit{})
		},
	)

	RegisterMigration(
		"015_add_operational_indexes",
		"Add composite indexes for attendance, leave permits, health statuses, and dormitory history",
		func(db *gorm.DB) error {
			stmts := []string{
				"CREATE INDEX IF NOT EXISTS idx_attendance_sessions_schedule_date_status ON attendance_sessions (class_schedule_id, date, status)",
				"CREATE INDEX IF NOT EXISTS idx_attendance_sessions_teacher_date ON attendance_sessions (teacher_id, date)",
				"CREATE INDEX IF NOT EXISTS idx_leave_permits_student_status_dates ON leave_permits (student_id, status, start_date, end_date)",
				"CREATE INDEX IF NOT EXISTS idx_health_statuses_student_status_dates ON health_statuses (student_id, status, start_date, COALESCE(end_date, '9999-12-31'))",
				"CREATE INDEX IF NOT EXISTS idx_student_dormitory_history_student_end_start ON student_dormitory_history (student_id, end_date, start_date DESC)",
			}
			for _, stmt := range stmts {
				if err := db.Exec(stmt).Error; err != nil {
					return err
				}
			}
			return nil
		},
		func(db *gorm.DB) error {
			stmts := []string{
				"DROP INDEX IF EXISTS idx_student_dormitory_history_student_end_start",
				"DROP INDEX IF EXISTS idx_health_statuses_student_status_dates",
				"DROP INDEX IF EXISTS idx_leave_permits_student_status_dates",
				"DROP INDEX IF EXISTS idx_attendance_sessions_teacher_date",
				"DROP INDEX IF EXISTS idx_attendance_sessions_schedule_date_status",
			}
			for _, stmt := range stmts {
				if err := db.Exec(stmt).Error; err != nil {
					return err
				}
			}
			return nil
		},
	)
}
