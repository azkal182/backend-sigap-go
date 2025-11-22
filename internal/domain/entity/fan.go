package entity

import (
	"time"

	"github.com/google/uuid"
)

// Fan represents an academic fan/stream within the pesantren.
type Fan struct {
	ID          uuid.UUID  `json:"id"`
	DormitoryID uuid.UUID  `json:"dormitory_id" gorm:"type:uuid;index"`
	Name        string     `json:"name" gorm:"size:150;not null"`
	Level       string     `json:"level" gorm:"size:50;not null"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName specifies the table name for GORM.
func (Fan) TableName() string {
	return "fans"
}
