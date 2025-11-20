package entity

import (
	"time"

	"github.com/google/uuid"
)

// SKSDefinition captures a competency or exam definition per FAN.
type SKSDefinition struct {
	ID          uuid.UUID  `json:"id"`
	FanID       uuid.UUID  `json:"fan_id" gorm:"not null;index"`
	SubjectID   *uuid.UUID `json:"subject_id" gorm:"index"`
	Code        string     `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Name        string     `json:"name" gorm:"size:150;not null"`
	KKM         float64    `json:"kkm"`
	Description string     `json:"description" gorm:"size:255"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName overrides the default table name.
func (SKSDefinition) TableName() string {
	return "sks_definitions"
}
