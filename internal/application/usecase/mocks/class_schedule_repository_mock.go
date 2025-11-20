package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// ClassScheduleRepositoryMock mocks ClassScheduleRepository behavior.
type ClassScheduleRepositoryMock struct {
	mock.Mock
}

var _ repository.ClassScheduleRepository = (*ClassScheduleRepositoryMock)(nil)

func (m *ClassScheduleRepositoryMock) Create(ctx context.Context, schedule *entity.ClassSchedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *ClassScheduleRepositoryMock) Update(ctx context.Context, schedule *entity.ClassSchedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *ClassScheduleRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSchedule, error) {
	args := m.Called(ctx, id)
	if sched, ok := args.Get(0).(*entity.ClassSchedule); ok {
		return sched, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ClassScheduleRepositoryMock) List(ctx context.Context, filter repository.ClassScheduleFilter) ([]*entity.ClassSchedule, int64, error) {
	args := m.Called(ctx, filter)
	schedules, _ := args.Get(0).([]*entity.ClassSchedule)
	total := args.Get(1).(int64)
	return schedules, total, args.Error(2)
}

func (m *ClassScheduleRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
