package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	infraRepo "github.com/your-org/go-backend-starter/internal/infrastructure/repository"
	"golang.org/x/crypto/bcrypt"
)

func findPermissionByName(perms []*entity.Permission, name string) *entity.Permission {
	for _, perm := range perms {
		if perm.Name == name {
			return perm
		}
	}
	return nil
}

func mustPermission(perms []*entity.Permission, name string) *entity.Permission {
	perm := findPermissionByName(perms, name)
	if perm == nil {
		log.Fatalf("permission %s not found", name)
	}
	return perm
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Initialize repositories
	permissionRepo := infraRepo.NewPermissionRepository()
	roleRepo := infraRepo.NewRoleRepository()
	userRepo := infraRepo.NewUserRepository()
	dormitoryRepo := infraRepo.NewDormitoryRepository()

	// Create permissions
	permissions := []*entity.Permission{
		// User permissions
		{ID: uuid.New(), Name: "user:read", Slug: "user-read", Resource: "user", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "user:create", Slug: "user-create", Resource: "user", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "user:update", Slug: "user-update", Resource: "user", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "user:delete", Slug: "user-delete", Resource: "user", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Dormitory permissions
		{ID: uuid.New(), Name: "dorm:read", Slug: "dorm-read", Resource: "dorm", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "dorm:create", Slug: "dorm-create", Resource: "dorm", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "dorm:update", Slug: "dorm-update", Resource: "dorm", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "dorm:delete", Slug: "dorm-delete", Resource: "dorm", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Role permissions
		{ID: uuid.New(), Name: "role:read", Slug: "role-read", Resource: "role", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "role:create", Slug: "role-create", Resource: "role", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "role:update", Slug: "role-update", Resource: "role", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "role:delete", Slug: "role-delete", Resource: "role", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Audit permissions
		{ID: uuid.New(), Name: "audit:read", Slug: "audit-read", Resource: "audit_log", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Student permissions
		{ID: uuid.New(), Name: "student:read", Slug: "student-read", Resource: "student", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "student:create", Slug: "student-create", Resource: "student", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "student:update", Slug: "student-update", Resource: "student", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Fan permissions
		{ID: uuid.New(), Name: "fans:read", Slug: "fans-read", Resource: "fans", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "fans:create", Slug: "fans-create", Resource: "fans", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "fans:update", Slug: "fans-update", Resource: "fans", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "fans:delete", Slug: "fans-delete", Resource: "fans", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Class permissions
		{ID: uuid.New(), Name: "classes:read", Slug: "classes-read", Resource: "classes", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "classes:create", Slug: "classes-create", Resource: "classes", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "classes:update", Slug: "classes-update", Resource: "classes", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "classes:delete", Slug: "classes-delete", Resource: "classes", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Teacher permissions
		{ID: uuid.New(), Name: "teachers:read", Slug: "teachers-read", Resource: "teachers", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "teachers:create", Slug: "teachers-create", Resource: "teachers", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "teachers:update", Slug: "teachers-update", Resource: "teachers", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "teachers:delete", Slug: "teachers-delete", Resource: "teachers", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Schedule slot permissions
		{ID: uuid.New(), Name: "schedule_slots:read", Slug: "schedule-slots-read", Resource: "schedule_slots", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "schedule_slots:create", Slug: "schedule-slots-create", Resource: "schedule_slots", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "schedule_slots:update", Slug: "schedule-slots-update", Resource: "schedule_slots", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "schedule_slots:delete", Slug: "schedule-slots-delete", Resource: "schedule_slots", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Class schedule permissions
		{ID: uuid.New(), Name: "class_schedules:read", Slug: "class-schedules-read", Resource: "class_schedules", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "class_schedules:create", Slug: "class-schedules-create", Resource: "class_schedules", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "class_schedules:update", Slug: "class-schedules-update", Resource: "class_schedules", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "class_schedules:delete", Slug: "class-schedules-delete", Resource: "class_schedules", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// SKS definition permissions
		{ID: uuid.New(), Name: "sks_definitions:read", Slug: "sks-definitions-read", Resource: "sks_definitions", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_definitions:create", Slug: "sks-definitions-create", Resource: "sks_definitions", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_definitions:update", Slug: "sks-definitions-update", Resource: "sks_definitions", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_definitions:delete", Slug: "sks-definitions-delete", Resource: "sks_definitions", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// SKS exam schedule permissions
		{ID: uuid.New(), Name: "sks_exams:read", Slug: "sks-exams-read", Resource: "sks_exams", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_exams:create", Slug: "sks-exams-create", Resource: "sks_exams", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_exams:update", Slug: "sks-exams-update", Resource: "sks_exams", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "sks_exams:delete", Slug: "sks-exams-delete", Resource: "sks_exams", Action: "delete", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Student SKS results & FAN completion permissions
		{ID: uuid.New(), Name: "student_sks_results:read", Slug: "student-sks-results-read", Resource: "student_sks_results", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "student_sks_results:create", Slug: "student-sks-results-create", Resource: "student_sks_results", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "student_sks_results:update", Slug: "student-sks-results-update", Resource: "student_sks_results", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		// Attendance permissions
		{ID: uuid.New(), Name: "attendance_sessions:read", Slug: "attendance-sessions-read", Resource: "attendance_sessions", Action: "read", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "attendance_sessions:create", Slug: "attendance-sessions-create", Resource: "attendance_sessions", Action: "create", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "attendance_sessions:update", Slug: "attendance-sessions-update", Resource: "attendance_sessions", Action: "update", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "attendance_sessions:lock", Slug: "attendance-sessions-lock", Resource: "attendance_sessions", Action: "lock", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	log.Println("Creating permissions...")
	for _, perm := range permissions {
		// Check if permission already exists
		existingPerm, _ := permissionRepo.GetBySlug(ctx, perm.Slug)
		if existingPerm != nil {
			log.Printf("Permission %s already exists, skipping...", perm.Name)
			// Update the permission ID in our array for later use
			perm.ID = existingPerm.ID
		} else {
			if err := permissionRepo.Create(ctx, perm); err != nil {
				log.Printf("Failed to create permission %s: %v", perm.Name, err)
			} else {
				log.Printf("Created permission: %s", perm.Name)
			}
		}
	}

	// Check if roles already exist (from migration)
	var existingUserRole, existingAdminRole, existingSuperAdminRole *entity.Role
	userRolePtr, _ := roleRepo.GetBySlug(ctx, "user")
	adminRolePtr, _ := roleRepo.GetBySlug(ctx, "admin")
	superAdminRolePtr, _ := roleRepo.GetBySlug(ctx, "super_admin")

	userRoleExists := userRolePtr != nil
	adminRoleExists := adminRolePtr != nil
	superAdminRoleExists := superAdminRolePtr != nil

	if userRoleExists {
		existingUserRole = userRolePtr
	}
	if adminRoleExists {
		existingAdminRole = adminRolePtr
	}
	if superAdminRoleExists {
		existingSuperAdminRole = superAdminRolePtr
	}

	// Create or update roles with permissions
	log.Println("Setting up roles and permissions...")

	// User Role (default role, not protected)
	if !userRoleExists {
		userRole := &entity.Role{
			ID:          uuid.New(),
			Name:        "User",
			Slug:        "user",
			IsActive:    true,
			IsProtected: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{
				*permissions[4], // dorm:read
			},
		}
		if err := roleRepo.Create(ctx, userRole); err != nil {
			log.Printf("Failed to create user role: %v", err)
		} else {
			log.Println("Created user role")
		}
	} else {
		// Assign permissions to existing user role
		if existingUserRole != nil {
			roleRepo.AssignPermission(ctx, existingUserRole.ID, permissions[4].ID) // dorm:read
			log.Println("Updated user role permissions")
		}
	}

	// Teacher Role (assign teacher permissions, not protected)
	teacherPerms := []*entity.Permission{permissions[23], permissions[24], permissions[25], permissions[26]}
	teacherRole, _ := roleRepo.GetBySlug(ctx, "teacher")
	if teacherRole == nil {
		teacherRoleEntity := &entity.Role{
			ID:          uuid.New(),
			Name:        "Teacher",
			Slug:        "teacher",
			IsActive:    true,
			IsProtected: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{
				*teacherPerms[0], *teacherPerms[1], *teacherPerms[2], *teacherPerms[3],
			},
		}
		if err := roleRepo.Create(ctx, teacherRoleEntity); err != nil {
			log.Printf("Failed to create teacher role: %v", err)
		} else {
			log.Println("Created teacher role")
		}
	} else {
		for _, perm := range teacherPerms {
			roleRepo.AssignPermission(ctx, teacherRole.ID, perm.ID)
		}
		log.Println("Updated teacher role permissions")
	}

	attendanceReadPerm := mustPermission(permissions, "attendance_sessions:read")
	attendanceCreatePerm := mustPermission(permissions, "attendance_sessions:create")
	attendanceUpdatePerm := mustPermission(permissions, "attendance_sessions:update")
	attendanceLockPerm := mustPermission(permissions, "attendance_sessions:lock")

	// Admin Role (protected)
	if !adminRoleExists {
		adminRole := &entity.Role{
			ID:          uuid.New(),
			Name:        "Admin",
			Slug:        "admin",
			IsActive:    true,
			IsProtected: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{
				*permissions[0], *permissions[1], *permissions[2], *permissions[3], // user:*
				*permissions[4], *permissions[5], *permissions[6], *permissions[7], // dorm:*
				*permissions[8], *permissions[9], *permissions[10], *permissions[11], // role:*
				*permissions[12], *permissions[13], *permissions[14], // student:*
				*permissions[15], *permissions[16], *permissions[17], *permissions[18], // fans:*
				*permissions[19], *permissions[20], *permissions[21], *permissions[22], // classes:*
				*permissions[23], *permissions[24], *permissions[25], *permissions[26], // teachers:*
				*permissions[27], *permissions[28], *permissions[29], *permissions[30], // schedule slots:*
				*permissions[31], *permissions[32], *permissions[33], *permissions[34], // class schedules:*
				*permissions[35], *permissions[36], *permissions[37], *permissions[38], // sks definitions:*
				*permissions[39], *permissions[40], *permissions[41], *permissions[42], // sks exams:*
				*permissions[43], *permissions[44], *permissions[45], // student sks results:*
				*attendanceReadPerm, *attendanceCreatePerm, *attendanceUpdatePerm, *attendanceLockPerm, // attendance
			},
		}
		if err := roleRepo.Create(ctx, adminRole); err != nil {
			log.Printf("Failed to create admin role: %v", err)
		} else {
			log.Println("Created admin role")
		}
	} else {
		// Assign all permissions to existing admin role
		if existingAdminRole != nil {
			for _, perm := range permissions {
				roleRepo.AssignPermission(ctx, existingAdminRole.ID, perm.ID)
			}
			log.Println("Updated admin role permissions")
		}
	}

	// Super Admin Role (protected, has all permissions)
	if !superAdminRoleExists {
		superAdminRole := &entity.Role{
			ID:          uuid.New(),
			Name:        "Super Admin",
			Slug:        "super_admin",
			IsActive:    true,
			IsProtected: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{
				*permissions[0], *permissions[1], *permissions[2], *permissions[3], // user:*
				*permissions[4], *permissions[5], *permissions[6], *permissions[7], // dorm:*
				*permissions[8], *permissions[9], *permissions[10], *permissions[11], // role:*
				*permissions[12], *permissions[13], *permissions[14], // student:*
			},
		}
		if err := roleRepo.Create(ctx, superAdminRole); err != nil {
			log.Printf("Failed to create super admin role: %v", err)
		} else {
			log.Println("Created super admin role")
		}
	} else {
		// Assign all permissions to existing super admin role
		if existingSuperAdminRole != nil {
			for _, perm := range permissions {
				roleRepo.AssignPermission(ctx, existingSuperAdminRole.ID, perm.ID)
			}
			log.Println("Updated super admin role permissions")
		}
	}

	// Academic SKS role (attendance management)
	academicRole, _ := roleRepo.GetBySlug(ctx, "academic_sks")
	academicPerms := []*entity.Permission{attendanceReadPerm, attendanceCreatePerm, attendanceUpdatePerm}
	if academicRole == nil {
		roleEntity := &entity.Role{
			ID:          uuid.New(),
			Name:        "Academic SKS",
			Slug:        "academic_sks",
			IsActive:    true,
			IsProtected: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{*academicPerms[0], *academicPerms[1], *academicPerms[2]},
		}
		if err := roleRepo.Create(ctx, roleEntity); err != nil {
			log.Printf("Failed to create academic SKS role: %v", err)
		} else {
			log.Println("Created academic SKS role")
		}
	} else {
		for _, perm := range academicPerms {
			roleRepo.AssignPermission(ctx, academicRole.ID, perm.ID)
		}
		log.Println("Updated academic SKS role permissions")
	}

	// Attendance cron role (lock sessions)
	cronRole, _ := roleRepo.GetBySlug(ctx, "attendance_cron")
	cronPerms := []*entity.Permission{attendanceReadPerm, attendanceLockPerm}
	if cronRole == nil {
		roleEntity := &entity.Role{
			ID:          uuid.New(),
			Name:        "Attendance Cron",
			Slug:        "attendance_cron",
			IsActive:    true,
			IsProtected: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []entity.Permission{*cronPerms[0], *cronPerms[1]},
		}
		if err := roleRepo.Create(ctx, roleEntity); err != nil {
			log.Printf("Failed to create attendance cron role: %v", err)
		} else {
			log.Println("Created attendance cron role")
		}
	} else {
		for _, perm := range cronPerms {
			roleRepo.AssignPermission(ctx, cronRole.ID, perm.ID)
		}
		log.Println("Updated attendance cron role permissions")
	}

	// Get roles for user assignment
	var adminRoleEntity, superAdminRoleEntity *entity.Role
	if adminRoleExists && existingAdminRole != nil {
		adminRoleEntity = existingAdminRole
	} else {
		adminRoleEntity, _ = roleRepo.GetBySlug(ctx, "admin")
	}
	if superAdminRoleExists && existingSuperAdminRole != nil {
		superAdminRoleEntity = existingSuperAdminRole
	} else {
		superAdminRoleEntity, _ = roleRepo.GetBySlug(ctx, "super_admin")
	}

	// Create admin user (only if doesn't exist)
	existingAdminUser, _ := userRepo.GetByUsername(ctx, "admin")
	if existingAdminUser == nil && adminRoleEntity != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser := &entity.User{
			ID:        uuid.New(),
			Username:  "admin",
			Password:  string(hashedPassword),
			Name:      "Admin User",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Roles:     []entity.Role{*adminRoleEntity},
		}

		log.Println("Creating admin user...")
		if err := userRepo.Create(ctx, adminUser); err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("Created admin user: admin / admin123")
		}
	} else {
		log.Println("Admin user already exists, skipping...")
	}

	// Create super admin user (only if doesn't exist)
	existingSuperAdminUser, _ := userRepo.GetByUsername(ctx, "superadmin")
	if existingSuperAdminUser == nil && superAdminRoleEntity != nil {
		hashedPasswordSuper, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)
		superAdminUser := &entity.User{
			ID:        uuid.New(),
			Username:  "superadmin",
			Password:  string(hashedPasswordSuper),
			Name:      "Super Admin User",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Roles:     []entity.Role{*superAdminRoleEntity},
		}

		log.Println("Creating super admin user...")
		if err := userRepo.Create(ctx, superAdminUser); err != nil {
			log.Printf("Failed to create super admin user: %v", err)
		} else {
			log.Println("Created super admin user: superadmin / superadmin123")
		}
	} else {
		log.Println("Super admin user already exists, skipping...")
	}

	// Create sample dormitories
	dormitories := []*entity.Dormitory{
		{
			ID:          uuid.New(),
			Name:        "Dormitory A",
			Gender:      "male",
			Level:       "senior",
			Code:        "DORMA",
			Description: "Main dormitory building",
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Dormitory B",
			Gender:      "female",
			Level:       "junior",
			Code:        "DORMB",
			Description: "Secondary dormitory building",
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	log.Println("Creating dormitories...")
	for _, dorm := range dormitories {
		if err := dormitoryRepo.Create(ctx, dorm); err != nil {
			log.Printf("Failed to create dormitory %s: %v", dorm.Name, err)
		} else {
			log.Printf("Created dormitory: %s", dorm.Name)
		}
	}

	log.Println("Seed data created successfully!")
}
