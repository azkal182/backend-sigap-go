package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ScheduleSlotRepository defines persistence operations for schedule slots.
type ScheduleSlotRepository interface {
	Create(ctx context.Context, slot *entity.ScheduleSlot) error
	Update(ctx context.Context, slot *entity.ScheduleSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error)
	GetByDormAndNumber(ctx context.Context, dormitoryID uuid.UUID, slotNumber int) (*entity.ScheduleSlot, error)
	List(ctx context.Context, filter ScheduleSlotFilter) ([]*entity.ScheduleSlot, int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// ScheduleSlotFilter encapsulates query filters for listing schedule slots.
type ScheduleSlotFilter struct {
	DormitoryID uuid.UUID
	IsActive    *bool
	Page        int
	PageSize    int
}
