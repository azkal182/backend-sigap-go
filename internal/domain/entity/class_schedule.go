package entity

import (
	"time"

	"github.com/google/uuid"
)

// ClassSchedule represents a recurring schedule entry for a class.
type ClassSchedule struct {
	ID          uuid.UUID  `json:"id"`
	ClassID     uuid.UUID  `json:"class_id" gorm:"not null;index"`
	DormitoryID uuid.UUID  `json:"dormitory_id" gorm:"not null;index"`
	SubjectID   *uuid.UUID `json:"subject_id" gorm:"index"`
	TeacherID   uuid.UUID  `json:"teacher_id" gorm:"not null;index"`
	SlotID      *uuid.UUID `json:"slot_id" gorm:"index"`
	DayOfWeek   string     `json:"day_of_week" gorm:"size:16;not null"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Location    string     `json:"location" gorm:"size:150"`
	Notes       string     `json:"notes" gorm:"size:255"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName overrides the default table name for GORM.
func (ClassSchedule) TableName() string {
	return "class_schedules"
}
