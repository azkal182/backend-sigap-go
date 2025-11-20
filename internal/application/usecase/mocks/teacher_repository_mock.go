package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// MockTeacherRepository is a testify mock for TeacherRepository.
type MockTeacherRepository struct {
	mock.Mock
}

var _ repository.TeacherRepository = (*MockTeacherRepository)(nil)

func (m *MockTeacherRepository) Create(ctx context.Context, teacher *entity.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *MockTeacherRepository) Update(ctx context.Context, teacher *entity.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *MockTeacherRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error) {
	args := m.Called(ctx, id)
	if teacher, ok := args.Get(0).(*entity.Teacher); ok {
		return teacher, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTeacherRepository) GetByCode(ctx context.Context, code string) (*entity.Teacher, error) {
	args := m.Called(ctx, code)
	if teacher, ok := args.Get(0).(*entity.Teacher); ok {
		return teacher, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTeacherRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Teacher, error) {
	args := m.Called(ctx, userID)
	if teacher, ok := args.Get(0).(*entity.Teacher); ok {
		return teacher, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTeacherRepository) List(ctx context.Context, filter repository.TeacherFilter) ([]*entity.Teacher, int64, error) {
	args := m.Called(ctx, filter)
	teachers, _ := args.Get(0).([]*entity.Teacher)
	total := args.Get(1).(int64)
	return teachers, total, args.Error(2)
}

func (m *MockTeacherRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
