package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// MockStudentRepository is a mock implementation of StudentRepository
// generated manually for unit testing use cases.
type MockStudentRepository struct {
	mock.Mock
}

// Ensure interface compliance
var _ repository.StudentRepository = (*MockStudentRepository)(nil)

func (m *MockStudentRepository) Create(ctx context.Context, student *entity.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Student), args.Error(1)
}

func (m *MockStudentRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*entity.Student, error) {
	args := m.Called(ctx, studentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Student), args.Error(1)
}

func (m *MockStudentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Student, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Student), args.Get(1).(int64), args.Error(2)
}

func (m *MockStudentRepository) Update(ctx context.Context, student *entity.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, isActive bool) error {
	args := m.Called(ctx, id, status, isActive)
	return args.Error(0)
}

func (m *MockStudentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStudentRepository) CreateHistory(ctx context.Context, history *entity.StudentDormitoryHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockStudentRepository) GetActiveHistory(ctx context.Context, studentID uuid.UUID) (*entity.StudentDormitoryHistory, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.StudentDormitoryHistory), args.Error(1)
}

func (m *MockStudentRepository) ListHistory(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentDormitoryHistory, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.StudentDormitoryHistory), args.Error(1)
}

func (m *MockStudentRepository) CloseHistory(ctx context.Context, historyID uuid.UUID, endDate time.Time) error {
	args := m.Called(ctx, historyID, endDate)
	return args.Error(0)
}
