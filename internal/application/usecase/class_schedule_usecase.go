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

// ClassScheduleUseCase orchestrates class schedule flows.
type ClassScheduleUseCase struct {
	scheduleRepo repository.ClassScheduleRepository
	classRepo    repository.ClassRepository
	teacherRepo  repository.TeacherRepository
	subjectRepo  repository.SubjectRepository
	slotRepo     repository.ScheduleSlotRepository
	dormRepo     repository.DormitoryRepository
	auditLogger  appService.AuditLogger
}

// NewClassScheduleUseCase wires dependencies for class schedules.
func NewClassScheduleUseCase(
	scheduleRepo repository.ClassScheduleRepository,
	classRepo repository.ClassRepository,
	teacherRepo repository.TeacherRepository,
	subjectRepo repository.SubjectRepository,
	slotRepo repository.ScheduleSlotRepository,
	dormRepo repository.DormitoryRepository,
	auditLogger appService.AuditLogger,
) *ClassScheduleUseCase {
	return &ClassScheduleUseCase{
		scheduleRepo: scheduleRepo,
		classRepo:    classRepo,
		teacherRepo:  teacherRepo,
		subjectRepo:  subjectRepo,
		slotRepo:     slotRepo,
		dormRepo:     dormRepo,
		auditLogger:  auditLogger,
	}
}

