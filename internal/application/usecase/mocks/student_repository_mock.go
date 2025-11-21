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

// StudentSKSResultRepositoryMock mocks StudentSKSResultRepository behavior.
type StudentSKSResultRepositoryMock struct {
	mock.Mock
}

var _ repository.StudentSKSResultRepository = (*StudentSKSResultRepositoryMock)(nil)

func (m *StudentSKSResultRepositoryMock) Create(ctx context.Context, result *entity.StudentSKSResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *StudentSKSResultRepositoryMock) Update(ctx context.Context, result *entity.StudentSKSResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *StudentSKSResultRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentSKSResult, error) {
	args := m.Called(ctx, id)
	if obj, ok := args.Get(0).(*entity.StudentSKSResult); ok {
		return obj, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *StudentSKSResultRepositoryMock) ListByStudent(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID, limit, offset int) ([]*entity.StudentSKSResult, int64, error) {
	args := m.Called(ctx, studentID, fanID, limit, offset)
	results, _ := args.Get(0).([]*entity.StudentSKSResult)
	total := args.Get(1).(int64)
	return results, total, args.Error(2)
}

func (m *StudentSKSResultRepositoryMock) CountPassedByStudentFan(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID) (int64, error) {
	args := m.Called(ctx, studentID, fanID)
	count := args.Get(0).(int64)
	return count, args.Error(1)
}

// FanCompletionStatusRepositoryMock mocks FanCompletionStatusRepository behavior.
type FanCompletionStatusRepositoryMock struct {
	mock.Mock
}

var _ repository.FanCompletionStatusRepository = (*FanCompletionStatusRepositoryMock)(nil)

func (m *FanCompletionStatusRepositoryMock) Upsert(ctx context.Context, status *entity.FanCompletionStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *FanCompletionStatusRepositoryMock) GetByStudentFan(ctx context.Context, studentID, fanID uuid.UUID) (*entity.FanCompletionStatus, error) {
	args := m.Called(ctx, studentID, fanID)
	if obj, ok := args.Get(0).(*entity.FanCompletionStatus); ok {
		return obj, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *FanCompletionStatusRepositoryMock) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.FanCompletionStatus, error) {
	args := m.Called(ctx, studentID)
	statuses, _ := args.Get(0).([]*entity.FanCompletionStatus)
	return statuses, args.Error(1)
}
