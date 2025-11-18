package repository

import (
	"context"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
)

type provinceRepository struct{}

type regencyRepository struct{}

type districtRepository struct{}

type villageRepository struct{}

func NewProvinceRepository() repository.ProvinceRepository { return &provinceRepository{} }
func NewRegencyRepository() repository.RegencyRepository   { return &regencyRepository{} }
func NewDistrictRepository() repository.DistrictRepository { return &districtRepository{} }
func NewVillageRepository() repository.VillageRepository   { return &villageRepository{} }

func (r *provinceRepository) GetByID(ctx context.Context, id int) (*entity.Province, error) {
	var province entity.Province
	if err := database.DB.WithContext(ctx).First(&province, id).Error; err != nil {
		return nil, err
	}
	return &province, nil
}

func (r *provinceRepository) List(ctx context.Context, page, pageSize int, search string) ([]*entity.Province, int64, error) {
	db := database.DB.WithContext(ctx).Model(&entity.Province{})
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var provinces []*entity.Province
	if err := db.Limit(pageSize).Offset(offset).Order("id ASC").Find(&provinces).Error; err != nil {
		return nil, 0, err
	}
	return provinces, total, nil
}

func (r *regencyRepository) GetByID(ctx context.Context, id int) (*entity.Regency, error) {
	var regency entity.Regency
	if err := database.DB.WithContext(ctx).First(&regency, id).Error; err != nil {
		return nil, err
	}
	return &regency, nil
}

func (r *regencyRepository) List(ctx context.Context, page, pageSize int, provinceID *int, search string) ([]*entity.Regency, int64, error) {
	db := database.DB.WithContext(ctx).Model(&entity.Regency{})
	if provinceID != nil {
		db = db.Where("province_id = ?", *provinceID)
	}
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var regencies []*entity.Regency
	if err := db.Limit(pageSize).Offset(offset).Order("id ASC").Find(&regencies).Error; err != nil {
		return nil, 0, err
	}
	return regencies, total, nil
}

func (r *districtRepository) GetByID(ctx context.Context, id int) (*entity.District, error) {
	var district entity.District
	if err := database.DB.WithContext(ctx).First(&district, id).Error; err != nil {
		return nil, err
	}
	return &district, nil
}

func (r *districtRepository) List(ctx context.Context, page, pageSize int, regencyID *int, search string) ([]*entity.District, int64, error) {
	db := database.DB.WithContext(ctx).Model(&entity.District{})
	if regencyID != nil {
		db = db.Where("regency_id = ?", *regencyID)
	}
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var districts []*entity.District
	if err := db.Limit(pageSize).Offset(offset).Order("id ASC").Find(&districts).Error; err != nil {
		return nil, 0, err
	}
	return districts, total, nil
}

func (r *villageRepository) GetByID(ctx context.Context, id int) (*entity.Village, error) {
	var village entity.Village
	if err := database.DB.WithContext(ctx).First(&village, id).Error; err != nil {
		return nil, err
	}
	return &village, nil
}

func (r *villageRepository) List(ctx context.Context, page, pageSize int, districtID *int, search string) ([]*entity.Village, int64, error) {
	db := database.DB.WithContext(ctx).Model(&entity.Village{})
	if districtID != nil {
		db = db.Where("district_id = ?", *districtID)
	}
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var villages []*entity.Village
	if err := db.Limit(pageSize).Offset(offset).Order("id ASC").Find(&villages).Error; err != nil {
		return nil, 0, err
	}
	return villages, total, nil
}
