package dto

// CreateScheduleSlotRequest represents payload to create a slot.
type CreateScheduleSlotRequest struct {
	DormitoryID string `json:"dormitory_id" binding:"required,uuid4"`
	SlotNumber  int    `json:"slot_number" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=3,max=100"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// UpdateScheduleSlotRequest updates a slot.
type UpdateScheduleSlotRequest struct {
	SlotNumber  *int    `json:"slot_number" binding:"omitempty,min=1"`
	Name        *string `json:"name" binding:"omitempty,min=3,max=100"`
	StartTime   *string `json:"start_time" binding:"omitempty"`
	EndTime     *string `json:"end_time" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active"`
}

// ScheduleSlotResponse represents slot data.
type ScheduleSlotResponse struct {
	ID          string `json:"id"`
	DormitoryID string `json:"dormitory_id"`
	SlotNumber  int    `json:"slot_number"`
	Name        string `json:"name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListScheduleSlotsResponse paginated response.
type ListScheduleSlotsResponse struct {
	Slots      []ScheduleSlotResponse `json:"slots"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}
