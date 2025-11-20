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

// SubjectUseCase orchestrates subject operations.
type SubjectUseCase struct {
	subjectRepo repository.SubjectRepository
	auditLogger appService.AuditLogger
}

// NewSubjectUseCase creates a new SubjectUseCase instance.
func NewSubjectUseCase(subjectRepo repository.SubjectRepository, auditLogger appService.AuditLogger) *SubjectUseCase {
	return &SubjectUseCase{subjectRepo: subjectRepo, auditLogger: auditLogger}
}

// CreateSubject creates a new subject entry.
func (uc *SubjectUseCase) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	now := time.Now()
	subject := &entity.Subject{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive == nil || *req.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.subjectRepo.Create(ctx, subject); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "subject", "subject:create", subject.ID.String(), map[string]string{
		"name": subject.Name,
	})

	return uc.toSubjectResponse(subject), nil
}

// GetSubject retrieves a subject by ID.
func (uc *SubjectUseCase) GetSubject(ctx context.Context, id uuid.UUID) (*dto.SubjectResponse, error) {
	subject, err := uc.subjectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSubjectNotFound
	}
	return uc.toSubjectResponse(subject), nil
}

// ListSubjects returns paginated subject data.
func (uc *SubjectUseCase) ListSubjects(ctx context.Context, page, pageSize int) (*dto.ListSubjectsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	subjects, total, err := uc.subjectRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.SubjectResponse, 0, len(subjects))
	for _, subject := range subjects {
		responses = append(responses, *uc.toSubjectResponse(subject))
	}

	return &dto.ListSubjectsResponse{
		Subjects:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateSubject updates an existing subject.
func (uc *SubjectUseCase) UpdateSubject(ctx context.Context, id uuid.UUID, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	subject, err := uc.subjectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSubjectNotFound
	}

	if req.Name != nil {
		subject.Name = *req.Name
	}
	if req.Description != nil {
		subject.Description = *req.Description
	}
	if req.IsActive != nil {
		subject.IsActive = *req.IsActive
	}
	subject.UpdatedAt = time.Now()

	if err := uc.subjectRepo.Update(ctx, subject); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "subject", "subject:update", subject.ID.String(), map[string]string{
		"name": subject.Name,
	})

	return uc.toSubjectResponse(subject), nil
}

// DeleteSubject deletes a subject.
func (uc *SubjectUseCase) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.subjectRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrSubjectNotFound
	}

	if err := uc.subjectRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "subject", "subject:delete", id.String(), nil)
	return nil
}

func (uc *SubjectUseCase) toSubjectResponse(subject *entity.Subject) *dto.SubjectResponse {
	return &dto.SubjectResponse{
		ID:          subject.ID.String(),
		Name:        subject.Name,
		Description: subject.Description,
		IsActive:    subject.IsActive,
		CreatedAt:   subject.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   subject.UpdatedAt.Format(time.RFC3339),
	}
}
