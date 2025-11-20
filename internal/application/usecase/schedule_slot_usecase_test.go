package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

type scheduleSlotNoopAuditLogger struct{}

func (n *scheduleSlotNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestScheduleSlotUseCase_CreateSlot_Success(t *testing.T) {
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	logger := &scheduleSlotNoopAuditLogger{}
	uc := NewScheduleSlotUseCase(slotRepo, dormRepo, logger)

	dormID := uuid.New()
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)
	slotRepo.On("GetByDormAndNumber", mock.Anything, dormID, 1).Return(nil, domainErrors.ErrScheduleSlotNotFound)
	slotRepo.On("List", mock.Anything, mock.Anything).
		Return([]*entity.ScheduleSlot{}, int64(0), nil)
	slotRepo.On("Create", mock.Anything, mock.MatchedBy(func(slot *entity.ScheduleSlot) bool {
		return slot.DormitoryID == dormID && slot.SlotNumber == 1
	})).Return(nil)

	resp, err := uc.CreateScheduleSlot(context.Background(), dto.CreateScheduleSlotRequest{
		DormitoryID: dormID.String(),
		SlotNumber:  1,
		Name:        "Morning Slot",
		StartTime:   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
		EndTime:     time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.SlotNumber)
	slotRepo.AssertExpectations(t)
	dormRepo.AssertExpectations(t)
}

func TestScheduleSlotUseCase_CreateSlot_OverlapConflict(t *testing.T) {
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	logger := &scheduleSlotNoopAuditLogger{}
	uc := NewScheduleSlotUseCase(slotRepo, dormRepo, logger)

	dormID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(time.Hour)

	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)
	slotRepo.On("GetByDormAndNumber", mock.Anything, dormID, 1).Return(nil, domainErrors.ErrScheduleSlotNotFound)
	slotRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ScheduleSlotFilter) bool {
		return filter.DormitoryID == dormID
	})).Return([]*entity.ScheduleSlot{{
		ID:          uuid.New(),
		DormitoryID: dormID,
		SlotNumber:  2,
		StartTime:   start,
		EndTime:     end,
	}}, int64(1), nil)

	_, err := uc.CreateScheduleSlot(context.Background(), dto.CreateScheduleSlotRequest{
		DormitoryID: dormID.String(),
		SlotNumber:  1,
		Name:        "Overlap",
		StartTime:   start.Add(30 * time.Minute).Format(time.RFC3339),
		EndTime:     end.Add(30 * time.Minute).Format(time.RFC3339),
	})

	assert.ErrorIs(t, err, domainErrors.ErrScheduleSlotConflict)
	slotRepo.AssertExpectations(t)
}

func TestScheduleSlotUseCase_ListSlots_NormalizesPagination(t *testing.T) {
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	logger := &scheduleSlotNoopAuditLogger{}
	uc := NewScheduleSlotUseCase(slotRepo, dormRepo, logger)

	slot := &entity.ScheduleSlot{
		ID:          uuid.New(),
		DormitoryID: uuid.New(),
		SlotNumber:  3,
		Name:        "Slot",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		IsActive:    true,
	}

	slotRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ScheduleSlotFilter) bool {
		return filter.Page == 1 && filter.PageSize == 10
	})).Return([]*entity.ScheduleSlot{slot}, int64(1), nil)

	resp, err := uc.ListScheduleSlots(context.Background(), "", 0, 0, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Slots, 1)
	assert.Equal(t, slot.SlotNumber, resp.Slots[0].SlotNumber)
	slotRepo.AssertExpectations(t)
}

func TestScheduleSlotUseCase_UpdateSlot_ChangeNumberAndTimes(t *testing.T) {
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	logger := &scheduleSlotNoopAuditLogger{}
	uc := NewScheduleSlotUseCase(slotRepo, dormRepo, logger)

	dormID := uuid.New()
	slotID := uuid.New()
	currentStart := time.Now().UTC()
	currentEnd := currentStart.Add(time.Hour)

	slotRepo.On("GetByID", mock.Anything, slotID).Return(&entity.ScheduleSlot{
		ID:          slotID,
		DormitoryID: dormID,
		SlotNumber:  1,
		Name:        "Old Slot",
		StartTime:   currentStart,
		EndTime:     currentEnd,
		IsActive:    true,
		CreatedAt:   currentStart,
		UpdatedAt:   currentStart,
	}, nil)
	slotRepo.On("GetByDormAndNumber", mock.Anything, dormID, 2).Return(nil, domainErrors.ErrScheduleSlotNotFound)
	slotRepo.On("List", mock.Anything, mock.Anything).Return([]*entity.ScheduleSlot{}, int64(0), nil)
	slotRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	newStart := currentStart.Add(2 * time.Hour)
	newEnd := newStart.Add(time.Hour)
	newName := "Updated Slot"
	desc := "desc"
	isActive := false
	newSlotNum := 2

	resp, err := uc.UpdateScheduleSlot(context.Background(), slotID, dto.UpdateScheduleSlotRequest{
		SlotNumber:  &newSlotNum,
		Name:        &newName,
		StartTime:   stringPtr(newStart.Format(time.RFC3339)),
		EndTime:     stringPtr(newEnd.Format(time.RFC3339)),
		Description: &desc,
		IsActive:    &isActive,
	})

	assert.NoError(t, err)
	assert.Equal(t, newSlotNum, resp.SlotNumber)
	assert.Equal(t, newName, resp.Name)
	assert.Equal(t, isActive, resp.IsActive)
	slotRepo.AssertExpectations(t)
}

func TestScheduleSlotUseCase_DeleteSlot(t *testing.T) {
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	logger := &scheduleSlotNoopAuditLogger{}
	uc := NewScheduleSlotUseCase(slotRepo, dormRepo, logger)

	slotID := uuid.New()
	slotRepo.On("GetByID", mock.Anything, slotID).Return(&entity.ScheduleSlot{ID: slotID}, nil)
	slotRepo.On("SoftDelete", mock.Anything, slotID).Return(nil)

	err := uc.DeleteScheduleSlot(context.Background(), slotID)

	assert.NoError(t, err)
	slotRepo.AssertExpectations(t)
}

func stringPtr(val string) *string {
	return &val
}
