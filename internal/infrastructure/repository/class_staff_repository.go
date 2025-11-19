package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ domainRepo.ClassStaffRepository = (*classStaffRepository)(nil)

type classStaffRepository struct {
	db *gorm.DB
}

// NewClassStaffRepository returns a repository backed by the shared DB instance.
func NewClassStaffRepository() domainRepo.ClassStaffRepository {
	return &classStaffRepository{db: database.DB}
}

func (r *classStaffRepository) Assign(ctx context.Context, staff *entity.ClassStaff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *classStaffRepository) ListByClass(ctx context.Context, classID uuid.UUID) ([]*entity.ClassStaff, error) {
	var staff []*entity.ClassStaff
	err := r.db.WithContext(ctx).
		Where("class_id = ?", classID).
		Order("created_at DESC").
		Find(&staff).Error
	return staff, err
}

func (r *classStaffRepository) Remove(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ClassStaff{}, id).Error
}
