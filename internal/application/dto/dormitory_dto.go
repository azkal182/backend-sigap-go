package dto

// CreateDormitoryRequest represents the request to create a dormitory
type CreateDormitoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Gender      string `json:"gender" binding:"required,oneof=male female"`
	Level       string `json:"level" binding:"required,min=2,max=50"`
	Code        string `json:"code" binding:"required,alphanumunicode,min=2,max=16"`
	Description string `json:"description"`
}

// UpdateDormitoryRequest represents the request to update a dormitory
type UpdateDormitoryRequest struct {
	Name        string  `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Gender      *string `json:"gender,omitempty" binding:"omitempty,oneof=male female"`
	Level       *string `json:"level,omitempty" binding:"omitempty,min=2,max=50"`
	Code        *string `json:"code,omitempty" binding:"omitempty,alphanumunicode,min=2,max=16"`
	Description string  `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// DormitoryResponse represents dormitory data in responses
type DormitoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Level       string `json:"level"`
	Code        string `json:"code"`
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

// AssignDormitoryUserRequest represents a request to assign a user to a dormitory
type AssignDormitoryUserRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
}
