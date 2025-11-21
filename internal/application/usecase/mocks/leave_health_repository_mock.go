package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// LeavePermitRepositoryMock mocks repository.LeavePermitRepository.
type LeavePermitRepositoryMock struct {
	mock.Mock
}

var _ repository.LeavePermitRepository = (*LeavePermitRepositoryMock)(nil)

func (m *LeavePermitRepositoryMock) Create(ctx context.Context, permit *entity.LeavePermit) error {
	return m.Called(ctx, permit).Error(0)
}

func (m *LeavePermitRepositoryMock) Update(ctx context.Context, permit *entity.LeavePermit) error {
	return m.Called(ctx, permit).Error(0)
}

func (m *LeavePermitRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.LeavePermit, error) {
	args := m.Called(ctx, id)
	if permit, ok := args.Get(0).(*entity.LeavePermit); ok {
		return permit, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LeavePermitRepositoryMock) List(ctx context.Context, filter repository.LeavePermitFilter) ([]*entity.LeavePermit, int64, error) {
	args := m.Called(ctx, filter)
	permits, _ := args.Get(0).([]*entity.LeavePermit)
	var total int64
	if v := args.Get(1); v != nil {
		total = v.(int64)
	}
	return permits, total, args.Error(2)
}

func (m *LeavePermitRepositoryMock) HasOverlap(ctx context.Context, studentID uuid.UUID, startDate, endDate time.Time, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, studentID, startDate, endDate, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *LeavePermitRepositoryMock) ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error) {
	args := m.Called(ctx, studentID, date)
	if permit, ok := args.Get(0).(*entity.LeavePermit); ok {
		return permit, args.Error(1)
	}
	return nil, args.Error(1)
}

// HealthStatusRepositoryMock mocks repository.HealthStatusRepository.
type HealthStatusRepositoryMock struct {
	mock.Mock
}

var _ repository.HealthStatusRepository = (*HealthStatusRepositoryMock)(nil)

func (m *HealthStatusRepositoryMock) Create(ctx context.Context, status *entity.HealthStatus) error {
	return m.Called(ctx, status).Error(0)
}

func (m *HealthStatusRepositoryMock) Update(ctx context.Context, status *entity.HealthStatus) error {
	return m.Called(ctx, status).Error(0)
}

func (m *HealthStatusRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.HealthStatus, error) {
	args := m.Called(ctx, id)
	if status, ok := args.Get(0).(*entity.HealthStatus); ok {
		return status, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *HealthStatusRepositoryMock) List(ctx context.Context, filter repository.HealthStatusFilter) ([]*entity.HealthStatus, int64, error) {
	args := m.Called(ctx, filter)
	statuses, _ := args.Get(0).([]*entity.HealthStatus)
	var total int64
	if v := args.Get(1); v != nil {
		total = v.(int64)
	}
	return statuses, total, args.Error(2)
}

func (m *HealthStatusRepositoryMock) ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error) {
	args := m.Called(ctx, studentID, date)
	if status, ok := args.Get(0).(*entity.HealthStatus); ok {
		return status, args.Error(1)
	}
	return nil, args.Error(1)
}
