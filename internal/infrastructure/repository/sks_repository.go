package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

var (
	_ domainRepo.SKSDefinitionRepository   = (*sksDefinitionRepository)(nil)
	_ domainRepo.SKSExamScheduleRepository = (*sksExamScheduleRepository)(nil)
)

type sksDefinitionRepository struct {
	db *gorm.DB
}

type sksExamScheduleRepository struct {
	db *gorm.DB
}

// NewSKSDefinitionRepository wires a repository backed by the shared DB instance.
func NewSKSDefinitionRepository() domainRepo.SKSDefinitionRepository {
	return &sksDefinitionRepository{db: database.DB}
}

// NewSKSExamScheduleRepository wires a repository backed by the shared DB instance.
func NewSKSExamScheduleRepository() domainRepo.SKSExamScheduleRepository {
	return &sksExamScheduleRepository{db: database.DB}
}

// SKS Definition operations

func (r *sksDefinitionRepository) Create(ctx context.Context, sks *entity.SKSDefinition) error {
	return r.db.WithContext(ctx).Create(sks).Error
}

func (r *sksDefinitionRepository) Update(ctx context.Context, sks *entity.SKSDefinition) error {
	return r.db.WithContext(ctx).Save(sks).Error
}

func (r *sksDefinitionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSDefinition, error) {
	var definition entity.SKSDefinition
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&definition).Error; err != nil {
		return nil, err
	}
	return &definition, nil
}

func (r *sksDefinitionRepository) GetByCode(ctx context.Context, code string) (*entity.SKSDefinition, error) {
	var definition entity.SKSDefinition
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&definition).Error; err != nil {
		return nil, err
	}
	return &definition, nil
}

func (r *sksDefinitionRepository) List(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.SKSDefinition, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.SKSDefinition{})
	if fanID != uuid.Nil {
		query = query.Where("fan_id = ?", fanID)
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var definitions []*entity.SKSDefinition
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&definitions).Error; err != nil {
		return nil, 0, err
	}

	return definitions, total, nil
}

func (r *sksDefinitionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SKSDefinition{}, id).Error
}

// SKS Exam schedule operations

func (r *sksExamScheduleRepository) Create(ctx context.Context, exam *entity.SKSExamSchedule) error {
	return r.db.WithContext(ctx).Create(exam).Error
}

func (r *sksExamScheduleRepository) Update(ctx context.Context, exam *entity.SKSExamSchedule) error {
	return r.db.WithContext(ctx).Save(exam).Error
}

func (r *sksExamScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SKSExamSchedule, error) {
	var exam entity.SKSExamSchedule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&exam).Error; err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *sksExamScheduleRepository) ListBySKS(ctx context.Context, sksID uuid.UUID, limit, offset int) ([]*entity.SKSExamSchedule, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.SKSExamSchedule{}).Where("sks_id = ?", sksID)

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var exams []*entity.SKSExamSchedule
	if err := query.Order("exam_date ASC, exam_time ASC").Limit(limit).Offset(offset).Find(&exams).Error; err != nil {
		return nil, 0, err
	}

	return exams, total, nil
}

func (r *sksExamScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SKSExamSchedule{}, id).Error
}
