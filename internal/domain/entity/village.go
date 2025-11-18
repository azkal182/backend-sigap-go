package entity

// Village represents a village entity (desa/kelurahan)
type Village struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	PosCode    string `json:"pos_code"`
	DistrictID int    `json:"kecamatan_id" gorm:"column:district_id"`
}

func (Village) TableName() string {
	return "villages"
}
