package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

type teacherRepository struct {
	db *gorm.DB
}

// NewTeacherRepository creates a TeacherRepository backed by GORM.
func NewTeacherRepository() domainRepo.TeacherRepository {
	return &teacherRepository{db: database.DB}
}

func (r *teacherRepository) Create(ctx context.Context, teacher *entity.Teacher) error {
	return r.db.WithContext(ctx).Create(teacher).Error
}

func (r *teacherRepository) Update(ctx context.Context, teacher *entity.Teacher) error {
	return r.db.WithContext(ctx).Save(teacher).Error
}

func (r *teacherRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error) {
	var teacher entity.Teacher
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) GetByCode(ctx context.Context, code string) (*entity.Teacher, error) {
	var teacher entity.Teacher
	if err := r.db.WithContext(ctx).Where("teacher_code = ?", code).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Teacher, error) {
	var teacher entity.Teacher
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherRepository) List(ctx context.Context, filter domainRepo.TeacherFilter) ([]*entity.Teacher, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Teacher{})

	if filter.Keyword != "" {
		keyword := strings.ToLower(filter.Keyword)
		like := "%" + keyword + "%"
		query = query.Where("LOWER(full_name) LIKE ? OR LOWER(teacher_code) LIKE ?", like, like)
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

	var teachers []*entity.Teacher
	if err := query.Order("full_name ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&teachers).Error; err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}

func (r *teacherRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Teacher{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}
