package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// StudentClassEnrollmentRepositoryMock mocks enrollment repository operations.
type StudentClassEnrollmentRepositoryMock struct {
	mock.Mock
}

func (m *StudentClassEnrollmentRepositoryMock) Create(ctx context.Context, enrollment *entity.StudentClassEnrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *StudentClassEnrollmentRepositoryMock) GetActiveByStudentAndClass(ctx context.Context, studentID, classID uuid.UUID) (*entity.StudentClassEnrollment, error) {
	args := m.Called(ctx, studentID, classID)
	if enrollment, ok := args.Get(0).(*entity.StudentClassEnrollment); ok {
		return enrollment, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *StudentClassEnrollmentRepositoryMock) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentClassEnrollment, error) {
	args := m.Called(ctx, studentID)
	enrollments, _ := args.Get(0).([]*entity.StudentClassEnrollment)
	return enrollments, args.Error(1)
}

func (m *StudentClassEnrollmentRepositoryMock) CloseEnrollment(ctx context.Context, id uuid.UUID, leftAt time.Time) error {
	args := m.Called(ctx, id, leftAt)
	return args.Error(0)
}
