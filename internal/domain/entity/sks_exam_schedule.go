package entity

import (
	"time"

	"github.com/google/uuid"
)

// SKSExamSchedule represents an exam schedule tied to an SKS definition.
type SKSExamSchedule struct {
	ID         uuid.UUID  `json:"id"`
	SKSID      uuid.UUID  `json:"sks_id" gorm:"not null;index"`
	ExaminerID *uuid.UUID `json:"examiner_id" gorm:"index"`
	ExamDate   time.Time  `json:"exam_date" gorm:"not null"`
	ExamTime   time.Time  `json:"exam_time" gorm:"not null"`
	Location   string     `json:"location" gorm:"size:150"`
	Notes      string     `json:"notes" gorm:"size:255"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// TableName overrides default table name.
func (SKSExamSchedule) TableName() string {
	return "sks_exam_schedules"
}
