package entity

import (
	"time"

	"github.com/google/uuid"
)

// StudentClassEnrollment tracks a student's enrollment in a class over time.
type StudentClassEnrollment struct {
	ID         uuid.UUID  `json:"id"`
	ClassID    uuid.UUID  `json:"class_id" gorm:"index;not null"`
	StudentID  uuid.UUID  `json:"student_id" gorm:"index;not null"`
	EnrolledAt time.Time  `json:"enrolled_at"`
	LeftAt     *time.Time `json:"left_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName specifies the table name for GORM.
func (StudentClassEnrollment) TableName() string {
	return "student_class_enrollments"
}
