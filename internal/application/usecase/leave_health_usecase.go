package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

const isoDateLayout = "2006-01-02"

// LeavePermitUseCase orchestrates leave permit workflows.
type LeavePermitUseCase struct {
	leaveRepo   repository.LeavePermitRepository
	studentRepo repository.StudentRepository
	auditLogger appService.AuditLogger
}

// HealthStatusUseCase orchestrates student health status workflows.
type HealthStatusUseCase struct {
	healthRepo  repository.HealthStatusRepository
	studentRepo repository.StudentRepository
	auditLogger appService.AuditLogger
}

// NewLeavePermitUseCase builds a LeavePermitUseCase instance.
func NewLeavePermitUseCase(
	leaveRepo repository.LeavePermitRepository,
	studentRepo repository.StudentRepository,
	auditLogger appService.AuditLogger,
) *LeavePermitUseCase {
	return &LeavePermitUseCase{leaveRepo: leaveRepo, studentRepo: studentRepo, auditLogger: auditLogger}
}

// NewHealthStatusUseCase builds a HealthStatusUseCase instance.
func NewHealthStatusUseCase(
	healthRepo repository.HealthStatusRepository,
	studentRepo repository.StudentRepository,
	auditLogger appService.AuditLogger,
) *HealthStatusUseCase {
	return &HealthStatusUseCase{healthRepo: healthRepo, studentRepo: studentRepo, auditLogger: auditLogger}
}

// CreateLeavePermit registers a new leave permit in pending status.
func (uc *LeavePermitUseCase) CreateLeavePermit(ctx context.Context, req dto.CreateLeavePermitRequest) (*dto.LeavePermitResponse, error) {
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	leaveType, err := parseLeavePermitType(req.Type)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	startDate, endDate, err := parseDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	actorID, err := requireActorID(ctx)
	if err != nil {
		return nil, err
	}

	overlap, err := uc.leaveRepo.HasOverlap(ctx, studentID, startDate, endDate, nil)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}
	if overlap {
		return nil, domainErrors.ErrLeavePermitConflict
	}

	now := time.Now()
	permit := &entity.LeavePermit{
		ID:        uuid.New(),
		StudentID: studentID,
		Type:      leaveType,
		Reason:    req.Reason,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    entity.LeavePermitStatusPending,
		CreatedBy: actorID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.leaveRepo.Create(ctx, permit); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "leave_permit", "leave_permit:create", permit.ID.String(), map[string]string{
		"student_id": studentID.String(),
		"type":       string(permit.Type),
	})

	resp := toLeavePermitResponse(permit)
	return &resp, nil
}

