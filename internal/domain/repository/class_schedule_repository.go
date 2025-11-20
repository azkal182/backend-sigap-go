package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ClassScheduleFilter encapsulates query parameters for listing schedules.
type ClassScheduleFilter struct {
	ClassID     uuid.UUID
	TeacherID   uuid.UUID
	DormitoryID uuid.UUID
	DayOfWeek   string
	IsActive    *bool
	Page        int
	PageSize    int
}

// ClassScheduleRepository defines persistence operations for class schedules.
type ClassScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.ClassSchedule) error
	Update(ctx context.Context, schedule *entity.ClassSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSchedule, error)
	List(ctx context.Context, filter ClassScheduleFilter) ([]*entity.ClassSchedule, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
