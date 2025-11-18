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
			Description: "Main dormitory building",
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Dormitory B",
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
