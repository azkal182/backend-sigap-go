package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	provinceRepo := infraRepo.NewProvinceRepository()
	regencyRepo := infraRepo.NewRegencyRepository()
	districtRepo := infraRepo.NewDistrictRepository()
	villageRepo := infraRepo.NewVillageRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo)
	locationUseCase := usecase.NewLocationUseCase(provinceRepo, regencyRepo, districtRepo, villageRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	permissionHandler := handler.NewPermissionHandler(permissionUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router (includes global CORS middleware inside SetupRouter)
	r := router.SetupRouter(authHandler, userHandler, dormitoryHandler, roleHandler, locationHandler, permissionHandler, authMiddleware)

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
