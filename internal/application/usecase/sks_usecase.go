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

const (
	examDateLayout = "2006-01-02"
	examTimeLayout = "15:04"
)

// SKSDefinitionUseCase orchestrates SKS definition operations.
type SKSDefinitionUseCase struct {
	sksRepo     repository.SKSDefinitionRepository
	fanRepo     repository.FanRepository
	subjectRepo repository.SubjectRepository
	auditLogger appService.AuditLogger
}

// NewSKSDefinitionUseCase wires dependencies for SKS definitions.
func NewSKSDefinitionUseCase(
	sksRepo repository.SKSDefinitionRepository,
	fanRepo repository.FanRepository,
	subjectRepo repository.SubjectRepository,
	auditLogger appService.AuditLogger,
) *SKSDefinitionUseCase {
	return &SKSDefinitionUseCase{
		sksRepo:     sksRepo,
		fanRepo:     fanRepo,
		subjectRepo: subjectRepo,
		auditLogger: auditLogger,
	}
}

// CreateSKSDefinition creates a new SKS definition entry.
func (uc *SKSDefinitionUseCase) CreateSKSDefinition(ctx context.Context, req dto.CreateSKSDefinitionRequest) (*dto.SKSDefinitionResponse, error) {
	fanID, err := uuid.Parse(req.FanID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.fanRepo.GetByID(ctx, fanID); err != nil {
		return nil, domainErrors.ErrFanNotFound
	}

	var subjectID *uuid.UUID
	if req.SubjectID != nil && *req.SubjectID != "" {
		parsed, err := uuid.Parse(*req.SubjectID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if _, err := uc.subjectRepo.GetByID(ctx, parsed); err != nil {
			return nil, domainErrors.ErrSubjectNotFound
		}
		subjectID = &parsed
	}

	if existing, _ := uc.sksRepo.GetByCode(ctx, req.Code); existing != nil {
		return nil, domainErrors.ErrSKSDefinitionAlreadyExist
	}

	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	definition := &entity.SKSDefinition{
		ID:          uuid.New(),
		FanID:       fanID,
		SubjectID:   subjectID,
		Code:        req.Code,
		Name:        req.Name,
		KKM:         req.KKM,
		Description: req.Description,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.sksRepo.Create(ctx, definition); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "sks_definition", "sks_definition:create", definition.ID.String(), map[string]string{
		"code":   definition.Code,
		"fan_id": definition.FanID.String(),
	})

	return uc.toSKSDefinitionResponse(definition), nil
}

// GetSKSDefinition fetches SKS definition by ID.
func (uc *SKSDefinitionUseCase) GetSKSDefinition(ctx context.Context, id uuid.UUID) (*dto.SKSDefinitionResponse, error) {
	definition, err := uc.sksRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSKSDefinitionNotFound
	}
	return uc.toSKSDefinitionResponse(definition), nil
}

// ListSKSDefinitions lists definitions filtered by fan.
func (uc *SKSDefinitionUseCase) ListSKSDefinitions(ctx context.Context, fanIDStr string, page, pageSize int) (*dto.ListSKSDefinitionsResponse, error) {
	var fanID uuid.UUID
	if fanIDStr != "" {
		parsed, err := uuid.Parse(fanIDStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		fanID = parsed
	}

	page, pageSize = normalizePagination(page, pageSize)
	limit := pageSize
	offset := (page - 1) * pageSize

	definitions, total, err := uc.sksRepo.List(ctx, fanID, limit, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.SKSDefinitionResponse, 0, len(definitions))
	for _, definition := range definitions {
		responses = append(responses, *uc.toSKSDefinitionResponse(definition))
	}

	return &dto.ListSKSDefinitionsResponse{
		Definitions: responses,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  calcTotalPages(total, pageSize),
	}, nil
}

// UpdateSKSDefinition updates mutable fields of definition.
func (uc *SKSDefinitionUseCase) UpdateSKSDefinition(ctx context.Context, id uuid.UUID, req dto.UpdateSKSDefinitionRequest) (*dto.SKSDefinitionResponse, error) {
	definition, err := uc.sksRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSKSDefinitionNotFound
	}

	if req.SubjectID != nil {
		if *req.SubjectID == "" {
			definition.SubjectID = nil
		} else {
			parsed, err := uuid.Parse(*req.SubjectID)
			if err != nil {
				return nil, domainErrors.ErrBadRequest
			}
			if _, err := uc.subjectRepo.GetByID(ctx, parsed); err != nil {
				return nil, domainErrors.ErrSubjectNotFound
			}
			definition.SubjectID = &parsed
		}
	}

	if req.Name != nil {
		definition.Name = *req.Name
	}
	if req.KKM != nil {
		definition.KKM = *req.KKM
	}
	if req.Description != nil {
		definition.Description = *req.Description
	}
	if req.IsActive != nil {
		definition.IsActive = *req.IsActive
	}

	definition.UpdatedAt = time.Now()

	if err := uc.sksRepo.Update(ctx, definition); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "sks_definition", "sks_definition:update", definition.ID.String(), map[string]string{
		"code": definition.Code,
	})

	return uc.toSKSDefinitionResponse(definition), nil
}

