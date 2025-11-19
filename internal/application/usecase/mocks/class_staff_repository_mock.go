package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ClassStaffRepositoryMock mocks class staff repository operations.
type ClassStaffRepositoryMock struct {
	mock.Mock
}

func (m *ClassStaffRepositoryMock) Assign(ctx context.Context, staff *entity.ClassStaff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}

func (m *ClassStaffRepositoryMock) ListByClass(ctx context.Context, classID uuid.UUID) ([]*entity.ClassStaff, error) {
	args := m.Called(ctx, classID)
	staff, _ := args.Get(0).([]*entity.ClassStaff)
	return staff, args.Error(1)
}

func (m *ClassStaffRepositoryMock) Remove(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
