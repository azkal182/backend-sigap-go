package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// ScheduleSlotUseCase orchestrates slot operations.
type ScheduleSlotUseCase struct {
	slotRepo    repository.ScheduleSlotRepository
	dormRepo    repository.DormitoryRepository
	auditLogger appService.AuditLogger
}

// NewScheduleSlotUseCase constructs slot use case.
func NewScheduleSlotUseCase(
	slotRepo repository.ScheduleSlotRepository,
	dormRepo repository.DormitoryRepository,
	auditLogger appService.AuditLogger,
) *ScheduleSlotUseCase {
	return &ScheduleSlotUseCase{slotRepo: slotRepo, dormRepo: dormRepo, auditLogger: auditLogger}
}

// CreateScheduleSlot creates a new slot.
func (uc *ScheduleSlotUseCase) CreateScheduleSlot(ctx context.Context, req dto.CreateScheduleSlotRequest) (*dto.ScheduleSlotResponse, error) {
	dormID, err := uuid.Parse(req.DormitoryID)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}
	if _, err := uc.dormRepo.GetByID(ctx, dormID); err != nil {
		return nil, domainErrors.ErrDormitoryNotFound
	}

	startTime, endTime, err := parseSlotTimes(req.StartTime, req.EndTime)
	if err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	if existing, _ := uc.slotRepo.GetByDormAndNumber(ctx, dormID, req.SlotNumber); existing != nil {
		return nil, domainErrors.ErrScheduleSlotConflict
	}

	if err := uc.ensureNoOverlap(ctx, dormID, startTime, endTime, uuid.Nil); err != nil {
		return nil, err
	}

	slot := &entity.ScheduleSlot{
		ID:          uuid.New(),
		DormitoryID: dormID,
		SlotNumber:  req.SlotNumber,
		Name:        req.Name,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.slotRepo.Create(ctx, slot); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "schedule_slot", "schedule_slot:create", slot.ID.String(), map[string]string{
		"dormitory_id": dormID.String(),
		"slot_number":  formatInt(req.SlotNumber),
	})

	return uc.toScheduleSlotResponse(slot), nil
}