// ListLeavePermits returns paginated leave permits.
func (uc *LeavePermitUseCase) ListLeavePermits(ctx context.Context, req dto.ListLeavePermitsRequest) (*dto.ListLeavePermitsResponse, error) {
	filter := repository.LeavePermitFilter{}

	if req.StudentID != nil && *req.StudentID != "" {
		studentID, err := uuid.Parse(*req.StudentID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.StudentID = &studentID
	}

	if req.Status != nil && *req.Status != "" {
		status, err := parseLeavePermitStatus(*req.Status)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Status = &status
	}

	if req.Type != nil && *req.Type != "" {
		permitType, err := parseLeavePermitType(*req.Type)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Type = &permitType
	}

	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse(isoDateLayout, *req.Date)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Date = &parsed
	}

	page, pageSize := normalizePagination(req.Page, req.PageSize)
	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	permits, total, err := uc.leaveRepo.List(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.LeavePermitResponse, 0, len(permits))
	for _, permit := range permits {
		resp := toLeavePermitResponse(permit)
		responses = append(responses, resp)
	}

	return &dto.ListLeavePermitsResponse{
		Permits:    responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateLeavePermitStatus processes workflow transitions (approve/reject/complete).
func (uc *LeavePermitUseCase) UpdateLeavePermitStatus(ctx context.Context, permitID uuid.UUID, req dto.UpdateLeavePermitStatusRequest) (*dto.LeavePermitResponse, error) {
	permit, err := uc.leaveRepo.GetByID(ctx, permitID)
	if err != nil {
		return nil, domainErrors.ErrLeavePermitNotFound
	}

	newStatus, err := parseLeavePermitStatus(req.Status)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	actorID, err := requireActorID(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	switch newStatus {
	case entity.LeavePermitStatusApproved:
		if permit.Status != entity.LeavePermitStatusPending {
			return nil, domainErrors.ErrLeavePermitStatus
		}
		permit.Status = entity.LeavePermitStatusApproved
		permit.ApprovedBy = &actorID
		permit.ApprovedAt = &now
	case entity.LeavePermitStatusRejected:
		if permit.Status != entity.LeavePermitStatusPending {
			return nil, domainErrors.ErrLeavePermitStatus
		}
		permit.Status = entity.LeavePermitStatusRejected
		permit.ApprovedBy = &actorID
		permit.ApprovedAt = &now
	case entity.LeavePermitStatusCompleted:
		if permit.Status != entity.LeavePermitStatusApproved {
			return nil, domainErrors.ErrLeavePermitStatus
		}
		permit.Status = entity.LeavePermitStatusCompleted
	default:
		return nil, domainErrors.ErrLeavePermitStatus
	}

	permit.UpdatedAt = now

	if err := uc.leaveRepo.Update(ctx, permit); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "leave_permit", "leave_permit:update_status", permit.ID.String(), map[string]string{
		"status": req.Status,
	})

	resp := toLeavePermitResponse(permit)
	return &resp, nil
}

// GetActivePermitForDate returns active permit overlapping a date (attendance hook helper).
func (uc *LeavePermitUseCase) GetActivePermitForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error) {
	permit, err := uc.leaveRepo.ActiveByDate(ctx, studentID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, domainErrors.ErrInternalServer
	}
	return permit, nil
}

// CreateHealthStatus registers a new active health status.
func (uc *HealthStatusUseCase) CreateHealthStatus(ctx context.Context, req dto.CreateHealthStatusRequest) (*dto.HealthStatusResponse, error) {
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	startDate, err := time.Parse(isoDateLayout, req.StartDate)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	var endDatePtr *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse(isoDateLayout, *req.EndDate)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if parsed.Before(startDate) {
			return nil, domainErrors.ErrBadRequest
		}
		endDatePtr = &parsed
	}

	if _, err := uc.studentRepo.GetByID(ctx, studentID); err != nil {
		return nil, domainErrors.ErrStudentNotFound
	}

	actorID, err := requireActorID(ctx)
	if err != nil {
		return nil, err
	}

	if existing, err := uc.healthRepo.ActiveByDate(ctx, studentID, startDate); err == nil && existing != nil {
		return nil, domainErrors.ErrHealthStatusActive
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainErrors.ErrInternalServer
	}

	now := time.Now()
	status := &entity.HealthStatus{
		ID:        uuid.New(),
		StudentID: studentID,
		Diagnosis: req.Diagnosis,
		Notes:     req.Notes,
		StartDate: startDate,
		EndDate:   endDatePtr,
		Status:    entity.HealthStatusStateActive,
		CreatedBy: actorID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.healthRepo.Create(ctx, status); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "health_status", "health_status:create", status.ID.String(), map[string]string{
		"student_id": studentID.String(),
	})

	resp := toHealthStatusResponse(status)
	return &resp, nil
}