// CreateClassSchedule creates schedule entries linking a class to a teacher/slot.
func (uc *ClassScheduleUseCase) CreateClassSchedule(ctx context.Context, req dto.CreateClassScheduleRequest) (*dto.ClassScheduleResponse, error) {
	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.classRepo.GetByID(ctx, classID); err != nil {
		return nil, domainErrors.ErrClassNotFound
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if teacher, err := uc.teacherRepo.GetByID(ctx, teacherID); err != nil || teacher == nil {
		return nil, domainErrors.ErrTeacherNotFound
	} else if !teacher.IsActive {
		return nil, domainErrors.ErrBadRequest
	}

	dormID, err := uuid.Parse(req.DormitoryID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.dormRepo.GetByID(ctx, dormID); err != nil {
		return nil, domainErrors.ErrDormitoryNotFound
	}

	var subjectID *uuid.UUID
	if req.SubjectID != nil {
		parsed, err := uuid.Parse(*req.SubjectID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if _, err := uc.subjectRepo.GetByID(ctx, parsed); err != nil {
			return nil, domainErrors.ErrSubjectNotFound
		}
		subjectID = &parsed
	}

	startTime, endTime, slotID, err := uc.resolveScheduleTiming(ctx, dormID, req.SlotID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	schedule := &entity.ClassSchedule{
		ID:          uuid.New(),
		ClassID:     classID,
		DormitoryID: dormID,
		SubjectID:   subjectID,
		TeacherID:   teacherID,
		SlotID:      slotID,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    req.Location,
		Notes:       req.Notes,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class_schedule", "class_schedule:create", schedule.ID.String(), map[string]string{
		"class_id":   schedule.ClassID.String(),
		"teacher_id": schedule.TeacherID.String(),
	})

	return uc.toClassScheduleResponse(schedule), nil
}

// GetClassSchedule fetches a schedule by ID.
func (uc *ClassScheduleUseCase) GetClassSchedule(ctx context.Context, id uuid.UUID) (*dto.ClassScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrClassScheduleNotFound
	}
	return uc.toClassScheduleResponse(schedule), nil
}

// ListClassSchedules supports filtering + pagination.
func (uc *ClassScheduleUseCase) ListClassSchedules(
	ctx context.Context,
	classIDStr, teacherIDStr, dormitoryIDStr, dayOfWeek string,
	page, pageSize int,
	isActive *bool,
) (*dto.ListClassSchedulesResponse, error) {
	filter := repository.ClassScheduleFilter{
		DayOfWeek: dayOfWeek,
		IsActive:  isActive,
		Page:      page,
		PageSize:  pageSize,
	}
	if classIDStr != "" {
		classID, err := uuid.Parse(classIDStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.ClassID = classID
	}
	if teacherIDStr != "" {
		teacherID, err := uuid.Parse(teacherIDStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.TeacherID = teacherID
	}
	if dormitoryIDStr != "" {
		dormID, err := uuid.Parse(dormitoryIDStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		filter.DormitoryID = dormID
	}

	schedules, total, err := uc.scheduleRepo.List(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.ClassScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		responses = append(responses, *uc.toClassScheduleResponse(schedule))
	}

	page, pageSize = normalizePagination(page, pageSize)
	return &dto.ListClassSchedulesResponse{
		Schedules:  responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// UpdateClassSchedule updates mutable fields.
func (uc *ClassScheduleUseCase) UpdateClassSchedule(ctx context.Context, id uuid.UUID, req dto.UpdateClassScheduleRequest) (*dto.ClassScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrClassScheduleNotFound
	}

	if req.SubjectID != nil {
		parsed, err := uuid.Parse(*req.SubjectID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if _, err := uc.subjectRepo.GetByID(ctx, parsed); err != nil {
			return nil, domainErrors.ErrSubjectNotFound
		}
		schedule.SubjectID = &parsed
	}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		if teacher, err := uc.teacherRepo.GetByID(ctx, parsed); err != nil || teacher == nil {
			return nil, domainErrors.ErrTeacherNotFound
		} else if !teacher.IsActive {
			return nil, domainErrors.ErrBadRequest
		}
		schedule.TeacherID = parsed
	}

	if req.SlotID != nil || req.StartTime != nil || req.EndTime != nil {
		start, end, slotID, err := uc.resolveScheduleTiming(ctx, schedule.DormitoryID, req.SlotID, req.StartTime, req.EndTime)
		if err != nil {
			return nil, err
		}
		schedule.SlotID = slotID
		schedule.StartTime = start
		schedule.EndTime = end
	}

	if req.DayOfWeek != nil {
		schedule.DayOfWeek = *req.DayOfWeek
	}
	if req.Location != nil {
		schedule.Location = *req.Location
	}
	if req.Notes != nil {
		schedule.Notes = *req.Notes
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}
	schedule.UpdatedAt = time.Now()

	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "class_schedule", "class_schedule:update", schedule.ID.String(), map[string]string{
		"class_id": schedule.ClassID.String(),
	})

	return uc.toClassScheduleResponse(schedule), nil
}

// DeleteClassSchedule deletes schedule entry.
func (uc *ClassScheduleUseCase) DeleteClassSchedule(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.scheduleRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrClassScheduleNotFound
	}
	if err := uc.scheduleRepo.Delete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}
	_ = uc.auditLogger.Log(ctx, "class_schedule", "class_schedule:delete", id.String(), nil)
	return nil
}

func (uc *ClassScheduleUseCase) resolveScheduleTiming(
	ctx context.Context,
	dormID uuid.UUID,
	slotIDStr *string,
	startStr, endStr *string,
) (*time.Time, *time.Time, *uuid.UUID, error) {
	if slotIDStr != nil {
		slotID, err := uuid.Parse(*slotIDStr)
		if err != nil {
			return nil, nil, nil, domainErrors.ErrBadRequest
		}
		slot, err := uc.slotRepo.GetByID(ctx, slotID)
		if err != nil {
			return nil, nil, nil, domainErrors.ErrScheduleSlotNotFound
		}
		if !slot.IsActive {
			return nil, nil, nil, domainErrors.ErrScheduleSlotInactive
		}
		if slot.DormitoryID != dormID {
			return nil, nil, nil, domainErrors.ErrBadRequest
		}
		start := slot.StartTime
		end := slot.EndTime
		return &start, &end, &slotID, nil
	}

	if startStr == nil || endStr == nil {
		return nil, nil, nil, domainErrors.ErrBadRequest
	}
	startTime, err := time.Parse(time.RFC3339, *startStr)
	if err != nil {
		return nil, nil, nil, domainErrors.ErrBadRequest
	}
	endTime, err := time.Parse(time.RFC3339, *endStr)
	if err != nil {
		return nil, nil, nil, domainErrors.ErrBadRequest
	}
	if !startTime.Before(endTime) {
		return nil, nil, nil, domainErrors.ErrBadRequest
	}
	return &startTime, &endTime, nil, nil
}

func (uc *ClassScheduleUseCase) toClassScheduleResponse(schedule *entity.ClassSchedule) *dto.ClassScheduleResponse {
	var subjectID *string
	if schedule.SubjectID != nil {
		val := schedule.SubjectID.String()
		subjectID = &val
	}
	var slotID *string
	if schedule.SlotID != nil {
		val := schedule.SlotID.String()
		slotID = &val
	}
	var startStr, endStr *string
	if schedule.StartTime != nil {
		formatted := schedule.StartTime.Format(time.RFC3339)
		startStr = &formatted
	}
	if schedule.EndTime != nil {
		formatted := schedule.EndTime.Format(time.RFC3339)
		endStr = &formatted
	}

	return &dto.ClassScheduleResponse{
		ID:          schedule.ID.String(),
		ClassID:     schedule.ClassID.String(),
		DormitoryID: schedule.DormitoryID.String(),
		SubjectID:   subjectID,
		TeacherID:   schedule.TeacherID.String(),
		SlotID:      slotID,
		DayOfWeek:   schedule.DayOfWeek,
		StartTime:   startStr,
		EndTime:     endStr,
		Location:    schedule.Location,
		Notes:       schedule.Notes,
		IsActive:    schedule.IsActive,
		CreatedAt:   schedule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   schedule.UpdatedAt.Format(time.RFC3339),
	}
}
