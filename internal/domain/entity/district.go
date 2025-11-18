package entity

// District represents a district entity (kecamatan)
type District struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	FullCode  string `json:"full_code"`
	RegencyID int    `json:"kabupaten_id" gorm:"column:regency_id"`
}

func (District) TableName() string {
	return "districts"
}