// ListHealthStatuses returns paginated health statuses.
func (uc *HealthStatusUseCase) ListHealthStatuses(ctx context.Context, req dto.ListHealthStatusesRequest) (*dto.ListHealthStatusesResponse, error) {
	filter := repository.HealthStatusFilter{}

	if req.StudentID != nil && *req.StudentID != "" {
		studentID, err := uuid.Parse(*req.StudentID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.StudentID = &studentID
	}

	if req.Status != nil && *req.Status != "" {
		status, err := parseHealthStatusState(*req.Status)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Status = &status
	}

	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse(isoDateLayout, *req.Date)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Date = &parsed
	}

	page, pageSize := normalizePagination(req.Page, req.PageSize)
	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	statuses, total, err := uc.healthRepo.List(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.HealthStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		resp := toHealthStatusResponse(status)
		responses = append(responses, resp)
	}

	return &dto.ListHealthStatusesResponse{
		Statuses:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// RevokeHealthStatus marks an active health status as revoked.
func (uc *HealthStatusUseCase) RevokeHealthStatus(ctx context.Context, statusID uuid.UUID, req dto.RevokeHealthStatusRequest) (*dto.HealthStatusResponse, error) {
	status, err := uc.healthRepo.GetByID(ctx, statusID)
	if err != nil {
		return nil, domainErrors.ErrHealthStatusNotFound
	}

	if status.Status != entity.HealthStatusStateActive {
		return nil, domainErrors.ErrHealthStatusForbidden
	}

	actorID, err := requireActorID(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	status.Status = entity.HealthStatusStateRevoked
	status.RevokedBy = &actorID
	status.RevokedAt = &now
	status.UpdatedAt = now

	if err := uc.healthRepo.Update(ctx, status); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	metadata := map[string]string{}
	if req.Reason != "" {
		metadata["reason"] = req.Reason
	}
	_ = uc.auditLogger.Log(ctx, "health_status", "health_status:revoke", status.ID.String(), metadata)

	resp := toHealthStatusResponse(status)
	return &resp, nil
}

// GetActiveHealthStatusForDate returns an active health status overlapping a date.
func (uc *HealthStatusUseCase) GetActiveHealthStatusForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error) {
	status, err := uc.healthRepo.ActiveByDate(ctx, studentID, date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, domainErrors.ErrInternalServer
	}
	return status, nil
}

func parseDateRange(start, end string) (time.Time, time.Time, error) {
	startDate, err := time.Parse(isoDateLayout, start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endDateParsed, err := time.Parse(isoDateLayout, end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if endDateParsed.Before(startDate) {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	return startDate, endDateParsed, nil
}

func parseLeavePermitStatus(status string) (entity.LeavePermitStatus, error) {
	switch entity.LeavePermitStatus(status) {
	case entity.LeavePermitStatusPending,
		entity.LeavePermitStatusApproved,
		entity.LeavePermitStatusRejected,
		entity.LeavePermitStatusCompleted:
		return entity.LeavePermitStatus(status), nil
	default:
		return "", errors.New("invalid leave permit status")
	}
}

func parseLeavePermitType(t string) (entity.LeavePermitType, error) {
	switch entity.LeavePermitType(t) {
	case entity.LeavePermitTypeHomeLeave, entity.LeavePermitTypeOfficialDuty:
		return entity.LeavePermitType(t), nil
	default:
		return "", errors.New("invalid leave permit type")
	}
}

func parseHealthStatusState(s string) (entity.HealthStatusState, error) {
	switch entity.HealthStatusState(s) {
	case entity.HealthStatusStateActive, entity.HealthStatusStateRevoked:
		return entity.HealthStatusState(s), nil
	default:
		return "", errors.New("invalid health status state")
	}
}

func toLeavePermitResponse(permit *entity.LeavePermit) dto.LeavePermitResponse {
	resp := dto.LeavePermitResponse{
		ID:        permit.ID.String(),
		StudentID: permit.StudentID.String(),
		Type:      string(permit.Type),
		Reason:    permit.Reason,
		StartDate: permit.StartDate.Format(isoDateLayout),
		EndDate:   permit.EndDate.Format(isoDateLayout),
		Status:    string(permit.Status),
		CreatedBy: permit.CreatedBy.String(),
	}

	resp.ApprovedBy = dto.OptionalUUIDToString(permit.ApprovedBy)
	if permit.ApprovedAt != nil {
		formatted := permit.ApprovedAt.Format(time.RFC3339)
		resp.ApprovedAt = &formatted
	}

	return resp
}

func toHealthStatusResponse(status *entity.HealthStatus) dto.HealthStatusResponse {
	resp := dto.HealthStatusResponse{
		ID:        status.ID.String(),
		StudentID: status.StudentID.String(),
		Diagnosis: status.Diagnosis,
		Notes:     status.Notes,
		StartDate: status.StartDate.Format(isoDateLayout),
		Status:    string(status.Status),
		CreatedBy: status.CreatedBy.String(),
	}

	if status.EndDate != nil {
		formatted := status.EndDate.Format(isoDateLayout)
		resp.EndDate = &formatted
	}

	resp.RevokedBy = dto.OptionalUUIDToString(status.RevokedBy)
	if status.RevokedAt != nil {
		formatted := status.RevokedAt.Format(time.RFC3339)
		resp.RevokedAt = &formatted
	}

	return resp
}
