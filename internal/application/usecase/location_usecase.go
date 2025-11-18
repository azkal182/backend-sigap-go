package usecase

import (
	"context"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// LocationUseCase handles read-only location lookups
type LocationUseCase struct {
	provinceRepo repository.ProvinceRepository
	regencyRepo  repository.RegencyRepository
	districtRepo repository.DistrictRepository
	villageRepo  repository.VillageRepository
}

func NewLocationUseCase(
	provinceRepo repository.ProvinceRepository,
	regencyRepo repository.RegencyRepository,
	districtRepo repository.DistrictRepository,
	villageRepo repository.VillageRepository,
) *LocationUseCase {
	return &LocationUseCase{
		provinceRepo: provinceRepo,
		regencyRepo:  regencyRepo,
		districtRepo: districtRepo,
		villageRepo:  villageRepo,
	}
}

// helper to compute total pages
func computeTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}

// Provinces

func (uc *LocationUseCase) ListProvinces(ctx context.Context, page, pageSize int, search string) (*dto.PaginatedProvinceResponse, error) {
	provinces, total, err := uc.provinceRepo.List(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProvinceResponse, 0, len(provinces))
	for _, p := range provinces {
		items = append(items, dto.ProvinceResponse{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		})
	}

	return &dto.PaginatedProvinceResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: computeTotalPages(total, pageSize),
	}, nil
}

func (uc *LocationUseCase) GetProvinceByID(ctx context.Context, id int) (*dto.ProvinceResponse, error) {
	p, err := uc.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.ProvinceResponse{ID: p.ID, Name: p.Name, Code: p.Code}, nil
}

// Regencies

func (uc *LocationUseCase) ListRegencies(ctx context.Context, page, pageSize int, provinceID *int, search string) (*dto.PaginatedRegencyResponse, error) {
	regencies, total, err := uc.regencyRepo.List(ctx, page, pageSize, provinceID, search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RegencyResponse, 0, len(regencies))
	for _, r := range regencies {
		items = append(items, dto.RegencyResponse{
			ID:         r.ID,
			Type:       r.Type,
			Name:       r.Name,
			Code:       r.Code,
			FullCode:   r.FullCode,
			ProvinceID: r.ProvinceID,
		})
	}

	return &dto.PaginatedRegencyResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: computeTotalPages(total, pageSize),
	}, nil
}

func (uc *LocationUseCase) GetRegencyByID(ctx context.Context, id int) (*dto.RegencyResponse, error) {
	r, err := uc.regencyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.RegencyResponse{
		ID:         r.ID,
		Type:       r.Type,
		Name:       r.Name,
		Code:       r.Code,
		FullCode:   r.FullCode,
		ProvinceID: r.ProvinceID,
	}, nil
}

// Districts

func (uc *LocationUseCase) ListDistricts(ctx context.Context, page, pageSize int, regencyID *int, search string) (*dto.PaginatedDistrictResponse, error) {
	districts, total, err := uc.districtRepo.List(ctx, page, pageSize, regencyID, search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DistrictResponse, 0, len(districts))
	for _, d := range districts {
		items = append(items, dto.DistrictResponse{
			ID:        d.ID,
			Name:      d.Name,
			Code:      d.Code,
			FullCode:  d.FullCode,
			RegencyID: d.RegencyID,
		})
	}

	return &dto.PaginatedDistrictResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: computeTotalPages(total, pageSize),
	}, nil
}

func (uc *LocationUseCase) GetDistrictByID(ctx context.Context, id int) (*dto.DistrictResponse, error) {
	d, err := uc.districtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.DistrictResponse{
		ID:        d.ID,
		Name:      d.Name,
		Code:      d.Code,
		FullCode:  d.FullCode,
		RegencyID: d.RegencyID,
	}, nil
}

// Villages

func (uc *LocationUseCase) ListVillages(ctx context.Context, page, pageSize int, districtID *int, search string) (*dto.PaginatedVillageResponse, error) {
	villages, total, err := uc.villageRepo.List(ctx, page, pageSize, districtID, search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.VillageResponse, 0, len(villages))
	for _, v := range villages {
		items = append(items, dto.VillageResponse{
			ID:         v.ID,
			Name:       v.Name,
			Code:       v.Code,
			FullCode:   v.FullCode,
			PosCode:    v.PosCode,
			DistrictID: v.DistrictID,
		})
	}

	return &dto.PaginatedVillageResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: computeTotalPages(total, pageSize),
	}, nil
}

func (uc *LocationUseCase) GetVillageByID(ctx context.Context, id int) (*dto.VillageResponse, error) {
	v, err := uc.villageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.VillageResponse{
		ID:         v.ID,
		Name:       v.Name,
		Code:       v.Code,
		FullCode:   v.FullCode,
		PosCode:    v.PosCode,
		DistrictID: v.DistrictID,
	}, nil
}
