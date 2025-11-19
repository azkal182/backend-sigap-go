package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ClassStaffRepository defines operations for class staff assignments.
type ClassStaffRepository interface {
	Assign(ctx context.Context, staff *entity.ClassStaff) error
	ListByClass(ctx context.Context, classID uuid.UUID) ([]*entity.ClassStaff, error)
	Remove(ctx context.Context, id uuid.UUID) error
}
