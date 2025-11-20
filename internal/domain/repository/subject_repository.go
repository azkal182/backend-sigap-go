package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// SubjectRepository defines operations for subject data.
type SubjectRepository interface {
	Create(ctx context.Context, subject *entity.Subject) error
	Update(ctx context.Context, subject *entity.Subject) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error)
	GetByName(ctx context.Context, name string) (*entity.Subject, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Subject, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
