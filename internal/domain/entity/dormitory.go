package entity

import (
	"time"

	"github.com/google/uuid"
)

// Dormitory represents a dormitory entity in the domain
type Dormitory struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Users []User `gorm:"many2many:user_dormitories;" json:"users,omitempty"`
}

// TableName specifies the table name for GORM
func (Dormitory) TableName() string {
	return "dormitories"
}
