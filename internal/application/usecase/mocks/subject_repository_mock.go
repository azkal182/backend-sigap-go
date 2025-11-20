package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// SubjectRepositoryMock mocks the SubjectRepository interface.
type SubjectRepositoryMock struct {
	mock.Mock
}

var _ repository.SubjectRepository = (*SubjectRepositoryMock)(nil)

func (m *SubjectRepositoryMock) Create(ctx context.Context, subject *entity.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *SubjectRepositoryMock) Update(ctx context.Context, subject *entity.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *SubjectRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error) {
	args := m.Called(ctx, id)
	if subject, ok := args.Get(0).(*entity.Subject); ok {
		return subject, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubjectRepositoryMock) GetByName(ctx context.Context, name string) (*entity.Subject, error) {
	args := m.Called(ctx, name)
	if subject, ok := args.Get(0).(*entity.Subject); ok {
		return subject, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubjectRepositoryMock) List(ctx context.Context, limit, offset int) ([]*entity.Subject, int64, error) {
	args := m.Called(ctx, limit, offset)
	subjects, _ := args.Get(0).([]*entity.Subject)
	total := args.Get(1).(int64)
	return subjects, total, args.Error(2)
}

func (m *SubjectRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
