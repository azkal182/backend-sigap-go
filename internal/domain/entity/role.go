package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role entity in the domain
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	IsActive    bool      `json:"is_active"`
	IsProtected bool      `json:"is_protected"` // Roles that cannot have permissions edited
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName specifies the table name for GORM
func (Role) TableName() string {
	return "roles"
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, perm := range r.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}
