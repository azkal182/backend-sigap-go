package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// FanRepositoryMock implements repository.FanRepository for tests.
type FanRepositoryMock struct {
	mock.Mock
}

func (m *FanRepositoryMock) Create(ctx context.Context, fan *entity.Fan) error {
	args := m.Called(ctx, fan)
	return args.Error(0)
}

func (m *FanRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Fan, error) {
	args := m.Called(ctx, id)
	if fan, ok := args.Get(0).(*entity.Fan); ok {
		return fan, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *FanRepositoryMock) List(ctx context.Context, limit, offset int) ([]*entity.Fan, int64, error) {
	args := m.Called(ctx, limit, offset)
	fans, _ := args.Get(0).([]*entity.Fan)
	total := args.Get(1).(int64)
	return fans, total, args.Error(2)
}

func (m *FanRepositoryMock) Update(ctx context.Context, fan *entity.Fan) error {
	args := m.Called(ctx, fan)
	return args.Error(0)
}

func (m *FanRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
