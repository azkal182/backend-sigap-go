package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// SKSDefinitionRepositoryMock mocks SKSDefinitionRepository behavior.
type SKSDefinitionRepositoryMock struct {
	mock.Mock
}

var _ repository.SKSDefinitionRepository = (*SKSDefinitionRepositoryMock)(nil)

func (m *SKSDefinitionRepositoryMock) Create(ctx context.Context, sks *entity.SKSDefinition) error {
	args := m.Called(ctx, sks)
	return args.Error(0)
}

func (m *SKSDefinitionRepositoryMock) Update(ctx context.Context, sks *entity.SKSDefinition) error {
	args := m.Called(ctx, sks)
	return args.Error(0)
}

func (m *SKSDefinitionRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSDefinition, error) {
	args := m.Called(ctx, id)
	if def, ok := args.Get(0).(*entity.SKSDefinition); ok {
		return def, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SKSDefinitionRepositoryMock) GetByCode(ctx context.Context, code string) (*entity.SKSDefinition, error) {
	args := m.Called(ctx, code)
	if def, ok := args.Get(0).(*entity.SKSDefinition); ok {
		return def, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SKSDefinitionRepositoryMock) List(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.SKSDefinition, int64, error) {
	args := m.Called(ctx, fanID, limit, offset)
	defs, _ := args.Get(0).([]*entity.SKSDefinition)
	total := args.Get(1).(int64)
	return defs, total, args.Error(2)
}

func (m *SKSDefinitionRepositoryMock) CountByFan(ctx context.Context, fanID uuid.UUID) (int64, error) {
	args := m.Called(ctx, fanID)
	total := args.Get(0).(int64)
	return total, args.Error(1)
}

func (m *SKSDefinitionRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SKSExamScheduleRepositoryMock mocks SKSExamScheduleRepository behavior.
type SKSExamScheduleRepositoryMock struct {
	mock.Mock
}

var _ repository.SKSExamScheduleRepository = (*SKSExamScheduleRepositoryMock)(nil)

func (m *SKSExamScheduleRepositoryMock) Create(ctx context.Context, exam *entity.SKSExamSchedule) error {
	args := m.Called(ctx, exam)
	return args.Error(0)
}

func (m *SKSExamScheduleRepositoryMock) Update(ctx context.Context, exam *entity.SKSExamSchedule) error {
	args := m.Called(ctx, exam)
	return args.Error(0)
}

func (m *SKSExamScheduleRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSExamSchedule, error) {
	args := m.Called(ctx, id)
	if exam, ok := args.Get(0).(*entity.SKSExamSchedule); ok {
		return exam, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SKSExamScheduleRepositoryMock) ListBySKS(ctx context.Context, sksID uuid.UUID, limit, offset int) ([]*entity.SKSExamSchedule, int64, error) {
	args := m.Called(ctx, sksID, limit, offset)
	exams, _ := args.Get(0).([]*entity.SKSExamSchedule)
	total := args.Get(1).(int64)
	return exams, total, args.Error(2)
}

func (m *SKSExamScheduleRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
