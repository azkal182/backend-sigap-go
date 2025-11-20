package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

// Ensure classScheduleRepository satisfies the domain interface.
var _ domainRepo.ClassScheduleRepository = (*classScheduleRepository)(nil)

type classScheduleRepository struct {
	db *gorm.DB
}

// NewClassScheduleRepository wires GORM-backed repository instance.
func NewClassScheduleRepository() domainRepo.ClassScheduleRepository {
	return &classScheduleRepository{db: database.DB}
}

func (r *classScheduleRepository) Create(ctx context.Context, schedule *entity.ClassSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *classScheduleRepository) Update(ctx context.Context, schedule *entity.ClassSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *classScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSchedule, error) {
	var schedule entity.ClassSchedule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *classScheduleRepository) List(ctx context.Context, filter domainRepo.ClassScheduleFilter) ([]*entity.ClassSchedule, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.ClassSchedule{})

	if filter.ClassID != uuid.Nil {
		query = query.Where("class_id = ?", filter.ClassID)
	}
	if filter.TeacherID != uuid.Nil {
		query = query.Where("teacher_id = ?", filter.TeacherID)
	}
	if filter.DormitoryID != uuid.Nil {
		query = query.Where("dormitory_id = ?", filter.DormitoryID)
	}
	if filter.DayOfWeek != "" {
		query = query.Where("day_of_week = ?", filter.DayOfWeek)
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

	var schedules []*entity.ClassSchedule
	if err := query.Order("day_of_week ASC, start_time ASC NULLS LAST").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *classScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ClassSchedule{}, id).Error
}
