package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ClassRepositoryMock implements repository.ClassRepository for tests.
type ClassRepositoryMock struct {
	mock.Mock
}

func (m *ClassRepositoryMock) Create(ctx context.Context, class *entity.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}

func (m *ClassRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error) {
	args := m.Called(ctx, id)
	if class, ok := args.Get(0).(*entity.Class); ok {
		return class, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ClassRepositoryMock) ListByFan(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.Class, int64, error) {
	args := m.Called(ctx, fanID, limit, offset)
	classes, _ := args.Get(0).([]*entity.Class)
	total := args.Get(1).(int64)
	return classes, total, args.Error(2)
}

func (m *ClassRepositoryMock) Update(ctx context.Context, class *entity.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}

func (m *ClassRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
