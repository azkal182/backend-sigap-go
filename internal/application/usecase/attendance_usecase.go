package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

type leavePermitStatusProvider interface {
	GetActivePermitForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error)
}

func (uc *AttendanceUseCase) getDerivedStatus(
	ctx context.Context,
	studentID uuid.UUID,
	date time.Time,
	cache map[uuid.UUID]derivedStatusCacheEntry,
) (derivedStatusCacheEntry, error) {
	if entry, ok := cache[studentID]; ok {
		return entry, nil
	}

	result := derivedStatusCacheEntry{}

	if uc.healthStatusProvider != nil {
		record, err := uc.healthStatusProvider.GetActiveHealthStatusForDate(ctx, studentID, date)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, domainErrors.ErrInternalServer
		}
		if record != nil {
			result.override = true
			result.status = entity.StudentAttendanceSick
			cache[studentID] = result
			return result, nil
		}
	}

	if uc.leavePermitProvider != nil {
		record, err := uc.leavePermitProvider.GetActivePermitForDate(ctx, studentID, date)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, domainErrors.ErrInternalServer
		}
		if record != nil {
			result.override = true
			result.status = entity.StudentAttendancePermit
			cache[studentID] = result
			return result, nil
		}
	}

	cache[studentID] = result
	return result, nil
}

type healthStatusProvider interface {
	GetActiveHealthStatusForDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error)
}

type derivedStatusCacheEntry struct {
	override bool
	status   entity.StudentAttendanceStatus
}

// AttendanceUseCase orchestrates attendance operations.
type AttendanceUseCase struct {
	sessionRepo           repository.AttendanceSessionRepository
	studentAttendanceRepo repository.StudentAttendanceRepository
	teacherAttendanceRepo repository.TeacherAttendanceRepository
	classScheduleRepo     repository.ClassScheduleRepository
	leavePermitProvider   leavePermitStatusProvider
	healthStatusProvider  healthStatusProvider
	auditLogger           appService.AuditLogger
}

// NewAttendanceUseCase builds AttendanceUseCase instance.
func NewAttendanceUseCase(
	sessionRepo repository.AttendanceSessionRepository,
	studentAttendanceRepo repository.StudentAttendanceRepository,
	teacherAttendanceRepo repository.TeacherAttendanceRepository,
	classScheduleRepo repository.ClassScheduleRepository,
	leavePermitProvider leavePermitStatusProvider,
	healthStatusProvider healthStatusProvider,
	auditLogger appService.AuditLogger,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		sessionRepo:           sessionRepo,
		studentAttendanceRepo: studentAttendanceRepo,
		teacherAttendanceRepo: teacherAttendanceRepo,
		classScheduleRepo:     classScheduleRepo,
		leavePermitProvider:   leavePermitProvider,
		healthStatusProvider:  healthStatusProvider,
		auditLogger:           auditLogger,
	}
}

// OpenSessions opens attendance sessions for the provided schedules on a given date.
func (uc *AttendanceUseCase) OpenSessions(ctx context.Context, req dto.OpenAttendanceSessionRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return domainErrors.ErrBadRequest
	}

	for _, scheduleIDStr := range req.ClassScheduleIDs {
		classScheduleID, err := uuid.Parse(scheduleIDStr)
		if err != nil {
			return domainErrors.ErrBadRequest
		}

		schedule, err := uc.classScheduleRepo.GetByID(ctx, classScheduleID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return domainErrors.ErrClassScheduleNotFound
			}
			return domainErrors.ErrInternalServer
		}

		existing, err := uc.sessionRepo.GetOpenByScheduleAndDate(ctx, classScheduleID, date)
		if err == nil && existing != nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return domainErrors.ErrInternalServer
		}

		now := time.Now()
		session := &entity.AttendanceSession{
			ID:              uuid.New(),
			ClassScheduleID: classScheduleID,
			Date:            date,
			StartTime:       schedule.StartTime,
			EndTime:         schedule.EndTime,
			TeacherID:       schedule.TeacherID,
			Status:          entity.AttendanceSessionStatusOpen,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := uc.sessionRepo.Create(ctx, session); err != nil {
			return domainErrors.ErrInternalServer
		}

		_ = uc.auditLogger.Log(ctx, "attendance_session", "attendance_session:open", session.ID.String(), map[string]string{
			"class_schedule_id": classScheduleID.String(),
			"date":              date.Format("2006-01-02"),
		})
	}

	return nil
}

// SubmitStudentAttendance bulk-submits student attendance for a session.
func (uc *AttendanceUseCase) SubmitStudentAttendance(ctx context.Context, sessionID uuid.UUID, req dto.SubmitStudentAttendanceRequest) error {
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domainErrors.ErrAttendanceSessionNotFound
		}
		return domainErrors.ErrInternalServer
	}

	if session.Status == entity.AttendanceSessionStatusLocked {
		return domainErrors.ErrAttendanceAlreadyLocked
	}
	if len(req.Records) == 0 {
		return domainErrors.ErrBadRequest
	}

	now := time.Now()
	attendances := make([]*entity.StudentAttendance, 0, len(req.Records))
	derivedCache := make(map[uuid.UUID]derivedStatusCacheEntry)
	for _, record := range req.Records {
		studentID, err := uuid.Parse(record.StudentID)
		if err != nil {
			return domainErrors.ErrBadRequest
		}
		status, err := mapStudentStatus(record.Status)
		if err != nil {
			return err
		}

		if uc.leavePermitProvider != nil || uc.healthStatusProvider != nil {
			entry, err := uc.getDerivedStatus(ctx, studentID, session.Date, derivedCache)
			if err != nil {
				return err
			}
			if entry.override {
				status = entry.status
			}
		}
		attendances = append(attendances, &entity.StudentAttendance{
			ID:                  uuid.New(),
			AttendanceSessionID: sessionID,
			StudentID:           studentID,
			Status:              status,
			Note:                record.Note,
			CreatedAt:           now,
			UpdatedAt:           now,
		})
	}

	if err := uc.studentAttendanceRepo.BulkUpsert(ctx, attendances); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "attendance", "attendance:students:update", sessionID.String(), map[string]string{
		"count": fmt.Sprintf("%d", len(attendances)),
	})
	return nil
}

