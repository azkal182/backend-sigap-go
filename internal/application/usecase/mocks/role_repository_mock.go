package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// MockRoleRepository is a mock implementation of RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

// Ensure MockRoleRepository implements repository.RoleRepository
var _ repository.RoleRepository = (*MockRoleRepository)(nil)

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetWithPermissions(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) List(ctx context.Context, limit, offset int) ([]*entity.Role, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Role), args.Get(1).(int64), args.Error(2)
}
