package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// StudentUseCase orchestrates student-related use cases.
type StudentUseCase struct {
	studentRepo repository.StudentRepository
	dormRepo    repository.DormitoryRepository
	auditLogger appService.AuditLogger
}

// NewStudentUseCase builds StudentUseCase instance.
func NewStudentUseCase(
	studentRepo repository.StudentRepository,
	dormRepo repository.DormitoryRepository,
	auditLogger appService.AuditLogger,
) *StudentUseCase {
	return &StudentUseCase{studentRepo: studentRepo, dormRepo: dormRepo, auditLogger: auditLogger}
}

// CreateStudent creates a new student record.
func (uc *StudentUseCase) CreateStudent(ctx context.Context, req dto.CreateStudentRequest) (*dto.StudentResponse, error) {
	existing, _ := uc.studentRepo.GetByStudentNumber(ctx, req.StudentNumber)
	if existing != nil {
		return nil, domainErrors.ErrStudentAlreadyExists
	}

	now := time.Now()
	student := &entity.Student{
		ID:            uuid.New(),
		StudentNumber: req.StudentNumber,
		FullName:      req.FullName,
		BirthDate:     req.BirthDate,
		Gender:        req.Gender,
		ParentName:    req.ParentName,
		Status:        entity.StudentStatusActive,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.studentRepo.Create(ctx, student); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "student", "student:create", student.ID.String(), map[string]string{
		"student_number": student.StudentNumber,
		"full_name":      student.FullName,
	})

	return uc.toStudentResponse(student, nil), nil
}

// GetStudentByID retrieves a student along with history.
func (uc *StudentUseCase) GetStudentByID(ctx context.Context, id uuid.UUID) (*dto.StudentResponse, error) {
	student, err := uc.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	histories, _ := uc.studentRepo.ListHistory(ctx, id)
	return uc.toStudentResponse(student, histories), nil
}

// ListStudents returns paginated students.
func (uc *StudentUseCase) ListStudents(ctx context.Context, page, pageSize int) (*dto.ListStudentsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	students, total, err := uc.studentRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.StudentResponse, 0, len(students))
	for _, student := range students {
		responses = append(responses, *uc.toStudentResponse(student, nil))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListStudentsResponse{
		Students:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateStudent updates student profile fields.
func (uc *StudentUseCase) UpdateStudent(ctx context.Context, id uuid.UUID, req dto.UpdateStudentRequest) (*dto.StudentResponse, error) {
	student, err := uc.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	if req.FullName != nil {
		student.FullName = *req.FullName
	}
	if req.BirthDate != nil {
		student.BirthDate = *req.BirthDate
	}
	if req.Gender != nil {
		student.Gender = *req.Gender
	}
	if req.ParentName != nil {
		student.ParentName = *req.ParentName
	}
	student.UpdatedAt = time.Now()

	if err := uc.studentRepo.Update(ctx, student); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	histories, _ := uc.studentRepo.ListHistory(ctx, id)
	_ = uc.auditLogger.Log(ctx, "student", "student:update", student.ID.String(), map[string]string{
		"student_number": student.StudentNumber,
	})

	return uc.toStudentResponse(student, histories), nil
}

// UpdateStudentStatus updates lifecycle status and related active flag.
func (uc *StudentUseCase) UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) (*dto.StudentResponse, error) {
	student, err := uc.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	isActive := status == entity.StudentStatusActive
	if err := uc.studentRepo.UpdateStatus(ctx, id, status, isActive); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	student.Status = status
	student.IsActive = isActive
	student.UpdatedAt = time.Now()

	histories, _ := uc.studentRepo.ListHistory(ctx, id)
	_ = uc.auditLogger.Log(ctx, "student", "student:update-status", student.ID.String(), map[string]string{
		"status": status,
	})

	return uc.toStudentResponse(student, histories), nil
}

// MutateStudentDormitory creates dormitory history entry.
func (uc *StudentUseCase) MutateStudentDormitory(ctx context.Context, studentID, dormitoryID uuid.UUID, startDate time.Time) (*dto.StudentResponse, error) {
	student, err := uc.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	if _, err := uc.dormRepo.GetByID(ctx, dormitoryID); err != nil {
		return nil, domainErrors.ErrDormitoryNotFound
	}

	if startDate.IsZero() {
		startDate = time.Now()
	}
	now := time.Now()

	if currentHistory, err := uc.studentRepo.GetActiveHistory(ctx, studentID); err == nil && currentHistory != nil {
		if err := uc.studentRepo.CloseHistory(ctx, currentHistory.ID, startDate); err != nil {
			return nil, domainErrors.ErrInternalServer
		}
	}

	history := &entity.StudentDormitoryHistory{
		ID:          uuid.New(),
		StudentID:   studentID,
		DormitoryID: dormitoryID,
		StartDate:   startDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.studentRepo.CreateHistory(ctx, history); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	histories, _ := uc.studentRepo.ListHistory(ctx, studentID)
	_ = uc.auditLogger.Log(ctx, "student", "student:mutate-dorm", student.ID.String(), map[string]string{
		"dormitory_id": dormitoryID.String(),
	})

	return uc.toStudentResponse(student, histories), nil
}

func (uc *StudentUseCase) toStudentResponse(student *entity.Student, histories []*entity.StudentDormitoryHistory) *dto.StudentResponse {
	var historyResponses []dto.StudentDormitoryEvent
	if len(histories) > 0 {
		historyResponses = make([]dto.StudentDormitoryEvent, 0, len(histories))
		for _, history := range histories {
			event := dto.StudentDormitoryEvent{
				DormitoryID: history.DormitoryID.String(),
				StartDate:   history.StartDate.Format(time.RFC3339),
			}
			if history.EndDate != nil {
				event.EndDate = history.EndDate.Format(time.RFC3339)
			}
			historyResponses = append(historyResponses, event)
		}
	}

	birthDate := ""
	if !student.BirthDate.IsZero() {
		birthDate = student.BirthDate.Format(time.RFC3339)
	}

	return &dto.StudentResponse{
		ID:               student.ID.String(),
		StudentNumber:    student.StudentNumber,
		FullName:         student.FullName,
		BirthDate:        birthDate,
		Gender:           student.Gender,
		ParentName:       student.ParentName,
		Status:           student.Status,
		IsActive:         student.IsActive,
		CreatedAt:        student.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        student.UpdatedAt.Format(time.RFC3339),
		DormitoryHistory: historyResponses,
	}
}
