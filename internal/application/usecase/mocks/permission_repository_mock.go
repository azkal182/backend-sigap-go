package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// PermissionRepositoryMock implements repository.PermissionRepository for tests.
type PermissionRepositoryMock struct {
	mock.Mock
}

var _ repository.PermissionRepository = (*PermissionRepositoryMock)(nil)

func (m *PermissionRepositoryMock) Create(ctx context.Context, permission *entity.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *PermissionRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	args := m.Called(ctx, id)
	if val := args.Get(0); val != nil {
		return val.(*entity.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *PermissionRepositoryMock) GetBySlug(ctx context.Context, slug string) (*entity.Permission, error) {
	args := m.Called(ctx, slug)
	if val := args.Get(0); val != nil {
		return val.(*entity.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *PermissionRepositoryMock) Update(ctx context.Context, permission *entity.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *PermissionRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *PermissionRepositoryMock) List(ctx context.Context, limit, offset int) ([]*entity.Permission, int64, error) {
	args := m.Called(ctx, limit, offset)
	perms, _ := args.Get(0).([]*entity.Permission)
	total, _ := args.Get(1).(int64)
	return perms, total, args.Error(2)
}
