package dto

// CreateDormitoryRequest represents the request to create a dormitory
type CreateDormitoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateDormitoryRequest represents the request to update a dormitory
type UpdateDormitoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// DormitoryResponse represents dormitory data in responses
type DormitoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListDormitoriesResponse represents paginated dormitory list response
type ListDormitoriesResponse struct {
	Dormitories []DormitoryResponse `json:"dormitories"`
	Total       int64               `json:"total"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	TotalPages  int                 `json:"total_pages"`
}
