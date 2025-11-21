package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	infraRepo "github.com/your-org/go-backend-starter/internal/infrastructure/repository"
)

func main() {
	var dateInput string
	flag.StringVar(&dateInput, "date", "", "Date (YYYY-MM-DD) to lock attendance sessions. Defaults to today in server timezone.")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if dateInput == "" {
		dateInput = time.Now().Format("2006-01-02")
	}

	ctx := context.Background()

	attendanceSessionRepo := infraRepo.NewAttendanceSessionRepository()
	studentAttendanceRepo := infraRepo.NewStudentAttendanceRepository()
	teacherAttendanceRepo := infraRepo.NewTeacherAttendanceRepository()
	classScheduleRepo := infraRepo.NewClassScheduleRepository()
	auditLogRepo := infraRepo.NewAuditLogRepository()

	auditLogger := service.NewAuditLogger(auditLogRepo)
	attendanceUseCase := usecase.NewAttendanceUseCase(attendanceSessionRepo, studentAttendanceRepo, teacherAttendanceRepo, classScheduleRepo, auditLogger)

	req := dto.LockAttendanceRequest{Date: dateInput}
	if err := attendanceUseCase.LockSessions(ctx, req); err != nil {
		log.Fatalf("Failed to lock attendance sessions: %v", err)
	}

	log.Printf("Attendance sessions locked for %s", dateInput)
}
