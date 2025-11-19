package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ domainRepo.StudentClassEnrollmentRepository = (*classEnrollmentRepository)(nil)

type classEnrollmentRepository struct {
	db *gorm.DB
}

// NewStudentClassEnrollmentRepository returns a repository backed by the shared DB instance.
func NewStudentClassEnrollmentRepository() domainRepo.StudentClassEnrollmentRepository {
	return &classEnrollmentRepository{db: database.DB}
}

func (r *classEnrollmentRepository) Create(ctx context.Context, enrollment *entity.StudentClassEnrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *classEnrollmentRepository) GetActiveByStudentAndClass(ctx context.Context, studentID, classID uuid.UUID) (*entity.StudentClassEnrollment, error) {
	var enrollment entity.StudentClassEnrollment
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND class_id = ? AND left_at IS NULL", studentID, classID).
		Order("enrolled_at DESC").
		First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *classEnrollmentRepository) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentClassEnrollment, error) {
	var enrollments []*entity.StudentClassEnrollment
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("enrolled_at DESC").
		Find(&enrollments).Error
	return enrollments, err
}

func (r *classEnrollmentRepository) CloseEnrollment(ctx context.Context, id uuid.UUID, leftAt time.Time) error {
	updatedAt := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.StudentClassEnrollment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"left_at":    leftAt,
			"updated_at": updatedAt,
		}).Error
}