// DeleteSKSDefinition deletes a definition by ID.
func (uc *SKSDefinitionUseCase) DeleteSKSDefinition(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.sksRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrSKSDefinitionNotFound
	}
	if err := uc.sksRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}
	_ = uc.auditLogger.Log(ctx, "sks_definition", "sks_definition:delete", id.String(), nil)
	return nil
}

func (uc *SKSDefinitionUseCase) toSKSDefinitionResponse(definition *entity.SKSDefinition) *dto.SKSDefinitionResponse {
	var subjectID *string
	if definition.SubjectID != nil {
		val := definition.SubjectID.String()
		subjectID = &val
	}

	return &dto.SKSDefinitionResponse{
		ID:          definition.ID.String(),
		FanID:       definition.FanID.String(),
		SubjectID:   subjectID,
		Code:        definition.Code,
		Name:        definition.Name,
		KKM:         definition.KKM,
		Description: definition.Description,
		IsActive:    definition.IsActive,
		CreatedAt:   definition.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   definition.UpdatedAt.Format(time.RFC3339),
	}
}

// SKSExamScheduleUseCase handles SKS exam schedule flows.
type SKSExamScheduleUseCase struct {
	examRepo    repository.SKSExamScheduleRepository
	sksRepo     repository.SKSDefinitionRepository
	teacherRepo repository.TeacherRepository
	auditLogger appService.AuditLogger
}

// NewSKSExamScheduleUseCase wires dependencies for exam schedules.
func NewSKSExamScheduleUseCase(
	examRepo repository.SKSExamScheduleRepository,
	sksRepo repository.SKSDefinitionRepository,
	teacherRepo repository.TeacherRepository,
	auditLogger appService.AuditLogger,
) *SKSExamScheduleUseCase {
	return &SKSExamScheduleUseCase{
		examRepo:    examRepo,
		sksRepo:     sksRepo,
		teacherRepo: teacherRepo,
		auditLogger: auditLogger,
	}
}

// CreateSKSExamSchedule creates an exam schedule for an SKS definition.
func (uc *SKSExamScheduleUseCase) CreateSKSExamSchedule(ctx context.Context, req dto.CreateSKSExamScheduleRequest) (*dto.SKSExamScheduleResponse, error) {
	sksID, err := uuid.Parse(req.SKSID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.sksRepo.GetByID(ctx, sksID); err != nil {
		return nil, domainErrors.ErrSKSDefinitionNotFound
	}

	var examinerID *uuid.UUID
	if req.ExaminerID != nil {
		if *req.ExaminerID != "" {
			parsed, err := uuid.Parse(*req.ExaminerID)
			if err != nil {
				return nil, domainErrors.ErrBadRequest
			}
			teacher, err := uc.teacherRepo.GetByID(ctx, parsed)
			if err != nil || teacher == nil {
				return nil, domainErrors.ErrTeacherNotFound
			}
			if !teacher.IsActive {
				return nil, domainErrors.ErrBadRequest
			}
			examinerID = &parsed
		}
	}

	examDate, err := time.Parse(examDateLayout, req.ExamDate)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	examTime, err := time.Parse(examTimeLayout, req.ExamTime)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	now := time.Now()
	exam := &entity.SKSExamSchedule{
		ID:         uuid.New(),
		SKSID:      sksID,
		ExaminerID: examinerID,
		ExamDate:   examDate,
		ExamTime:   examTime,
		Location:   req.Location,
		Notes:      req.Notes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.examRepo.Create(ctx, exam); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "sks_exam", "sks_exam:create", exam.ID.String(), map[string]string{
		"sks_id": exam.SKSID.String(),
	})

	return uc.toSKSExamScheduleResponse(exam), nil
}

// GetSKSExamSchedule fetches an exam schedule by ID.
func (uc *SKSExamScheduleUseCase) GetSKSExamSchedule(ctx context.Context, id uuid.UUID) (*dto.SKSExamScheduleResponse, error) {
	exam, err := uc.examRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSKSExamScheduleNotFound
	}
	return uc.toSKSExamScheduleResponse(exam), nil
}

