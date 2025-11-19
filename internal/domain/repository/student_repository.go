package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// StudentRepository defines persistence operations for students.
type StudentRepository interface {
	Create(ctx context.Context, student *entity.Student) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error)
	GetByStudentNumber(ctx context.Context, studentNumber string) (*entity.Student, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Student, int64, error)
	Update(ctx context.Context, student *entity.Student) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, isActive bool) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Dormitory history helpers
	CreateHistory(ctx context.Context, history *entity.StudentDormitoryHistory) error
	GetActiveHistory(ctx context.Context, studentID uuid.UUID) (*entity.StudentDormitoryHistory, error)
	ListHistory(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentDormitoryHistory, error)
	CloseHistory(ctx context.Context, historyID uuid.UUID, endDate time.Time) error
}
