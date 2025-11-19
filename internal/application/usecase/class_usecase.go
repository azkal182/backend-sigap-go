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

// ClassUseCase orchestrates class operations.
type ClassUseCase struct {
	classRepo      repository.ClassRepository
	fanRepo        repository.FanRepository
	studentRepo    repository.StudentRepository
	enrollmentRepo repository.StudentClassEnrollmentRepository
	staffRepo      repository.ClassStaffRepository
	auditLogger    appService.AuditLogger
}

// NewClassUseCase builds a ClassUseCase instance.
func NewClassUseCase(
	classRepo repository.ClassRepository,
	fanRepo repository.FanRepository,
	studentRepo repository.StudentRepository,
	enrollmentRepo repository.StudentClassEnrollmentRepository,
	staffRepo repository.ClassStaffRepository,
	auditLogger appService.AuditLogger,
) *ClassUseCase {
	return &ClassUseCase{
		classRepo:      classRepo,
		fanRepo:        fanRepo,
		studentRepo:    studentRepo,
		enrollmentRepo: enrollmentRepo,
		staffRepo:      staffRepo,
		auditLogger:    auditLogger,
	}
}

// CreateClass creates a class under a FAN.
func (uc *ClassUseCase) CreateClass(ctx context.Context, req dto.CreateClassRequest) (*dto.ClassResponse, error) {
	fanID, err := uuid.Parse(req.FanID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	if _, err := uc.fanRepo.GetByID(ctx, fanID); err != nil {
		return nil, domainErrors.ErrFanNotFound
	}

	now := time.Now()
	class := &entity.Class{
		ID:        uuid.New(),
		FanID:     fanID,
		Name:      req.Name,
		Capacity:  req.Capacity,
		IsActive:  req.IsActive == nil || *req.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.classRepo.Create(ctx, class); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class", "class:create", class.ID.String(), map[string]string{
		"fan_id": fanID.String(),
		"name":   class.Name,
	})

	return uc.toClassResponse(class), nil
}

// GetClass retrieves class by ID.
func (uc *ClassUseCase) GetClass(ctx context.Context, id uuid.UUID) (*dto.ClassResponse, error) {
	class, err := uc.classRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrClassNotFound
	}
	return uc.toClassResponse(class), nil
}

// ListClassesByFan lists classes for a given FAN.
func (uc *ClassUseCase) ListClassesByFan(ctx context.Context, fanID uuid.UUID, page, pageSize int) (*dto.ListClassesResponse, error) {
	if _, err := uc.fanRepo.GetByID(ctx, fanID); err != nil {
		return nil, domainErrors.ErrFanNotFound
	}

	page, pageSize = normalizePagination(page, pageSize)
	classes, total, err := uc.classRepo.ListByFan(ctx, fanID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	resp := make([]dto.ClassResponse, 0, len(classes))
	for _, class := range classes {
		resp = append(resp, *uc.toClassResponse(class))
	}

	return &dto.ListClassesResponse{
		Classes:    resp,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateClass updates class attributes.
func (uc *ClassUseCase) UpdateClass(ctx context.Context, id uuid.UUID, req dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	class, err := uc.classRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrClassNotFound
	}

	if req.Name != nil {
		class.Name = *req.Name
	}
	if req.Capacity != nil {
		class.Capacity = *req.Capacity
	}
	if req.IsActive != nil {
		class.IsActive = *req.IsActive
	}
	class.UpdatedAt = time.Now()

	if err := uc.classRepo.Update(ctx, class); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class", "class:update", class.ID.String(), map[string]string{
		"name": class.Name,
	})

	return uc.toClassResponse(class), nil
}

// DeleteClass deletes a class by ID.
func (uc *ClassUseCase) DeleteClass(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.classRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrClassNotFound
	}

	if err := uc.classRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class", "class:delete", id.String(), nil)
	return nil
}

// EnrollStudent assigns a student to a class.
func (uc *ClassUseCase) EnrollStudent(ctx context.Context, classID uuid.UUID, req dto.EnrollStudentRequest) error {
	class, err := uc.classRepo.GetByID(ctx, classID)
	if err != nil {
		return domainErrors.ErrClassNotFound
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return domainErrors.ErrBadRequest
	}

	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return domainErrors.ErrStudentNotFound
	}

	// Ensure student not already active in class
	if enrollment, err := uc.enrollmentRepo.GetActiveByStudentAndClass(ctx, studentID, classID); err == nil && enrollment != nil {
		return domainErrors.ErrStudentAlreadyEnrolled
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return domainErrors.ErrBadRequest
	}
	now := time.Now()

	enrollment := &entity.StudentClassEnrollment{
		ID:         uuid.New(),
		ClassID:    classID,
		StudentID:  studentID,
		EnrolledAt: startDate,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class", "class:enroll-student", class.ID.String(), map[string]string{
		"student_id": studentID.String(),
	})
	return nil
}

// AssignStaff assigns staff to a class.
func (uc *ClassUseCase) AssignStaff(ctx context.Context, classID uuid.UUID, req dto.AssignClassStaffRequest) error {
	if _, err := uc.classRepo.GetByID(ctx, classID); err != nil {
		return domainErrors.ErrClassNotFound
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return domainErrors.ErrBadRequest
	}

	// We assume staff users exist in user repo; verifying optional.
	now := time.Now()
	staff := &entity.ClassStaff{
		ID:        uuid.New(),
		ClassID:   classID,
		UserID:    userID,
		Role:      req.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.staffRepo.Assign(ctx, staff); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class", "class:assign-staff", classID.String(), map[string]string{
		"user_id": userID.String(),
		"role":    req.Role,
	})
	return nil
}

func (uc *ClassUseCase) toClassResponse(class *entity.Class) *dto.ClassResponse {
	return &dto.ClassResponse{
		ID:        class.ID.String(),
		FanID:     class.FanID.String(),
		Name:      class.Name,
		Capacity:  class.Capacity,
		IsActive:  class.IsActive,
		CreatedAt: class.CreatedAt.Format(time.RFC3339),
		UpdatedAt: class.UpdatedAt.Format(time.RFC3339),
	}
}
