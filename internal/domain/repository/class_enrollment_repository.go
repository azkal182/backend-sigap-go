package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// StudentClassEnrollmentRepository defines operations for class enrollments.
type StudentClassEnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entity.StudentClassEnrollment) error
	GetActiveByStudentAndClass(ctx context.Context, studentID, classID uuid.UUID) (*entity.StudentClassEnrollment, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentClassEnrollment, error)
	CloseEnrollment(ctx context.Context, id uuid.UUID, leftAt time.Time) error
}
