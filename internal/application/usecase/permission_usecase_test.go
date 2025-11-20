package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

func TestPermissionUseCase_ListPermissions_NormalizesPagination(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.PermissionRepositoryMock)

	mockPermissions := []*entity.Permission{
		{ID: uuid.New(), Name: "users:read", Slug: "users-read", Resource: "users", Action: "read"},
	}
	mockRepo.On("List", ctx, 10, 0).Return(mockPermissions, int64(1), nil)

	uc := NewPermissionUseCase(mockRepo)
	resp, err := uc.ListPermissions(ctx, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, resp.TotalPages)
	assert.Equal(t, "users:read", resp.Permissions[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestPermissionUseCase_ListPermissions_WhenRepoFails(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.PermissionRepositoryMock)
	mockRepo.On("List", ctx, 5, 5).Return(nil, int64(0), errors.New("db err"))

	uc := NewPermissionUseCase(mockRepo)
	resp, err := uc.ListPermissions(ctx, 2, 5)
	assert.Nil(t, resp)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
