package entity

import (
	"time"

	"github.com/google/uuid"
)

// Dormitory represents a dormitory entity in the domain
type Dormitory struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name" gorm:"size:100;not null"`
	Gender      string     `json:"gender" gorm:"size:10;not null"`
	Level       string     `json:"level" gorm:"size:50;not null"`
	Code        string     `json:"code" gorm:"size:16;not null;uniqueIndex"`
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
