package entity

import (
	"time"

	"github.com/google/uuid"
)

// Class represents a class under a specific fan.
type Class struct {
	ID        uuid.UUID  `json:"id"`
	FanID     uuid.UUID  `json:"fan_id" gorm:"index;not null"`
	Name      string     `json:"name" gorm:"size:150;not null"`
	Capacity  int        `json:"capacity"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM.
func (Class) TableName() string {
	return "classes"
}
