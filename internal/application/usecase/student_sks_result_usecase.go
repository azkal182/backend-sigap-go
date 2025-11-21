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

const studentExamDateLayout = "2006-01-02"

// StudentSKSResultUseCase orchestrates SKS exam result flows and FAN completion tracking.
type StudentSKSResultUseCase struct {
	resultRepo    repository.StudentSKSResultRepository
	fanStatusRepo repository.FanCompletionStatusRepository
	studentRepo   repository.StudentRepository
	sksRepo       repository.SKSDefinitionRepository
	teacherRepo   repository.TeacherRepository
	auditLogger   appService.AuditLogger
}

// NewStudentSKSResultUseCase wires dependencies for the SKS result use case.
func NewStudentSKSResultUseCase(
	resultRepo repository.StudentSKSResultRepository,
	fanStatusRepo repository.FanCompletionStatusRepository,
	studentRepo repository.StudentRepository,
	sksRepo repository.SKSDefinitionRepository,
	teacherRepo repository.TeacherRepository,
	auditLogger appService.AuditLogger,
) *StudentSKSResultUseCase {
	return &StudentSKSResultUseCase{
		resultRepo:    resultRepo,
		fanStatusRepo: fanStatusRepo,
		studentRepo:   studentRepo,
		sksRepo:       sksRepo,
		teacherRepo:   teacherRepo,
		auditLogger:   auditLogger,
	}
}

