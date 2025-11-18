package entity

// Province represents a province entity (static reference data)
type Province struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (Province) TableName() string {
	return "provinces"
}
