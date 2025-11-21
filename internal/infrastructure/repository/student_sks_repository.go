package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_ domainRepo.StudentSKSResultRepository    = (*studentSKSResultRepository)(nil)
	_ domainRepo.FanCompletionStatusRepository = (*fanCompletionStatusRepository)(nil)
)

type studentSKSResultRepository struct {
	db *gorm.DB
}

type fanCompletionStatusRepository struct {
	db *gorm.DB
}

// NewStudentSKSResultRepository wires the repository using the shared DB instance.
func NewStudentSKSResultRepository() domainRepo.StudentSKSResultRepository {
	return &studentSKSResultRepository{db: database.DB}
}

// NewFanCompletionStatusRepository wires the fan completion repository using the shared DB instance.
func NewFanCompletionStatusRepository() domainRepo.FanCompletionStatusRepository {
	return &fanCompletionStatusRepository{db: database.DB}
}

func (r *studentSKSResultRepository) Create(ctx context.Context, result *entity.StudentSKSResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *studentSKSResultRepository) Update(ctx context.Context, result *entity.StudentSKSResult) error {
	return r.db.WithContext(ctx).Save(result).Error
}

func (r *studentSKSResultRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentSKSResult, error) {
	var res entity.StudentSKSResult
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *studentSKSResultRepository) ListByStudent(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID, limit, offset int) ([]*entity.StudentSKSResult, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := r.db.WithContext(ctx).
		Model(&entity.StudentSKSResult{}).
		Where("student_id = ?", studentID)

	if fanID != uuid.Nil {
		query = query.Joins("JOIN sks_definitions ON sks_definitions.id = student_sks_results.sks_id").
			Where("sks_definitions.fan_id = ?", fanID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var results []*entity.StudentSKSResult
	if err := query.Order("student_sks_results.created_at DESC").Limit(limit).Offset(offset).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *studentSKSResultRepository) CountPassedByStudentFan(ctx context.Context, studentID uuid.UUID, fanID uuid.UUID) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.StudentSKSResult{}).
		Joins("JOIN sks_definitions ON sks_definitions.id = student_sks_results.sks_id").
		Where("student_sks_results.student_id = ?", studentID).
		Where("sks_definitions.fan_id = ?", fanID).
		Where("student_sks_results.is_passed = ?", true)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *fanCompletionStatusRepository) Upsert(ctx context.Context, status *entity.FanCompletionStatus) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "student_id"}, {Name: "fan_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_completed", "completed_at", "updated_at"}),
		}).
		Create(status).Error
}

func (r *fanCompletionStatusRepository) GetByStudentFan(ctx context.Context, studentID, fanID uuid.UUID) (*entity.FanCompletionStatus, error) {
	var status entity.FanCompletionStatus
	if err := r.db.WithContext(ctx).
		Where("student_id = ? AND fan_id = ?", studentID, fanID).
		First(&status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *fanCompletionStatusRepository) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.FanCompletionStatus, error) {
	var statuses []*entity.FanCompletionStatus
	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("fan_id ASC").
		Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
