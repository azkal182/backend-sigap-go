package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// FanRepository defines persistence operations for fans.
type FanRepository interface {
	Create(ctx context.Context, fan *entity.Fan) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Fan, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Fan, int64, error)
	Update(ctx context.Context, fan *entity.Fan) error
	Delete(ctx context.Context, id uuid.UUID) error
}
