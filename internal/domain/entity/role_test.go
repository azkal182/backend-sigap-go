package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRole_HasPermission(t *testing.T) {
	role := &Role{
		ID:   uuid.New(),
		Name: "admin",
		Permissions: []Permission{
			{ID: uuid.New(), Name: "users.read"},
			{ID: uuid.New(), Name: "users.write"},
			{ID: uuid.New(), Name: "users.delete"},
		},
	}

	tests := []struct {
		name           string
		permissionName string
		expectedResult bool
	}{
		{
			name:           "success - has permission",
			permissionName: "users.read",
			expectedResult: true,
		},
		{
			name:           "success - has another permission",
			permissionName: "users.write",
			expectedResult: true,
		},
		{
			name:           "failure - does not have permission",
			permissionName: "dormitories.read",
			expectedResult: false,
		},
		{
			name:           "failure - empty permission",
			permissionName: "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := role.HasPermission(tt.permissionName)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestRole_HasPermission_EmptyPermissions(t *testing.T) {
	role := &Role{
		ID:          uuid.New(),
		Name:        "user",
		Permissions: []Permission{},
	}

	result := role.HasPermission("users.read")
	assert.False(t, result)
}
