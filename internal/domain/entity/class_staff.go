package entity

import (
	"time"

	"github.com/google/uuid"
)

// ClassStaff links staff/users to classes with a specific academic role.
type ClassStaff struct {
	ID        uuid.UUID `json:"id"`
	ClassID   uuid.UUID `json:"class_id" gorm:"index;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"index;not null"`
	Role      string    `json:"role" gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM.
func (ClassStaff) TableName() string {
	return "class_staff"
}
