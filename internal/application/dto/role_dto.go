package dto

// CreateRoleRequest represents the request to create a role
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Slug        string   `json:"slug" binding:"required"`
	IsActive    bool     `json:"is_active"`
	IsProtected bool     `json:"is_protected"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

// UpdateRoleRequest represents the request to update a role
type UpdateRoleRequest struct {
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// RoleResponse represents role data in responses
type RoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	IsActive    bool     `json:"is_active"`
	IsProtected bool     `json:"is_protected"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// ListRolesResponse represents paginated role list response
type ListRolesResponse struct {
	Roles      []RoleResponse `json:"roles"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// AssignPermissionRequest represents the request to assign a permission to a role
type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id" binding:"required"`
}

// RemovePermissionRequest represents the request to remove a permission from a role
type RemovePermissionRequest struct {
	PermissionID string `json:"permission_id" binding:"required"`
}

// AssignRoleToUserRequest represents the request to assign a role to a user
type AssignRoleToUserRequest struct {
	RoleID string `json:"role_id" binding:"required"`
}
