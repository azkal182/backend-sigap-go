package dto

// CreateClassScheduleRequest represents payload to create a class schedule entry.
type CreateClassScheduleRequest struct {
	ClassID     string  `json:"class_id" binding:"required,uuid4"`
	DormitoryID string  `json:"dormitory_id" binding:"required,uuid4"`
	SubjectID   *string `json:"subject_id" binding:"omitempty,uuid4"`
	TeacherID   string  `json:"teacher_id" binding:"required,uuid4"`
	SlotID      *string `json:"slot_id" binding:"omitempty,uuid4"`
	DayOfWeek   string  `json:"day_of_week" binding:"required,oneof=mon tue wed thu fri sat sun"`
	StartTime   *string `json:"start_time" binding:"omitempty"`
	EndTime     *string `json:"end_time" binding:"omitempty"`
	Location    string  `json:"location" binding:"omitempty,max=150"`
	Notes       string  `json:"notes" binding:"omitempty,max=255"`
}

// UpdateClassScheduleRequest represents payload to update a schedule entry.
type UpdateClassScheduleRequest struct {
	SubjectID *string `json:"subject_id" binding:"omitempty,uuid4"`
	TeacherID *string `json:"teacher_id" binding:"omitempty,uuid4"`
	SlotID    *string `json:"slot_id" binding:"omitempty,uuid4"`
	DayOfWeek *string `json:"day_of_week" binding:"omitempty,oneof=mon tue wed thu fri sat sun"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
	Location  *string `json:"location" binding:"omitempty,max=150"`
	Notes     *string `json:"notes" binding:"omitempty,max=255"`
	IsActive  *bool   `json:"is_active"`
}

// ClassScheduleResponse represents read model for schedules.
type ClassScheduleResponse struct {
	ID          string  `json:"id"`
	ClassID     string  `json:"class_id"`
	DormitoryID string  `json:"dormitory_id"`
	SubjectID   *string `json:"subject_id"`
	TeacherID   string  `json:"teacher_id"`
	SlotID      *string `json:"slot_id"`
	DayOfWeek   string  `json:"day_of_week"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Location    string  `json:"location"`
	Notes       string  `json:"notes"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ListClassSchedulesResponse wraps paginated schedules.
type ListClassSchedulesResponse struct {
	Schedules  []ClassScheduleResponse `json:"schedules"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}