// ListSKSExamSchedules lists exam schedules for an SKS.
func (uc *SKSExamScheduleUseCase) ListSKSExamSchedules(ctx context.Context, sksIDStr string, page, pageSize int) (*dto.ListSKSExamSchedulesResponse, error) {
	if sksIDStr == "" {
		return nil, domainErrors.ErrBadRequest
	}
	sksID, err := uuid.Parse(sksIDStr)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	page, pageSize = normalizePagination(page, pageSize)
	limit := pageSize
	offset := (page - 1) * pageSize

	exams, total, err := uc.examRepo.ListBySKS(ctx, sksID, limit, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.SKSExamScheduleResponse, 0, len(exams))
	for _, exam := range exams {
		responses = append(responses, *uc.toSKSExamScheduleResponse(exam))
	}

	return &dto.ListSKSExamSchedulesResponse{
		Exams:      responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateSKSExamSchedule updates mutable fields on exam schedule.
func (uc *SKSExamScheduleUseCase) UpdateSKSExamSchedule(ctx context.Context, id uuid.UUID, req dto.UpdateSKSExamScheduleRequest) (*dto.SKSExamScheduleResponse, error) {
	exam, err := uc.examRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrSKSExamScheduleNotFound
	}

	if req.ExaminerID != nil {
		if *req.ExaminerID == "" {
			exam.ExaminerID = nil
		} else {
			parsed, err := uuid.Parse(*req.ExaminerID)
			if err != nil {
				return nil, domainErrors.ErrBadRequest
			}
			teacher, err := uc.teacherRepo.GetByID(ctx, parsed)
			if err != nil || teacher == nil {
				return nil, domainErrors.ErrTeacherNotFound
			}
			if !teacher.IsActive {
				return nil, domainErrors.ErrBadRequest
			}
			exam.ExaminerID = &parsed
		}
	}

	if req.ExamDate != nil {
		parsed, err := time.Parse(examDateLayout, *req.ExamDate)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		exam.ExamDate = parsed
	}
	if req.ExamTime != nil {
		parsed, err := time.Parse(examTimeLayout, *req.ExamTime)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		exam.ExamTime = parsed
	}
	if req.Location != nil {
		exam.Location = *req.Location
	}
	if req.Notes != nil {
		exam.Notes = *req.Notes
	}

	exam.UpdatedAt = time.Now()

	if err := uc.examRepo.Update(ctx, exam); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "sks_exam", "sks_exam:update", exam.ID.String(), map[string]string{
		"sks_id": exam.SKSID.String(),
	})

	return uc.toSKSExamScheduleResponse(exam), nil
}

// DeleteSKSExamSchedule deletes an exam schedule.
func (uc *SKSExamScheduleUseCase) DeleteSKSExamSchedule(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.examRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrSKSExamScheduleNotFound
	}
	if err := uc.examRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}
	_ = uc.auditLogger.Log(ctx, "sks_exam", "sks_exam:delete", id.String(), nil)
	return nil
}

func (uc *SKSExamScheduleUseCase) toSKSExamScheduleResponse(exam *entity.SKSExamSchedule) *dto.SKSExamScheduleResponse {
	var examinerID *string
	if exam.ExaminerID != nil {
		val := exam.ExaminerID.String()
		examinerID = &val
	}

	return &dto.SKSExamScheduleResponse{
		ID:         exam.ID.String(),
		SKSID:      exam.SKSID.String(),
		ExaminerID: examinerID,
		ExamDate:   exam.ExamDate.Format(examDateLayout),
		ExamTime:   exam.ExamTime.Format(examTimeLayout),
		Location:   exam.Location,
		Notes:      exam.Notes,
		CreatedAt:  exam.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  exam.UpdatedAt.Format(time.RFC3339),
	}
}
