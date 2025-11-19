package dto

// CreateFanRequest represents payload to create a FAN entity.
type CreateFanRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Level       string `json:"level" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// UpdateFanRequest represents payload to update a FAN entity.
type UpdateFanRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=150"`
	Level       *string `json:"level" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

// FanResponse represents FAN data returned to clients.
type FanResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListFansResponse wraps paginated fans result.
type ListFansResponse struct {
	Fans       []FanResponse `json:"fans"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}
