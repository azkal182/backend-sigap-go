package dto

// CreateClassRequest represents payload to create a class under a FAN.
type CreateClassRequest struct {
	FanID    string `json:"fan_id" binding:"required,uuid4"`
	Name     string `json:"name" binding:"required,min=3,max=150"`
	Capacity int    `json:"capacity" binding:"omitempty,min=0"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateClassRequest represents payload to update class attributes.
type UpdateClassRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=3,max=150"`
	Capacity *int    `json:"capacity" binding:"omitempty,min=0"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}

// ClassResponse represents class data returned to clients.
type ClassResponse struct {
	ID        string `json:"id"`
	FanID     string `json:"fan_id"`
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListClassesResponse wraps paginated class result.
type ListClassesResponse struct {
	Classes    []ClassResponse `json:"classes"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// EnrollStudentRequest represents payload to enroll a student into a class.
type EnrollStudentRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid4"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// AssignClassStaffRequest represents payload to assign staff to a class.
type AssignClassStaffRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
	Role   string `json:"role" binding:"required,oneof=class_manager homeroom_teacher"`
}
