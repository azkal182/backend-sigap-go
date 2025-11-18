package repository

import (
	"context"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// ProvinceRepository defines read-only operations for provinces
type ProvinceRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Province, error)
	List(ctx context.Context, page, pageSize int, search string) ([]*entity.Province, int64, error)
}

// RegencyRepository defines read-only operations for regencies
type RegencyRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Regency, error)
	List(ctx context.Context, page, pageSize int, provinceID *int, search string) ([]*entity.Regency, int64, error)
}

// DistrictRepository defines read-only operations for districts
type DistrictRepository interface {
	GetByID(ctx context.Context, id int) (*entity.District, error)
	List(ctx context.Context, page, pageSize int, regencyID *int, search string) ([]*entity.District, int64, error)
}

// VillageRepository defines read-only operations for villages
type VillageRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Village, error)
	List(ctx context.Context, page, pageSize int, districtID *int, search string) ([]*entity.Village, int64, error)
}
