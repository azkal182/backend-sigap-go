package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subject represents an academic subject that can be attached to class schedules or SKS definitions.
type Subject struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name" gorm:"size:150;not null;uniqueIndex"`
	Description string     `json:"description" gorm:"size:255"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName overrides the default table name for GORM.
func (Subject) TableName() string {
	return "subjects"
}
