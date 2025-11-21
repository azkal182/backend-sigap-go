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

var _ domainRepo.LeavePermitRepository = (*leavePermitRepository)(nil)
var _ domainRepo.HealthStatusRepository = (*healthStatusRepository)(nil)

type leavePermitRepository struct {
	db *gorm.DB
}

type healthStatusRepository struct {
	db *gorm.DB
}

func NewLeavePermitRepository() domainRepo.LeavePermitRepository {
	return &leavePermitRepository{db: database.DB}
}

func NewHealthStatusRepository() domainRepo.HealthStatusRepository {
	return &healthStatusRepository{db: database.DB}
}

func (r *leavePermitRepository) Create(ctx context.Context, permit *entity.LeavePermit) error {
	return r.db.WithContext(ctx).Create(permit).Error
}

func (r *leavePermitRepository) Update(ctx context.Context, permit *entity.LeavePermit) error {
	return r.db.WithContext(ctx).Save(permit).Error
}

func (r *leavePermitRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.LeavePermit, error) {
	var permit entity.LeavePermit
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&permit).Error; err != nil {
		return nil, err
	}
	return &permit, nil
}

func (r *leavePermitRepository) List(ctx context.Context, filter domainRepo.LeavePermitFilter) ([]*entity.LeavePermit, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.LeavePermit{})

	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Date != nil {
		query = query.Where("start_date <= ? AND end_date >= ?", *filter.Date, *filter.Date)
	}

	limit, offset := normalizePaging(filter.Limit, filter.Offset)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var permits []*entity.LeavePermit
	if err := query.Order("start_date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&permits).Error; err != nil {
		return nil, 0, err
	}

	return permits, total, nil
}

func (r *leavePermitRepository) HasOverlap(ctx context.Context, studentID uuid.UUID, startDate, endDate time.Time, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.LeavePermit{}).
		Where("student_id = ?", studentID).
		Where("status IN ?", []entity.LeavePermitStatus{
			entity.LeavePermitStatusPending,
			entity.LeavePermitStatusApproved,
		})

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	query = query.Where("NOT (end_date < ? OR start_date > ?)", startDate, endDate)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *leavePermitRepository) ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.LeavePermit, error) {
	var permit entity.LeavePermit
	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Where("status IN ?", []entity.LeavePermitStatus{
			entity.LeavePermitStatusPending,
			entity.LeavePermitStatusApproved,
		}).
		Where("start_date <= ? AND end_date >= ?", date, date).
		Order("start_date DESC").
		First(&permit).Error; err != nil {
		return nil, err
	}
	return &permit, nil
}

func (r *healthStatusRepository) Create(ctx context.Context, status *entity.HealthStatus) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *healthStatusRepository) Update(ctx context.Context, status *entity.HealthStatus) error {
	return r.db.WithContext(ctx).Save(status).Error
}

func (r *healthStatusRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.HealthStatus, error) {
	var record entity.HealthStatus
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *healthStatusRepository) List(ctx context.Context, filter domainRepo.HealthStatusFilter) ([]*entity.HealthStatus, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.HealthStatus{})

	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Date != nil {
		query = query.Where("start_date <= ?", *filter.Date).
			Where("(end_date IS NULL OR end_date >= ?)", *filter.Date)
	}

	limit, offset := normalizePaging(filter.Limit, filter.Offset)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var statuses []*entity.HealthStatus
	if err := query.Order("start_date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&statuses).Error; err != nil {
		return nil, 0, err
	}

	return statuses, total, nil
}

func (r *healthStatusRepository) ActiveByDate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.HealthStatus, error) {
	var record entity.HealthStatus
	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Where("status = ?", entity.HealthStatusStateActive).
		Where("start_date <= ?", date).
		Where("(end_date IS NULL OR end_date >= ?)", date).
		Order("start_date DESC").
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func normalizePaging(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
