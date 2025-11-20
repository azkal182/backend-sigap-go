package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// MockScheduleSlotRepository is a testify mock for ScheduleSlotRepository.
type MockScheduleSlotRepository struct {
	mock.Mock
}

var _ repository.ScheduleSlotRepository = (*MockScheduleSlotRepository)(nil)

func (m *MockScheduleSlotRepository) Create(ctx context.Context, slot *entity.ScheduleSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockScheduleSlotRepository) Update(ctx context.Context, slot *entity.ScheduleSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockScheduleSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error) {
	args := m.Called(ctx, id)
	if slot, ok := args.Get(0).(*entity.ScheduleSlot); ok {
		return slot, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleSlotRepository) GetByDormAndNumber(ctx context.Context, dormitoryID uuid.UUID, slotNumber int) (*entity.ScheduleSlot, error) {
	args := m.Called(ctx, dormitoryID, slotNumber)
	if slot, ok := args.Get(0).(*entity.ScheduleSlot); ok {
		return slot, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleSlotRepository) List(ctx context.Context, filter repository.ScheduleSlotFilter) ([]*entity.ScheduleSlot, int64, error) {
	args := m.Called(ctx, filter)
	slots, _ := args.Get(0).([]*entity.ScheduleSlot)
	total := args.Get(1).(int64)
	return slots, total, args.Error(2)
}

func (m *MockScheduleSlotRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
