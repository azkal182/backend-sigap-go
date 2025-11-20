package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// TeacherRepository defines persistence operations for teachers.
type TeacherRepository interface {
	Create(ctx context.Context, teacher *entity.Teacher) error
	Update(ctx context.Context, teacher *entity.Teacher) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error)
	GetByCode(ctx context.Context, code string) (*entity.Teacher, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Teacher, error)
	List(ctx context.Context, filter TeacherFilter) ([]*entity.Teacher, int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// TeacherFilter encapsulates query filters for listing teachers.
type TeacherFilter struct {
	Keyword  string
	IsActive *bool
	Page     int
	PageSize int
}
