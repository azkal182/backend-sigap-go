package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

type stubProvinceRepo struct {
	listFn func(ctx context.Context, page, pageSize int, search string) ([]*entity.Province, int64, error)
	getFn  func(ctx context.Context, id int) (*entity.Province, error)
}

func (s *stubProvinceRepo) List(ctx context.Context, page, pageSize int, search string) ([]*entity.Province, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, page, pageSize, search)
	}
	return nil, 0, nil
}

func (s *stubProvinceRepo) GetByID(ctx context.Context, id int) (*entity.Province, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

type stubRegencyRepo struct {
	listFn func(ctx context.Context, page, pageSize int, provinceID *int, search string) ([]*entity.Regency, int64, error)
	getFn  func(ctx context.Context, id int) (*entity.Regency, error)
}

func (s *stubRegencyRepo) List(ctx context.Context, page, pageSize int, provinceID *int, search string) ([]*entity.Regency, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, page, pageSize, provinceID, search)
	}
	return nil, 0, nil
}

func (s *stubRegencyRepo) GetByID(ctx context.Context, id int) (*entity.Regency, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

type stubDistrictRepo struct {
	listFn func(ctx context.Context, page, pageSize int, regencyID *int, search string) ([]*entity.District, int64, error)
	getFn  func(ctx context.Context, id int) (*entity.District, error)
}

func (s *stubDistrictRepo) List(ctx context.Context, page, pageSize int, regencyID *int, search string) ([]*entity.District, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, page, pageSize, regencyID, search)
	}
	return nil, 0, nil
}

func (s *stubDistrictRepo) GetByID(ctx context.Context, id int) (*entity.District, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

type stubVillageRepo struct {
	listFn func(ctx context.Context, page, pageSize int, districtID *int, search string) ([]*entity.Village, int64, error)
	getFn  func(ctx context.Context, id int) (*entity.Village, error)
}

func (s *stubVillageRepo) List(ctx context.Context, page, pageSize int, districtID *int, search string) ([]*entity.Village, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, page, pageSize, districtID, search)
	}
	return nil, 0, nil
}

func (s *stubVillageRepo) GetByID(ctx context.Context, id int) (*entity.Village, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}

func TestLocationUseCase_ListProvinces(t *testing.T) {
	ctx := context.Background()
	provinceRepo := &stubProvinceRepo{
		listFn: func(ctx context.Context, page, pageSize int, search string) ([]*entity.Province, int64, error) {
			assert.Equal(t, 2, page)
			assert.Equal(t, 5, pageSize)
			assert.Equal(t, "ja", search)
			return []*entity.Province{{ID: 1, Name: "Jabar", Code: "JB"}}, 11, nil
		},
	}

	uc := NewLocationUseCase(provinceRepo, &stubRegencyRepo{}, &stubDistrictRepo{}, &stubVillageRepo{})
	resp, err := uc.ListProvinces(ctx, 2, 5, "ja")
	assert.NoError(t, err)
	assert.Equal(t, int64(11), resp.Total)
	assert.Equal(t, 3, resp.TotalPages)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "Jabar", resp.Items[0].Name)
}

func TestLocationUseCase_ListRegencies(t *testing.T) {
	ctx := context.Background()
	provinceID := 10
	regencyRepo := &stubRegencyRepo{
		listFn: func(ctx context.Context, page, pageSize int, pid *int, search string) ([]*entity.Regency, int64, error) {
			assert.Equal(t, &provinceID, pid)
			return []*entity.Regency{{ID: 3, Name: "Bandung", Type: "city", Code: "BDG", FullCode: "32.73", ProvinceID: provinceID}}, 1, nil
		},
	}
	uc := NewLocationUseCase(&stubProvinceRepo{}, regencyRepo, &stubDistrictRepo{}, &stubVillageRepo{})
	resp, err := uc.ListRegencies(ctx, 1, 10, &provinceID, "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "Bandung", resp.Items[0].Name)
	assert.Equal(t, "city", resp.Items[0].Type)
}

func TestLocationUseCase_GetVillageByID(t *testing.T) {
	ctx := context.Background()
	villageRepo := &stubVillageRepo{
		getFn: func(ctx context.Context, id int) (*entity.Village, error) {
			assert.Equal(t, 7, id)
			return &entity.Village{ID: 7, Name: "Ciburial", Code: "001", FullCode: "32.73.01", PosCode: "40198", DistrictID: 5}, nil
		},
	}
	uc := NewLocationUseCase(&stubProvinceRepo{}, &stubRegencyRepo{}, &stubDistrictRepo{}, villageRepo)
	resp, err := uc.GetVillageByID(ctx, 7)
	assert.NoError(t, err)
	assert.Equal(t, "Ciburial", resp.Name)
	assert.Equal(t, "40198", resp.PosCode)
}
