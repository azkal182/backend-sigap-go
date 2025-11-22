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
	fanRepo := infraRepo.NewFanRepository()
	classRepo := infraRepo.NewClassRepository()
	enrollmentRepo := infraRepo.NewStudentClassEnrollmentRepository()
	classStaffRepo := infraRepo.NewClassStaffRepository()
	teacherRepo := infraRepo.NewTeacherRepository()
	subjectRepo := infraRepo.NewSubjectRepository()
	classScheduleRepo := infraRepo.NewClassScheduleRepository()
	leavePermitRepo := infraRepo.NewLeavePermitRepository()
	healthStatusRepo := infraRepo.NewHealthStatusRepository()
	sksDefinitionRepo := infraRepo.NewSKSDefinitionRepository()
	sksExamRepo := infraRepo.NewSKSExamScheduleRepository()
	studentSKSResultRepo := infraRepo.NewStudentSKSResultRepository()
	fanCompletionRepo := infraRepo.NewFanCompletionStatusRepository()
	attendanceSessionRepo := infraRepo.NewAttendanceSessionRepository()
	studentAttendanceRepo := infraRepo.NewStudentAttendanceRepository()
	teacherAttendanceRepo := infraRepo.NewTeacherAttendanceRepository()
	scheduleSlotRepo := infraRepo.NewScheduleSlotRepository()
	auditLogRepo := infraRepo.NewAuditLogRepository()
	provinceRepo := infraRepo.NewProvinceRepository()
	regencyRepo := infraRepo.NewRegencyRepository()
	districtRepo := infraRepo.NewDistrictRepository()
	villageRepo := infraRepo.NewVillageRepository()
	reportRepo := infraRepo.NewReportRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()
	auditLogger := service.NewAuditLogger(auditLogRepo)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, auditLogger)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, auditLogger)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo, auditLogger)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, dormitoryRepo, auditLogger)
	studentSKSResultUseCase := usecase.NewStudentSKSResultUseCase(studentSKSResultRepo, fanCompletionRepo, studentRepo, sksDefinitionRepo, teacherRepo, auditLogger)
	fanUseCase := usecase.NewFanUseCase(fanRepo, dormitoryRepo, auditLogger)
	classUseCase := usecase.NewClassUseCase(classRepo, fanRepo, studentRepo, enrollmentRepo, classStaffRepo, auditLogger)
	teacherUseCase := usecase.NewTeacherUseCase(teacherRepo, userRepo, roleRepo, auditLogger)
	scheduleSlotUseCase := usecase.NewScheduleSlotUseCase(scheduleSlotRepo, dormitoryRepo, auditLogger)
	classScheduleUseCase := usecase.NewClassScheduleUseCase(classScheduleRepo, classRepo, teacherRepo, subjectRepo, scheduleSlotRepo, dormitoryRepo, auditLogger)
	sksDefinitionUseCase := usecase.NewSKSDefinitionUseCase(sksDefinitionRepo, fanRepo, subjectRepo, auditLogger)
	sksExamUseCase := usecase.NewSKSExamScheduleUseCase(sksExamRepo, sksDefinitionRepo, teacherRepo, auditLogger)
	leavePermitUseCase := usecase.NewLeavePermitUseCase(leavePermitRepo, studentRepo, auditLogger)
	healthStatusUseCase := usecase.NewHealthStatusUseCase(healthStatusRepo, studentRepo, auditLogger)
	attendanceUseCase := usecase.NewAttendanceUseCase(attendanceSessionRepo, studentAttendanceRepo, teacherAttendanceRepo, classScheduleRepo, leavePermitUseCase, healthStatusUseCase, auditLogger)
	locationUseCase := usecase.NewLocationUseCase(provinceRepo, regencyRepo, districtRepo, villageRepo)
	auditLogUseCase := usecase.NewAuditLogUseCase(auditLogRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	reportUseCase := usecase.NewReportUseCase(reportRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)
	studentHandler := handler.NewStudentHandler(studentUseCase, studentSKSResultUseCase)
	fanHandler := handler.NewFanHandler(fanUseCase)
	classHandler := handler.NewClassHandler(classUseCase)
	teacherHandler := handler.NewTeacherHandler(teacherUseCase)
	scheduleSlotHandler := handler.NewScheduleSlotHandler(scheduleSlotUseCase)
	classScheduleHandler := handler.NewClassScheduleHandler(classScheduleUseCase)
	sksDefinitionHandler := handler.NewSKSDefinitionHandler(sksDefinitionUseCase)
	sksExamHandler := handler.NewSKSExamScheduleHandler(sksExamUseCase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUseCase)
	leavePermitHandler := handler.NewLeavePermitHandler(leavePermitUseCase)
	healthStatusHandler := handler.NewHealthStatusHandler(healthStatusUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	permissionHandler := handler.NewPermissionHandler(permissionUseCase)
	auditLogHandler := handler.NewAuditLogHandler(auditLogUseCase)
	reportHandler := handler.NewReportHandler(reportUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router (includes global CORS & audit context middleware inside SetupRouter)
	r := router.SetupRouter(
		authHandler,
		userHandler,
		dormitoryHandler,
		studentHandler,
		roleHandler,
		locationHandler,
		permissionHandler,
		auditLogHandler,
		fanHandler,
		classHandler,
		teacherHandler,
		classScheduleHandler,
		sksDefinitionHandler,
		sksExamHandler,
		attendanceHandler,
		scheduleSlotHandler,
		leavePermitHandler,
		healthStatusHandler,
		reportHandler,
		authMiddleware,
	)

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
