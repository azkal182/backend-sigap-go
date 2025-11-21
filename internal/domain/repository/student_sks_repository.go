package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// StudentSKSResultRepository handles persistence for student SKS exam outcomes.
type StudentSKSResultRepository interface {
	Create(ctx context.Context, result *entity.StudentSKSResult) error
	Update(ctx context.Context, result *entity.StudentSKSResult) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentSKSResult, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID, limit, offset int) ([]*entity.StudentSKSResult, int64, error)
	CountPassedByStudentFan(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID) (int64, error)
}

// FanCompletionStatusRepository stores the completion state per student/FAN.
type FanCompletionStatusRepository interface {
	Upsert(ctx context.Context, status *entity.FanCompletionStatus) error
	GetByStudentFan(ctx context.Context, studentID, fanID uuid.UUID) (*entity.FanCompletionStatus, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.FanCompletionStatus, error)
}
