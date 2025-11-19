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

// Ensure studentRepository implements StudentRepository interface.
var _ domainRepo.StudentRepository = (*studentRepository)(nil)

type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new Student repository backed by GORM.
func NewStudentRepository() domainRepo.StudentRepository {
	return &studentRepository{db: database.DB}
}

func (r *studentRepository) Create(ctx context.Context, student *entity.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error) {
	var student entity.Student
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*entity.Student, error) {
	var student entity.Student
	if err := r.db.WithContext(ctx).Where("student_number = ?", studentNumber).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Student, int64, error) {
	var (
		students []*entity.Student
		total    int64
	)

	db := r.db.WithContext(ctx)
	if err := db.Model(&entity.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (r *studentRepository) Update(ctx context.Context, student *entity.Student) error {
	return r.db.WithContext(ctx).Save(student).Error
}

func (r *studentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, isActive bool) error {
	updatedAt := time.Now()
	return r.db.WithContext(ctx).Model(&entity.Student{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"is_active":  isActive,
			"updated_at": updatedAt,
		}).Error
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Student{}, id).Error
}

func (r *studentRepository) CreateHistory(ctx context.Context, history *entity.StudentDormitoryHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *studentRepository) GetActiveHistory(ctx context.Context, studentID uuid.UUID) (*entity.StudentDormitoryHistory, error) {
	var history entity.StudentDormitoryHistory
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND end_date IS NULL", studentID).
		Order("start_date DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *studentRepository) ListHistory(ctx context.Context, studentID uuid.UUID) ([]*entity.StudentDormitoryHistory, error) {
	var histories []*entity.StudentDormitoryHistory
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("start_date DESC").
		Find(&histories).Error
	return histories, err
}

func (r *studentRepository) CloseHistory(ctx context.Context, historyID uuid.UUID, endDate time.Time) error {
	updatedAt := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.StudentDormitoryHistory{}).
		Where("id = ?", historyID).
		Updates(map[string]interface{}{
			"end_date":   endDate,
			"updated_at": updatedAt,
		}).Error
}
