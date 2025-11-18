package dto

// ProvinceResponse represents province data in responses
type ProvinceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// RegencyResponse represents regency data in responses
type RegencyResponse struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	ProvinceID int    `json:"province_id"`
}

// DistrictResponse represents district data in responses
type DistrictResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	FullCode  string `json:"full_code"`
	RegencyID int    `json:"regency_id"`
}

// VillageResponse represents village data in responses
type VillageResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	PosCode    string `json:"pos_code"`
	DistrictID int    `json:"district_id"`
}

// PaginatedProvinceResponse represents paginated provinces list
type PaginatedProvinceResponse struct {
	Items      []ProvinceResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// PaginatedRegencyResponse represents paginated regencies list
type PaginatedRegencyResponse struct {
	Items      []RegencyResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// PaginatedDistrictResponse represents paginated districts list
type PaginatedDistrictResponse struct {
	Items      []DistrictResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// PaginatedVillageResponse represents paginated villages list
type PaginatedVillageResponse struct {
	Items      []VillageResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
