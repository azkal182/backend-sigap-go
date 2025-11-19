package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ domainRepo.ClassRepository = (*classRepository)(nil)

type classRepository struct {
	db *gorm.DB
}

// NewClassRepository returns a ClassRepository backed by the shared DB instance.
func NewClassRepository() domainRepo.ClassRepository {
	return &classRepository{db: database.DB}
}

func (r *classRepository) Create(ctx context.Context, class *entity.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

func (r *classRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error) {
	var class entity.Class
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) ListByFan(ctx context.Context, fanID uuid.UUID, limit, offset int) ([]*entity.Class, int64, error) {
	var (
		classes []*entity.Class
		total   int64
	)

	db := r.db.WithContext(ctx).Model(&entity.Class{}).Where("fan_id = ?", fanID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&classes).Error; err != nil {
		return nil, 0, err
	}
	return classes, total, nil
}

func (r *classRepository) Update(ctx context.Context, class *entity.Class) error {
	return r.db.WithContext(ctx).Save(class).Error
}

func (r *classRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Class{}, id).Error
}
