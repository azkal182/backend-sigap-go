package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_HashPassword(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "success - hash password",
			user: &User{
				ID:       uuid.New(),
				Username: "test",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "success - empty password",
			user: &User{
				ID:       uuid.New(),
				Username: "test",
				Password: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPassword := tt.user.Password
			err := tt.user.HashPassword()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, originalPassword, tt.user.Password)
				assert.NotEmpty(t, tt.user.Password)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{
		ID:       uuid.New(),
		Username: "test",
		Password: "password123",
	}

	// Hash the password first
	err := user.HashPassword()
	require.NoError(t, err)

	tests := []struct {
		name           string
		password       string
		expectedResult bool
	}{
		{
			name:           "success - correct password",
			password:       "password123",
			expectedResult: true,
		},
		{
			name:           "failure - incorrect password",
			password:       "wrongpassword",
			expectedResult: false,
		},
		{
			name:           "failure - empty password",
			password:       "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.CheckPassword(tt.password)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	user := &User{
		ID:       uuid.New(),
		Username: "test",
		Roles: []Role{
			{
				ID:   uuid.New(),
				Name: "admin",
				Permissions: []Permission{
					{ID: uuid.New(), Name: "users.read"},
					{ID: uuid.New(), Name: "users.write"},
				},
			},
			{
				ID:   uuid.New(),
				Name: "user",
				Permissions: []Permission{
					{ID: uuid.New(), Name: "dormitories.read"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		permission     string
		expectedResult bool
	}{
		{
			name:           "success - has permission",
			permission:     "users.read",
			expectedResult: true,
		},
		{
			name:           "success - has permission from different role",
			permission:     "dormitories.read",
			expectedResult: true,
		},
		{
			name:           "failure - does not have permission",
			permission:     "users.delete",
			expectedResult: false,
		},
		{
			name:           "failure - empty permission",
			permission:     "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasPermission(tt.permission)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	user := &User{
		ID:       uuid.New(),
		Username: "test",
		Roles: []Role{
			{ID: uuid.New(), Name: "admin"},
			{ID: uuid.New(), Name: "user"},
		},
	}

	tests := []struct {
		name           string
		roleName       string
		expectedResult bool
	}{
		{
			name:           "success - has role",
			roleName:       "admin",
			expectedResult: true,
		},
		{
			name:           "success - has another role",
			roleName:       "user",
			expectedResult: true,
		},
		{
			name:           "failure - does not have role",
			roleName:       "super_admin",
			expectedResult: false,
		},
		{
			name:           "failure - empty role",
			roleName:       "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasRole(tt.roleName)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUser_CanAccessDormitory(t *testing.T) {
	dormitoryID := uuid.New()
	otherDormitoryID := uuid.New()

	tests := []struct {
		name           string
		user           *User
		dormitoryID    uuid.UUID
		expectedResult bool
	}{
		{
			name: "success - admin can access any dormitory",
			user: &User{
				ID:       uuid.New(),
				Username: "admin",
				Roles: []Role{
					{ID: uuid.New(), Name: "admin"},
				},
			},
			dormitoryID:    dormitoryID,
			expectedResult: true,
		},
		{
			name: "success - super_admin can access any dormitory",
			user: &User{
				ID:       uuid.New(),
				Username: "superadmin",
				Roles: []Role{
					{ID: uuid.New(), Name: "super_admin"},
				},
			},
			dormitoryID:    dormitoryID,
			expectedResult: true,
		},
		{
			name: "success - user can access assigned dormitory",
			user: &User{
				ID:       uuid.New(),
				Username: "user",
				Roles: []Role{
					{ID: uuid.New(), Name: "user"},
				},
				Dormitories: []Dormitory{
					{ID: dormitoryID},
				},
			},
			dormitoryID:    dormitoryID,
			expectedResult: true,
		},
		{
			name: "failure - user cannot access unassigned dormitory",
			user: &User{
				ID:       uuid.New(),
				Username: "user",
				Roles: []Role{
					{ID: uuid.New(), Name: "user"},
				},
				Dormitories: []Dormitory{
					{ID: dormitoryID},
				},
			},
			dormitoryID:    otherDormitoryID,
			expectedResult: false,
		},
		{
			name: "failure - user with no roles or dormitories",
			user: &User{
				ID:       uuid.New(),
				Username: "user",
			},
			dormitoryID:    dormitoryID,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.CanAccessDormitory(tt.dormitoryID)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