// ListScheduleSlots lists slots with filters.
func (uc *ScheduleSlotUseCase) ListScheduleSlots(ctx context.Context, dormitoryID string, page, pageSize int, isActive *bool) (*dto.ListScheduleSlotsResponse, error) {
	var dormID uuid.UUID
	var err error
	if dormitoryID != "" {
		dormID, err = uuid.Parse(dormitoryID)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
	}

	page, pageSize = normalizePagination(page, pageSize)
	slots, total, err := uc.slotRepo.List(ctx, repository.ScheduleSlotFilter{
		DormitoryID: dormID,
		IsActive:    isActive,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	responses := make([]dto.ScheduleSlotResponse, 0, len(slots))
	for _, slot := range slots {
		responses = append(responses, *uc.toScheduleSlotResponse(slot))
	}

	return &dto.ListScheduleSlotsResponse{
		Slots:      responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

// GetScheduleSlot retrieves slot by ID.
func (uc *ScheduleSlotUseCase) GetScheduleSlot(ctx context.Context, id uuid.UUID) (*dto.ScheduleSlotResponse, error) {
	slot, err := uc.slotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrScheduleSlotNotFound
	}
	return uc.toScheduleSlotResponse(slot), nil
}

// UpdateScheduleSlot updates slot fields.
func (uc *ScheduleSlotUseCase) UpdateScheduleSlot(ctx context.Context, id uuid.UUID, req dto.UpdateScheduleSlotRequest) (*dto.ScheduleSlotResponse, error) {
	slot, err := uc.slotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrScheduleSlotNotFound
	}

	if req.SlotNumber != nil && *req.SlotNumber != slot.SlotNumber {
		if existing, _ := uc.slotRepo.GetByDormAndNumber(ctx, slot.DormitoryID, *req.SlotNumber); existing != nil && existing.ID != slot.ID {
			return nil, domainErrors.ErrScheduleSlotConflict
		}
		slot.SlotNumber = *req.SlotNumber
	}

	if req.Name != nil {
		slot.Name = *req.Name
	}
	if req.Description != nil {
		slot.Description = *req.Description
	}

	updatedStart := slot.StartTime
	updatedEnd := slot.EndTime
	if req.StartTime != nil || req.EndTime != nil {
		startStr := slot.StartTime.Format(time.RFC3339)
		endStr := slot.EndTime.Format(time.RFC3339)
		if req.StartTime != nil {
			startStr = *req.StartTime
		}
		if req.EndTime != nil {
			endStr = *req.EndTime
		}
		start, end, err := parseSlotTimes(startStr, endStr)
		if err != nil {
			return nil, domainErrors.ErrBadRequest
		}
		updatedStart = start
		updatedEnd = end
	}

	if err := uc.ensureNoOverlap(ctx, slot.DormitoryID, updatedStart, updatedEnd, slot.ID); err != nil {
		return nil, err
	}

	slot.StartTime = updatedStart
	slot.EndTime = updatedEnd

	if req.IsActive != nil {
		slot.IsActive = *req.IsActive
	}

	slot.UpdatedAt = time.Now()
	if err := uc.slotRepo.Update(ctx, slot); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	_ = uc.auditLogger.Log(ctx, "schedule_slot", "schedule_slot:update", slot.ID.String(), map[string]string{
		"slot_number": formatInt(slot.SlotNumber),
	})

	return uc.toScheduleSlotResponse(slot), nil
}

// DeleteScheduleSlot soft deletes a slot.
func (uc *ScheduleSlotUseCase) DeleteScheduleSlot(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.slotRepo.GetByID(ctx, id); err != nil {
		return domainErrors.ErrScheduleSlotNotFound
	}
	if err := uc.slotRepo.SoftDelete(ctx, id); err != nil {
		return domainErrors.ErrInternalServer
	}
	_ = uc.auditLogger.Log(ctx, "schedule_slot", "schedule_slot:delete", id.String(), nil)
	return nil
}

func (uc *ScheduleSlotUseCase) ensureNoOverlap(ctx context.Context, dormID uuid.UUID, start, end time.Time, ignoreID uuid.UUID) error {
	page := 1
	pageSize := 100
	for {
		slots, total, err := uc.slotRepo.List(ctx, repository.ScheduleSlotFilter{
			DormitoryID: dormID,
			Page:        page,
			PageSize:    pageSize,
		})
		if err != nil {
			return domainErrors.ErrInternalServer
		}
		for _, slot := range slots {
			if slot.ID == ignoreID {
				continue
			}
			if timesOverlap(start, end, slot.StartTime, slot.EndTime) {
				return domainErrors.ErrScheduleSlotConflict
			}
		}
		if page*pageSize >= int(total) {
			break
		}
		page++
	}
	return nil
}

func parseSlotTimes(startStr, endStr string) (time.Time, time.Time, error) {
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if !start.Before(end) {
		return time.Time{}, time.Time{}, domainErrors.ErrBadRequest
	}
	return start, end, nil
}

func timesOverlap(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && endA.After(startB)
}

func (uc *ScheduleSlotUseCase) toScheduleSlotResponse(slot *entity.ScheduleSlot) *dto.ScheduleSlotResponse {
	return &dto.ScheduleSlotResponse{
		ID:          slot.ID.String(),
		DormitoryID: slot.DormitoryID.String(),
		SlotNumber:  slot.SlotNumber,
		Name:        slot.Name,
		StartTime:   slot.StartTime.Format(time.RFC3339),
		EndTime:     slot.EndTime.Format(time.RFC3339),
		Description: slot.Description,
		IsActive:    slot.IsActive,
		CreatedAt:   slot.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   slot.UpdatedAt.Format(time.RFC3339),
	}
}

func formatInt(val int) string {
	return strconv.FormatInt(int64(val), 10)
}
