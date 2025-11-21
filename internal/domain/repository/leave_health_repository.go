package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// LeavePermitFilter collects optional filters for listing permits.
type LeavePermitFilter struct {
	StudentID *uuid.UUID
	Status    *entity.LeavePermitStatus
	Type      *entity.LeavePermitType
	Date      *time.Time // filter permits overlapping specific date
	Limit     int
	Offset    int
}

// HealthStatusFilter collects filters for health status queries.
type HealthStatusFilter struct {
	StudentID *uuid.UUID
	Status    *entity.HealthStatusState
	Date      *time.Time // overlapping date
	Limit     int
	Offset    int
}

// LeavePermitRepository defines persistence behavior for leave permits.
type LeavePermitRepository interface {
	Create(ctx context.Context, permit *entity.LeavePermit) error
	Update(ctx context.Context, permit *entity.LeavePermit) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LeavePermit, error)
	List(ctx context.Context, filter LeavePermitFilter) ([]*entity.LeavePermit, int64, error)
	HasOverlap(ctx context.Context, studentID uuid.UUID, startDate, endDate time.Time, excludeID *uuid.UUID) (bool, error)
	ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error)
}

// HealthStatusRepository defines persistence behavior for health statuses.
type HealthStatusRepository interface {
	Create(ctx context.Context, status *entity.HealthStatus) error
	Update(ctx context.Context, status *entity.HealthStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.HealthStatus, error)
	List(ctx context.Context, filter HealthStatusFilter) ([]*entity.HealthStatus, int64, error)
	ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error)
}
