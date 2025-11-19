package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ domainRepo.FanRepository = (*fanRepository)(nil)

type fanRepository struct {
	db *gorm.DB
}

// NewFanRepository returns a FanRepository backed by the shared DB instance.
func NewFanRepository() domainRepo.FanRepository {
	return &fanRepository{db: database.DB}
}

func (r *fanRepository) Create(ctx context.Context, fan *entity.Fan) error {
	return r.db.WithContext(ctx).Create(fan).Error
}

func (r *fanRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Fan, error) {
	var fan entity.Fan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&fan).Error; err != nil {
		return nil, err
	}
	return &fan, nil
}

func (r *fanRepository) List(ctx context.Context, limit, offset int) ([]*entity.Fan, int64, error) {
	var (
		fans  []*entity.Fan
		total int64
	)

	db := r.db.WithContext(ctx)
	if err := db.Model(&entity.Fan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&fans).Error; err != nil {
		return nil, 0, err
	}

	return fans, total, nil
}

func (r *fanRepository) Update(ctx context.Context, fan *entity.Fan) error {
	return r.db.WithContext(ctx).Save(fan).Error
}

func (r *fanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Fan{}, id).Error
}
