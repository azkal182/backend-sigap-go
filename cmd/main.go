package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	infraRepo "github.com/your-org/go-backend-starter/internal/infrastructure/repository"
	infraService "github.com/your-org/go-backend-starter/internal/infrastructure/service"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/handler"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/middleware"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/router"
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

	// Run migrations (using versioned migrations)
	// For production, use: go run cmd/migrate/main.go -command up
	if err := database.MigrateUpVersioned(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := infraRepo.NewUserRepository()
	roleRepo := infraRepo.NewRoleRepository()
	permissionRepo := infraRepo.NewPermissionRepository()
	dormitoryRepo := infraRepo.NewDormitoryRepository()
	studentRepo := infraRepo.NewStudentRepository()
	auditLogRepo := infraRepo.NewAuditLogRepository()
	provinceRepo := infraRepo.NewProvinceRepository()
	regencyRepo := infraRepo.NewRegencyRepository()
	districtRepo := infraRepo.NewDistrictRepository()
	villageRepo := infraRepo.NewVillageRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()
	auditLogger := service.NewAuditLogger(auditLogRepo)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, auditLogger)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, auditLogger)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo, auditLogger)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, dormitoryRepo, auditLogger)
	locationUseCase := usecase.NewLocationUseCase(provinceRepo, regencyRepo, districtRepo, villageRepo)
	auditLogUseCase := usecase.NewAuditLogUseCase(auditLogRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)
	studentHandler := handler.NewStudentHandler(studentUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	permissionHandler := handler.NewPermissionHandler(permissionUseCase)
	auditLogHandler := handler.NewAuditLogHandler(auditLogUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router (includes global CORS & audit context middleware inside SetupRouter)
	r := router.SetupRouter(authHandler, userHandler, dormitoryHandler, studentHandler, roleHandler, locationHandler, permissionHandler, auditLogHandler, authMiddleware)

	// Get server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
