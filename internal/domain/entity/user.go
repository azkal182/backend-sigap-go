package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity in the domain
type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username" gorm:"uniqueIndex"`
	Password  string     `json:"-"` // Never expose password in JSON
	Name      string     `json:"name"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Roles       []Role      `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Dormitories []Dormitory `gorm:"many2many:user_dormitories;" json:"dormitories,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the user's hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HasPermission checks if user has a specific permission through their roles
func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission {
				return true
			}
		}
	}
	return false
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// CanAccessDormitory checks if user can access a specific dormitory
// Returns true if user has access to all dormitories or specific dormitory
func (u *User) CanAccessDormitory(dormitoryID uuid.UUID) bool {
	// Check if user has access to all dormitories (via special role or guard)
	for _, role := range u.Roles {
		if role.Name == "admin" || role.Name == "super_admin" {
			return true
		}
	}

	// Check if user has access to specific dormitory
	for _, dorm := range u.Dormitories {
		if dorm.ID == dormitoryID {
			return true
		}
	}

	return false
}
