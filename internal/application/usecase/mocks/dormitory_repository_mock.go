package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// MockDormitoryRepository is a mock implementation of DormitoryRepository
type MockDormitoryRepository struct {
	mock.Mock
}

// Ensure MockDormitoryRepository implements repository.DormitoryRepository
var _ repository.DormitoryRepository = (*MockDormitoryRepository)(nil)

func (m *MockDormitoryRepository) Create(ctx context.Context, dormitory *entity.Dormitory) error {
	args := m.Called(ctx, dormitory)
	return args.Error(0)
}

func (m *MockDormitoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Dormitory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Dormitory), args.Error(1)
}

func (m *MockDormitoryRepository) Update(ctx context.Context, dormitory *entity.Dormitory) error {
	args := m.Called(ctx, dormitory)
	return args.Error(0)
}

func (m *MockDormitoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDormitoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Dormitory, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Dormitory), args.Get(1).(int64), args.Error(2)
}

func (m *MockDormitoryRepository) AssignToUser(ctx context.Context, userID, dormitoryID uuid.UUID) error {
	args := m.Called(ctx, userID, dormitoryID)
	return args.Error(0)
}

func (m *MockDormitoryRepository) RemoveFromUser(ctx context.Context, userID, dormitoryID uuid.UUID) error {
	args := m.Called(ctx, userID, dormitoryID)
	return args.Error(0)
}

func (m *MockDormitoryRepository) GetUserDormitories(ctx context.Context, userID uuid.UUID) ([]*entity.Dormitory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Dormitory), args.Error(1)
}
