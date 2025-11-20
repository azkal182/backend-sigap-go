package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// SKSDefinitionRepository handles persistence for SKS definitions.
type SKSDefinitionRepository interface {
	Create(ctx context.Context, sks *entity.SKSDefinition) error
	Update(ctx context.Context, sks *entity.SKSDefinition) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSDefinition, error)
	GetByCode(ctx context.Context, code string) (*entity.SKSDefinition, error)
	List(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.SKSDefinition, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SKSExamScheduleRepository handles persistence for SKS exam schedules.
type SKSExamScheduleRepository interface {
	Create(ctx context.Context, exam *entity.SKSExamSchedule) error
	Update(ctx context.Context, exam *entity.SKSExamSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSExamSchedule, error)
	ListBySKS(ctx context.Context, sksID uuid.UUID, limit, offset int) ([]*entity.SKSExamSchedule, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