// CreateStudentSKSResult records a student's SKS outcome.
func (uc *StudentSKSResultUseCase) CreateStudentSKSResult(ctx context.Context, req dto.CreateStudentSKSResultRequest) (*dto.StudentSKSResultResponse, error) {
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	sksID, err := uuid.Parse(req.SKSID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	definition, err := uc.sksRepo.GetByID(ctx, sksID)
	if err != nil {
		return nil, domainErrors.ErrSKSDefinitionNotFound
	}

	examinerID, err := uc.validateExaminer(ctx, req.ExaminerID)
	if err != nil {
		return nil, err
	}

	examDate, err := parseOptionalDate(req.ExamDate)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	isPassed := uc.resolvePassFlag(req.IsPassed, req.Score, definition.KKM)
	now := time.Now()
	result := &entity.StudentSKSResult{
		ID:         uuid.New(),
		StudentID:  studentID,
		SKSID:      sksID,
		Score:      req.Score,
		IsPassed:   isPassed,
		ExamDate:   examDate,
		ExaminerID: examinerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.resultRepo.Create(ctx, result); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	uc.updateFanCompletion(ctx, studentID, definition.FanID)
	uc.logAudit(ctx, "sks_result:create", result.ID, map[string]string{
		"student_id": result.StudentID.String(),
		"sks_id":     result.SKSID.String(),
	})

	return uc.toStudentSKSResultResponse(result, definition.FanID)
}

// UpdateStudentSKSResult updates an existing SKS result.
func (uc *StudentSKSResultUseCase) UpdateStudentSKSResult(ctx context.Context, id uuid.UUID, req dto.UpdateStudentSKSResultRequest) (*dto.StudentSKSResultResponse, error) {
	result, err := uc.resultRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrStudentSKSResultNotFound
	}

	definition, err := uc.sksRepo.GetByID(ctx, result.SKSID)
	if err != nil {
		return nil, domainErrors.ErrSKSDefinitionNotFound
	}

	if req.Score != nil {
		result.Score = *req.Score
	}
	if req.IsPassed != nil {
		result.IsPassed = *req.IsPassed
	} else if req.Score != nil {
		result.IsPassed = *req.Score >= definition.KKM
	}

	examinerID, err := uc.validateExaminer(ctx, req.ExaminerID)
	if err != nil {
		return nil, err
	}
	if req.ExaminerID != nil {
		result.ExaminerID = examinerID
	}

	if req.ExamDate != nil {
		parsedDate, err := parseOptionalDate(req.ExamDate)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		result.ExamDate = parsedDate
	}

	result.UpdatedAt = time.Now()
	if err := uc.resultRepo.Update(ctx, result); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	uc.updateFanCompletion(ctx, result.StudentID, definition.FanID)
	uc.logAudit(ctx, "sks_result:update", result.ID, map[string]string{
		"student_id": result.StudentID.String(),
		"sks_id":     result.SKSID.String(),
	})

	return uc.toStudentSKSResultResponse(result, definition.FanID)
}

// ListStudentSKSResults lists SKS results for a student (optionally filtered by FAN).
func (uc *StudentSKSResultUseCase) ListStudentSKSResults(ctx context.Context, studentID uuid.UUID, fanIDStr string, page, pageSize int) (*dto.ListStudentSKSResultsResponse, error) {
	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	var fanID uuid.UUID
	if fanIDStr != "" {
		parsed, err := uuid.Parse(fanIDStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		fanID = parsed
	}

	page, pageSize = normalizePagination(page, pageSize)
	results, total, err := uc.resultRepo.ListByStudent(ctx, studentID, fanID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.StudentSKSResultResponse, 0, len(results))
	fanCache := make(map[uuid.UUID]uuid.UUID)
	for _, res := range results {
		resolvedFanID, err := uc.getFanID(ctx, fanCache, res.SKSID)
		if err != nil {
			return nil, domainErrors.ErrInternalServer
		}
		resp, err := uc.toStudentSKSResultResponse(res, resolvedFanID)
		if err != nil {
			return nil, domainErrors.ErrInternalServer
		}
		responses = append(responses, *resp)
	}

	return &dto.ListStudentSKSResultsResponse{
		Results:    responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// ListFanCompletionStatuses returns FAN completion summary for a student.
func (uc *StudentSKSResultUseCase) ListFanCompletionStatuses(ctx context.Context, studentID uuid.UUID) ([]dto.FanCompletionStatusResponse, error) {
	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}
	statuses, err := uc.fanStatusRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.FanCompletionStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		responses = append(responses, dto.FanCompletionStatusResponse{
			FanID:       status.FanID.String(),
			IsCompleted: status.IsCompleted,
			CompletedAt: formatTimePtr(status.CompletedAt),
		})
	}
	return responses, nil
}

func (uc *StudentSKSResultUseCase) getFanID(ctx context.Context, cache map[uuid.UUID]uuid.UUID, sksID uuid.UUID) (uuid.UUID, error) {
	if fanID, ok := cache[sksID]; ok {
		return fanID, nil
	}
	definition, err := uc.sksRepo.GetByID(ctx, sksID)
	if err != nil {
		return uuid.Nil, err
	}
	cache[sksID] = definition.FanID
	return definition.FanID, nil
}

func (uc *StudentSKSResultUseCase) updateFanCompletion(ctx context.Context, studentID, fanID uuid.UUID) {
	if fanID == uuid.Nil {
		return
	}
	totalSKS, err := uc.sksRepo.CountByFan(ctx, fanID)
	if err != nil || totalSKS == 0 {
		return
	}
	passed, err := uc.resultRepo.CountPassedByStudentFan(ctx, studentID, fanID)
	if err != nil {
		return
	}
	isCompleted := passed >= totalSKS
	now := time.Now()
	status := &entity.FanCompletionStatus{
		ID:          uuid.New(),
		StudentID:   studentID,
		FanID:       fanID,
		IsCompleted: isCompleted,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if isCompleted {
		status.CompletedAt = &now
	}
	_ = uc.fanStatusRepo.Upsert(ctx, status)
}

func (uc *StudentSKSResultUseCase) validateExaminer(ctx context.Context, examinerIDStr *string) (*uuid.UUID, error) {
	if examinerIDStr == nil || *examinerIDStr == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(*examinerIDStr)
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
	return &parsed, nil
}

func (uc *StudentSKSResultUseCase) resolvePassFlag(explicit *bool, score, kkm float64) bool {
	if explicit != nil {
		return *explicit
	}
	return score >= kkm
}

func parseOptionalDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}
	parsed, err := time.Parse(studentExamDateLayout, *dateStr)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (uc *StudentSKSResultUseCase) toStudentSKSResultResponse(result *entity.StudentSKSResult, fanID uuid.UUID) (*dto.StudentSKSResultResponse, error) {
	return &dto.StudentSKSResultResponse{
		ID:         result.ID.String(),
		StudentID:  result.StudentID.String(),
		FanID:      fanID.String(),
		SKSID:      result.SKSID.String(),
		Score:      result.Score,
		IsPassed:   result.IsPassed,
		ExamDate:   formatTimePtr(result.ExamDate),
		ExaminerID: uuidPtrToString(result.ExaminerID),
		CreatedAt:  result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  result.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *StudentSKSResultUseCase) logAudit(ctx context.Context, action string, id uuid.UUID, metadata map[string]string) {
	_ = uc.auditLogger.Log(ctx, "sks_result", action, id.String(), metadata)
}

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	val := id.String()
	return &val
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	val := t.Format(time.RFC3339)
	return &val
}
