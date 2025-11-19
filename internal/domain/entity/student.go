package entity

import (
	"time"

	"github.com/google/uuid"
)

// StudentStatus enumerates possible lifecycle states for a student.
const (
	StudentStatusActive    = "active"
	StudentStatusInactive  = "inactive"
	StudentStatusLeave     = "leave"
	StudentStatusGraduated = "graduated"
)

// Student represents a student entity in the domain.
type Student struct {
	ID            uuid.UUID `json:"id"`
	StudentNumber string    `json:"student_number" gorm:"size:50;uniqueIndex;not null"`
	FullName      string    `json:"full_name" gorm:"size:150;not null"`
	BirthDate     time.Time `json:"birth_date"`
	Gender        string    `json:"gender" gorm:"size:10;not null"`
	ParentName    string    `json:"parent_name" gorm:"size:150"`
	Status        string    `json:"status" gorm:"size:20;default:'active'"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	DormitoryHistories []StudentDormitoryHistory `json:"dormitory_histories,omitempty"`
}

// TableName defines students table name.
func (Student) TableName() string {
	return "students"
}

// StudentDormitoryHistory captures student dorm mutations over time.
type StudentDormitoryHistory struct {
	ID          uuid.UUID  `json:"id"`
	StudentID   uuid.UUID  `json:"student_id" gorm:"index;not null"`
	DormitoryID uuid.UUID  `json:"dormitory_id" gorm:"index;not null"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Student   *Student   `json:"student,omitempty"`
	Dormitory *Dormitory `json:"dormitory,omitempty"`
}

// TableName defines history table name.
func (StudentDormitoryHistory) TableName() string {
	return "student_dormitory_history"
}
