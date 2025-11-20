package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

type roleNoopAuditLogger struct{}

func (n *roleNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestRoleUseCase_CreateRole_Success(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)

	roleRepo.On("GetBySlug", ctx, "admin").Return(nil, domainErrors.ErrRoleNotFound)
	roleRepo.On("Create", ctx, mock.Anything).Return(nil)
	roleRepo.On("GetWithPermissions", ctx, mock.Anything).Return(&entity.Role{ID: uuid.New(), Name: "Admin", Slug: "admin"}, nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	resp, err := uc.CreateRole(ctx, dto.CreateRoleRequest{Name: "Admin", Slug: "admin"})
	assert.NoError(t, err)
	assert.Equal(t, "Admin", resp.Name)
	roleRepo.AssertExpectations(t)
}

func TestRoleUseCase_CreateRole_DuplicateSlug(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)
	roleRepo.On("GetBySlug", ctx, "existing").Return(&entity.Role{ID: uuid.New()}, nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	resp, err := uc.CreateRole(ctx, dto.CreateRoleRequest{Name: "Existing", Slug: "existing"})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrRoleAlreadyExists)
}

func TestRoleUseCase_UpdateRole_SlugConflict(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)
	roleID := uuid.New()
	roleRepo.On("GetByID", ctx, roleID).Return(&entity.Role{ID: roleID, Slug: "old"}, nil)
	roleRepo.On("GetBySlug", ctx, "taken").Return(&entity.Role{ID: uuid.New()}, nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	resp, err := uc.UpdateRole(ctx, roleID, dto.UpdateRoleRequest{Slug: "taken"})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, domainErrors.ErrRoleAlreadyExists)
}

func TestRoleUseCase_DeleteRole_Protected(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	roleID := uuid.New()
	roleRepo.On("GetByID", ctx, roleID).Return(&entity.Role{ID: roleID, IsProtected: true}, nil)

	uc := NewRoleUseCase(roleRepo, new(mocks.PermissionRepositoryMock), &roleNoopAuditLogger{})
	err := uc.DeleteRole(ctx, roleID)
	assert.ErrorIs(t, err, domainErrors.ErrProtectedRole)
}

func TestRoleUseCase_ListRoles_Paginates(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)
	roleRepo.On("List", ctx, 10, 0).Return([]*entity.Role{{ID: uuid.New(), Name: "Viewer"}}, int64(1), nil)
	roleRepo.On("GetWithPermissions", ctx, mock.Anything).Return(&entity.Role{ID: uuid.New(), Name: "Viewer", Permissions: []entity.Permission{{Name: "read"}}}, nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	resp, err := uc.ListRoles(ctx, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 1, resp.TotalPages)
	assert.Len(t, resp.Roles, 1)
	roleRepo.AssertExpectations(t)
}

func TestRoleUseCase_AssignPermission(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)
	roleID := uuid.New()
	permID := uuid.New()
	roleRepo.On("GetByID", ctx, roleID).Return(&entity.Role{ID: roleID, IsProtected: false}, nil)
	permissionRepo.On("GetByID", ctx, permID).Return(&entity.Permission{ID: permID}, nil)
	roleRepo.On("AssignPermission", ctx, roleID, permID).Return(nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	assert.NoError(t, uc.AssignPermission(ctx, roleID, permID))
	roleRepo.AssertExpectations(t)
	permissionRepo.AssertExpectations(t)
}

func TestRoleUseCase_RemovePermission_Protected(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	roleID := uuid.New()
	permissionID := uuid.New()
	roleRepo.On("GetByID", ctx, roleID).Return(&entity.Role{ID: roleID, IsProtected: true}, nil)

	uc := NewRoleUseCase(roleRepo, new(mocks.PermissionRepositoryMock), &roleNoopAuditLogger{})
	err := uc.RemovePermission(ctx, roleID, permissionID)
	assert.ErrorIs(t, err, domainErrors.ErrProtectedRole)
}

func TestRoleUseCase_UpdateRole_Success(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	permissionRepo := new(mocks.PermissionRepositoryMock)
	roleID := uuid.New()
	updated := &entity.Role{ID: roleID, Name: "Updated", Slug: "updated", UpdatedAt: time.Now()}
	roleRepo.On("GetByID", ctx, roleID).Return(&entity.Role{ID: roleID, Name: "Old"}, nil)
	roleRepo.On("Update", ctx, mock.Anything).Return(nil)
	roleRepo.On("GetWithPermissions", ctx, roleID).Return(updated, nil)

	uc := NewRoleUseCase(roleRepo, permissionRepo, &roleNoopAuditLogger{})
	resp, err := uc.UpdateRole(ctx, roleID, dto.UpdateRoleRequest{Name: "Updated"})
	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Name)
}

func TestRoleUseCase_ListRoles_RepoError(t *testing.T) {
	ctx := context.Background()
	roleRepo := new(mocks.MockRoleRepository)
	roleRepo.On("List", ctx, 5, 5).Return(nil, int64(0), errors.New("db"))

	uc := NewRoleUseCase(roleRepo, new(mocks.PermissionRepositoryMock), &roleNoopAuditLogger{})
	resp, err := uc.ListRoles(ctx, 2, 5)
	assert.Nil(t, resp)
	assert.Error(t, err)
}
