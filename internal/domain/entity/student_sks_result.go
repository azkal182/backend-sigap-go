package entity

import (
	"time"

	"github.com/google/uuid"
)

// StudentSKSResult represents a student's outcome for a specific SKS definition.
type StudentSKSResult struct {
	ID         uuid.UUID  `json:"id"`
	StudentID  uuid.UUID  `json:"student_id" gorm:"not null;uniqueIndex:uniq_student_sks"`
	SKSID      uuid.UUID  `json:"sks_id" gorm:"not null;uniqueIndex:uniq_student_sks;index"`
	Score      float64    `json:"score"`
	IsPassed   bool       `json:"is_passed"`
	ExamDate   *time.Time `json:"exam_date"`
	ExaminerID *uuid.UUID `json:"examiner_id" gorm:"index"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName satisfies gorm's table naming for StudentSKSResult.
func (StudentSKSResult) TableName() string {
	return "student_sks_results"
}

// FanCompletionStatus captures whether a student has completed a FAN.
type FanCompletionStatus struct {
	ID          uuid.UUID  `json:"id"`
	StudentID   uuid.UUID  `json:"student_id" gorm:"not null;uniqueIndex:uniq_student_fan"`
	FanID       uuid.UUID  `json:"fan_id" gorm:"not null;uniqueIndex:uniq_student_fan"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName satisfies gorm's table naming for FanCompletionStatus.
func (FanCompletionStatus) TableName() string {
	return "fan_completion_status"
}
