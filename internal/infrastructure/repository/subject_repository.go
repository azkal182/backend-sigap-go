package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

// Ensure subjectRepository implements SubjectRepository.
var _ domainRepo.SubjectRepository = (*subjectRepository)(nil)

type subjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository creates a new repository backed by the shared DB instance.
func NewSubjectRepository() domainRepo.SubjectRepository {
	return &subjectRepository{db: database.DB}
}

func (r *subjectRepository) Create(ctx context.Context, subject *entity.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *subjectRepository) Update(ctx context.Context, subject *entity.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

func (r *subjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error) {
	var subject entity.Subject
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&subject).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *subjectRepository) GetByName(ctx context.Context, name string) (*entity.Subject, error) {
	var subject entity.Subject
	if err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(&subject).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *subjectRepository) List(ctx context.Context, limit, offset int) ([]*entity.Subject, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var (
		subjects []*entity.Subject
		total    int64
	)

	query := r.db.WithContext(ctx).Model(&entity.Subject{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&subjects).Error; err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *subjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Subject{}, id).Error
}
