package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

type scheduleSlotRepository struct {
	db *gorm.DB
}

// NewScheduleSlotRepository creates a new repository backed by GORM.
func NewScheduleSlotRepository() domainRepo.ScheduleSlotRepository {
	return &scheduleSlotRepository{db: database.DB}
}

func (r *scheduleSlotRepository) Create(ctx context.Context, slot *entity.ScheduleSlot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *scheduleSlotRepository) Update(ctx context.Context, slot *entity.ScheduleSlot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}

func (r *scheduleSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleSlot, error) {
	var slot entity.ScheduleSlot
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&slot).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *scheduleSlotRepository) GetByDormAndNumber(ctx context.Context, dormitoryID uuid.UUID, slotNumber int) (*entity.ScheduleSlot, error) {
	var slot entity.ScheduleSlot
	if err := r.db.WithContext(ctx).
		Where("dormitory_id = ? AND slot_number = ?", dormitoryID, slotNumber).
		First(&slot).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *scheduleSlotRepository) List(ctx context.Context, filter domainRepo.ScheduleSlotFilter) ([]*entity.ScheduleSlot, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.ScheduleSlot{})

	if filter.DormitoryID != uuid.Nil {
		query = query.Where("dormitory_id = ?", filter.DormitoryID)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	page := filter.Page
	pageSize := filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var slots []*entity.ScheduleSlot
	if err := query.Order("slot_number ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&slots).Error; err != nil {
		return nil, 0, err
	}

	return slots, total, nil
}

func (r *scheduleSlotRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.ScheduleSlot{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}
