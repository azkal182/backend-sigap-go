package dto

import "time"

// CreateTeacherRequest represents payload to create a teacher.
type CreateTeacherRequest struct {
	TeacherCode      string     `json:"teacher_code" binding:"required,alphanumunicode,min=3,max=50"`
	FullName         string     `json:"full_name" binding:"required,min=3,max=150"`
	Gender           string     `json:"gender" binding:"omitempty,oneof=male female"`
	Phone            string     `json:"phone" binding:"omitempty,max=30"`
	Email            string     `json:"email" binding:"omitempty,email"`
	Specialization   string     `json:"specialization" binding:"omitempty,max=150"`
	EmploymentStatus string     `json:"employment_status" binding:"omitempty,max=50"`
	JoinedAt         *time.Time `json:"joined_at" binding:"omitempty"`
	ExistingUsername string     `json:"existing_username" binding:"omitempty,alphanumunicode,min=3,max=32"`
}

// UpdateTeacherRequest updates teacher profile data.
type UpdateTeacherRequest struct {
	FullName         *string    `json:"full_name" binding:"omitempty,min=3,max=150"`
	Gender           *string    `json:"gender" binding:"omitempty,oneof=male female"`
	Phone            *string    `json:"phone" binding:"omitempty,max=30"`
	Email            *string    `json:"email" binding:"omitempty,email"`
	Specialization   *string    `json:"specialization" binding:"omitempty,max=150"`
	EmploymentStatus *string    `json:"employment_status" binding:"omitempty,max=50"`
	JoinedAt         *time.Time `json:"joined_at" binding:"omitempty"`
	IsActive         *bool      `json:"is_active"`
}

// TeacherResponse represents teacher data for responses.
type TeacherResponse struct {
	ID               string `json:"id"`
	TeacherCode      string `json:"teacher_code"`
	FullName         string `json:"full_name"`
	Gender           string `json:"gender"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Specialization   string `json:"specialization"`
	EmploymentStatus string `json:"employment_status"`
	JoinedAt         string `json:"joined_at"`
	IsActive         bool   `json:"is_active"`
	UserID           string `json:"user_id,omitempty"`
	Username         string `json:"username,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// ListTeachersResponse paginated list of teachers.
type ListTeachersResponse struct {
	Teachers   []TeacherResponse `json:"teachers"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