// SubmitTeacherAttendance submits teacher attendance for a session.
func (uc *AttendanceUseCase) SubmitTeacherAttendance(ctx context.Context, sessionID uuid.UUID, req dto.SubmitTeacherAttendanceRequest) error {
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domainErrors.ErrAttendanceSessionNotFound
		}
		return domainErrors.ErrInternalServer
	}

	if session.Status == entity.AttendanceSessionStatusLocked {
		return domainErrors.ErrAttendanceAlreadyLocked
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return domainErrors.ErrBadRequest
	}

	status, err := mapTeacherStatus(req.Status)
	if err != nil {
		return err
	}

	record := &entity.TeacherAttendance{
		ID:                  uuid.New(),
		AttendanceSessionID: sessionID,
		TeacherID:           teacherID,
		Status:              status,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := uc.teacherAttendanceRepo.Upsert(ctx, record); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "attendance", "attendance:teacher:update", sessionID.String(), map[string]string{
		"teacher_id": req.TeacherID,
	})
	return nil
}

// LockSessions locks all sessions for a particular date.
func (uc *AttendanceUseCase) LockSessions(ctx context.Context, req dto.LockAttendanceRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return domainErrors.ErrBadRequest
	}

	if err := uc.sessionRepo.LockSessionsByDate(ctx, date); err != nil {
		return domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "attendance", "attendance:lock", req.Date, nil)
	return nil
}

// ListAttendanceSessions lists sessions by filters.
func (uc *AttendanceUseCase) ListAttendanceSessions(ctx context.Context, req dto.ListAttendanceSessionsRequest) (*dto.ListAttendanceSessionsResponse, error) {
	filter := repository.AttendanceSessionFilter{}
	if req.ClassScheduleID != nil && *req.ClassScheduleID != "" {
		parsed, err := uuid.Parse(*req.ClassScheduleID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.ClassScheduleID = &parsed
	}
	if req.TeacherID != nil && *req.TeacherID != "" {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.TeacherID = &parsed
	}
	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.Date = &parsed
	}
	if req.Status != nil && *req.Status != "" {
		status := entity.AttendanceSessionStatus(*req.Status)
		filter.Status = &status
	}

	page, pageSize := normalizePagination(req.Page, req.PageSize)
	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	sessions, total, err := uc.sessionRepo.List(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.AttendanceSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		resp := dto.AttendanceSessionResponse{
			ID:              session.ID.String(),
			ClassScheduleID: session.ClassScheduleID.String(),
			Date:            session.Date.Format("2006-01-02"),
			StartTime:       formatTimePtr(session.StartTime),
			EndTime:         formatTimePtr(session.EndTime),
			TeacherID:       session.TeacherID.String(),
			Status:          string(session.Status),
			LockedAt:        formatTimePtr(session.LockedAt),
		}

		studentRecords := make([]dto.StudentAttendanceRecordResponse, 0, len(session.StudentAttendances))
		for _, record := range session.StudentAttendances {
			studentRecords = append(studentRecords, dto.StudentAttendanceRecordResponse{
				StudentID: record.StudentID.String(),
				Status:    string(record.Status),
				Note:      record.Note,
			})
		}
		resp.StudentRecords = studentRecords

		if len(session.TeacherAttendances) > 0 {
			record := session.TeacherAttendances[0]
			resp.TeacherRecord = &dto.TeacherAttendanceRecordResponse{
				TeacherID: record.TeacherID.String(),
				Status:    string(record.Status),
			}
		}

		responses = append(responses, resp)
	}

	return &dto.ListAttendanceSessionsResponse{
		Sessions:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

func mapStudentStatus(status string) (entity.StudentAttendanceStatus, error) {
	switch status {
	case string(entity.StudentAttendancePresent):
		return entity.StudentAttendancePresent, nil
	case string(entity.StudentAttendanceAbsent):
		return entity.StudentAttendanceAbsent, nil
	case string(entity.StudentAttendancePermit):
		return entity.StudentAttendancePermit, nil
	case string(entity.StudentAttendanceSick):
		return entity.StudentAttendanceSick, nil
	default:
		return "", domainErrors.ErrAttendanceInvalidStatus
	}
}

func mapTeacherStatus(status string) (entity.TeacherAttendanceStatus, error) {
	switch status {
	case string(entity.TeacherAttendancePresent):
		return entity.TeacherAttendancePresent, nil
	case string(entity.TeacherAttendanceAbsent):
		return entity.TeacherAttendanceAbsent, nil
	default:
		return "", domainErrors.ErrAttendanceInvalidStatus
	}
}
