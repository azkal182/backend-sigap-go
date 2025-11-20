package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Teacher represents an instructor that can be scheduled into classes.
type Teacher struct {
	ID               uuid.UUID      `json:"id"`
	UserID           *uuid.UUID     `json:"user_id" gorm:"uniqueIndex"`
	TeacherCode      string         `json:"teacher_code" gorm:"size:50;uniqueIndex;not null"`
	FullName         string         `json:"full_name" gorm:"size:150;not null"`
	Gender           string         `json:"gender" gorm:"size:10"`
	Phone            string         `json:"phone" gorm:"size:30"`
	Email            string         `json:"email" gorm:"size:150"`
	Specialization   string         `json:"specialization" gorm:"size:150"`
	EmploymentStatus string         `json:"employment_status" gorm:"size:50"`
	JoinedAt         *time.Time     `json:"joined_at"`
	IsActive         bool           `json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies GORM table name.
func (Teacher) TableName() string {
	return "teachers"
}

// NormalizeUsername derives username from teacher full name.
func (t *Teacher) NormalizeUsername() string {
	username := strings.ToLower(strings.TrimSpace(t.FullName))
	username = strings.ReplaceAll(username, " ", "")
	return username
}
