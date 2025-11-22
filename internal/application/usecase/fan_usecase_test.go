package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

func TestFanUseCase_CreateFan(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.FanRepositoryMock)
	dormRepo := new(mocks.MockDormitoryRepository)
	dormID := uuid.New()

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)

	uc := NewFanUseCase(repo, dormRepo, &noopAuditLogger{})
	resp, err := uc.CreateFan(ctx, dto.CreateFanRequest{Name: "FAN A", Level: "junior", DormitoryID: dormID.String()})

	assert.NoError(t, err)
	assert.Equal(t, "FAN A", resp.Name)
	repo.AssertExpectations(t)
	dormRepo.AssertExpectations(t)
}

func TestFanUseCase_GetFan(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.FanRepositoryMock)
	dormRepo := new(mocks.MockDormitoryRepository)
	id := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(&entity.Fan{ID: id, Name: "Existing"}, nil)

	uc := NewFanUseCase(repo, dormRepo, &noopAuditLogger{})
	resp, err := uc.GetFan(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "Existing", resp.Name)

	repo.AssertExpectations(t)

	repoFail := new(mocks.FanRepositoryMock)
	repoFail.On("GetByID", mock.Anything, id).Return(nil, assert.AnError)
	ucFail := NewFanUseCase(repoFail, dormRepo, &noopAuditLogger{})
	resp, err = ucFail.GetFan(ctx, id)
	assert.ErrorIs(t, err, domainErrors.ErrFanNotFound)
	assert.Nil(t, resp)
}

func TestFanUseCase_ListFans(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.FanRepositoryMock)
	dormRepo := new(mocks.MockDormitoryRepository)
	fans := []*entity.Fan{{ID: uuid.New(), Name: "Fan1"}}
	repo.On("List", mock.Anything, 10, 0).Return(fans, int64(1), nil)

	uc := NewFanUseCase(repo, dormRepo, &noopAuditLogger{})
	resp, err := uc.ListFans(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, resp.Fans, 1)
}

func TestFanUseCase_UpdateFan(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.FanRepositoryMock)
	dormRepo := new(mocks.MockDormitoryRepository)
	id := uuid.New()
	name := "Updated"
	repo.On("GetByID", mock.Anything, id).Return(&entity.Fan{ID: id, Name: "Old"}, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(f *entity.Fan) bool {
		return f.Name == name
	})).Return(nil)

	uc := NewFanUseCase(repo, dormRepo, &noopAuditLogger{})
	resp, err := uc.UpdateFan(ctx, id, dto.UpdateFanRequest{Name: &name})
	assert.NoError(t, err)
	assert.Equal(t, name, resp.Name)
	repo.AssertExpectations(t)
}

func TestFanUseCase_DeleteFan(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.FanRepositoryMock)
	dormRepo := new(mocks.MockDormitoryRepository)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&entity.Fan{ID: id}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	uc := NewFanUseCase(repo, dormRepo, &noopAuditLogger{})
	err := uc.DeleteFan(ctx, id)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
