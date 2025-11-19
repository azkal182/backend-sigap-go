package dto

import "time"

// CreateStudentRequest represents payload to create a student.
type CreateStudentRequest struct {
	StudentNumber string    `json:"student_number" binding:"required,alphanumunicode,min=3,max=50"`
	FullName      string    `json:"full_name" binding:"required,min=3,max=150"`
	BirthDate     time.Time `json:"birth_date" binding:"required"`
	Gender        string    `json:"gender" binding:"required,oneof=male female"`
	ParentName    string    `json:"parent_name" binding:"required,min=3,max=150"`
}

// UpdateStudentRequest represents payload to update student profile.
type UpdateStudentRequest struct {
	FullName   *string    `json:"full_name" binding:"omitempty,min=3,max=150"`
	BirthDate  *time.Time `json:"birth_date" binding:"omitempty"`
	Gender     *string    `json:"gender" binding:"omitempty,oneof=male female"`
	ParentName *string    `json:"parent_name" binding:"omitempty,min=3,max=150"`
}

// UpdateStudentStatusRequest handles status changes.
type UpdateStudentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive leave graduated"`
}

// MutateStudentDormitoryRequest handles dormitory mutation.
type MutateStudentDormitoryRequest struct {
	DormitoryID string    `json:"dormitory_id" binding:"required,uuid4"`
	StartDate   time.Time `json:"start_date" binding:"required"`
}

// StudentResponse represents student response payload.
type StudentResponse struct {
	ID               string                  `json:"id"`
	StudentNumber    string                  `json:"student_number"`
	FullName         string                  `json:"full_name"`
	BirthDate        string                  `json:"birth_date"`
	Gender           string                  `json:"gender"`
	ParentName       string                  `json:"parent_name"`
	Status           string                  `json:"status"`
	IsActive         bool                    `json:"is_active"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
	DormitoryHistory []StudentDormitoryEvent `json:"dormitory_history"`
}

// StudentDormitoryEvent represents dormitory assignment history entry.
type StudentDormitoryEvent struct {
	DormitoryID string `json:"dormitory_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// ListStudentsResponse paginated response.
type ListStudentsResponse struct {
	Students   []StudentResponse `json:"students"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
