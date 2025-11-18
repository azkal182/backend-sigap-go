package entity

// Regency represents a regency/city entity (kabupaten/kota)
type Regency struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	FullCode   string `json:"full_code"`
	ProvinceID int    `json:"provinsi_id" gorm:"column:province_id"`
}

func (Regency) TableName() string {
	return "regencies"
}
