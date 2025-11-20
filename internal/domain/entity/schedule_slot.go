package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleSlot represents a shared time slot definition per dormitory.
type ScheduleSlot struct {
	ID          uuid.UUID      `json:"id"`
	DormitoryID uuid.UUID      `json:"dormitory_id" gorm:"not null;index;uniqueIndex:idx_slot_dorm_number"`
	SlotNumber  int            `json:"slot_number" gorm:"not null;uniqueIndex:idx_slot_dorm_number"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	StartTime   time.Time      `json:"start_time" gorm:"not null"`
	EndTime     time.Time      `json:"end_time" gorm:"not null"`
	Description string         `json:"description" gorm:"size:255"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName overrides GORM table name.
func (ScheduleSlot) TableName() string {
	return "schedule_slots"
}
