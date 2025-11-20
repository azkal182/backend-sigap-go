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

type subjectNoopAuditLogger struct{}

func (n *subjectNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func TestSubjectUseCase_CreateSubject(t *testing.T) {
	repo := new(mocks.SubjectRepositoryMock)
	logger := &subjectNoopAuditLogger{}
	uc := NewSubjectUseCase(repo, logger)

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := uc.CreateSubject(context.Background(), dto.CreateSubjectRequest{Name: "Fiqh"})
	assert.NoError(t, err)
	assert.Equal(t, "Fiqh", resp.Name)
	repo.AssertExpectations(t)
}

func TestSubjectUseCase_GetSubject(t *testing.T) {
	repo := new(mocks.SubjectRepositoryMock)
	logger := &subjectNoopAuditLogger{}
	uc := NewSubjectUseCase(repo, logger)

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&entity.Subject{ID: id, Name: "Fiqh"}, nil)

	resp, err := uc.GetSubject(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id.String(), resp.ID)

	repoFail := new(mocks.SubjectRepositoryMock)
	ucFail := NewSubjectUseCase(repoFail, logger)
	repoFail.On("GetByID", mock.Anything, id).Return(nil, assert.AnError)

	resp, err = ucFail.GetSubject(context.Background(), id)
	assert.ErrorIs(t, err, domainErrors.ErrSubjectNotFound)
	assert.Nil(t, resp)
}

func TestSubjectUseCase_ListSubjects(t *testing.T) {
	repo := new(mocks.SubjectRepositoryMock)
	logger := &subjectNoopAuditLogger{}
	uc := NewSubjectUseCase(repo, logger)

	repo.On("List", mock.Anything, 10, 0).Return([]*entity.Subject{{ID: uuid.New(), Name: "Fiqh"}}, int64(1), nil)

	resp, err := uc.ListSubjects(context.Background(), 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	repo.AssertExpectations(t)
}

func TestSubjectUseCase_UpdateSubject(t *testing.T) {
	repo := new(mocks.SubjectRepositoryMock)
	logger := &subjectNoopAuditLogger{}
	uc := NewSubjectUseCase(repo, logger)

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&entity.Subject{ID: id, Name: "Fiqh"}, nil)
	repo.On("Update", mock.Anything, mock.Anything).Return(nil)

	newName := "Tafsir"
	resp, err := uc.UpdateSubject(context.Background(), id, dto.UpdateSubjectRequest{Name: &newName})
	assert.NoError(t, err)
	assert.Equal(t, newName, resp.Name)
	repo.AssertExpectations(t)

	repoFail := new(mocks.SubjectRepositoryMock)
	ucFail := NewSubjectUseCase(repoFail, logger)
	repoFail.On("GetByID", mock.Anything, id).Return(nil, assert.AnError)

	resp, err = ucFail.UpdateSubject(context.Background(), id, dto.UpdateSubjectRequest{})
	assert.ErrorIs(t, err, domainErrors.ErrSubjectNotFound)
	assert.Nil(t, resp)
}

func TestSubjectUseCase_DeleteSubject(t *testing.T) {
	repo := new(mocks.SubjectRepositoryMock)
	logger := &subjectNoopAuditLogger{}
	uc := NewSubjectUseCase(repo, logger)

	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&entity.Subject{ID: id}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	err := uc.DeleteSubject(context.Background(), id)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	repoFail := new(mocks.SubjectRepositoryMock)
	ucFail := NewSubjectUseCase(repoFail, logger)
	repoFail.On("GetByID", mock.Anything, id).Return(nil, assert.AnError)

	err = ucFail.DeleteSubject(context.Background(), id)
	assert.ErrorIs(t, err, domainErrors.ErrSubjectNotFound)
}
