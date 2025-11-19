package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ClassRepository defines persistence operations for classes.
type ClassRepository interface {
	Create(ctx context.Context, class *entity.Class) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error)
	ListByFan(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.Class, int64, error)
	Update(ctx context.Context, class *entity.Class) error
	Delete(ctx context.Context, id uuid.UUID) error
}
