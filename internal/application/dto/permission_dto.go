package dto

// PermissionResponse represents permission data in responses
type PermissionResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// ListPermissionsResponse represents paginated permission list response
type ListPermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int                  `json:"total_pages"`
}
