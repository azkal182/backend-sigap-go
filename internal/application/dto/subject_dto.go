package dto

// CreateSubjectRequest captures payload for creating a subject.
type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateSubjectRequest captures payload for updating a subject.
type UpdateSubjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
}

// SubjectResponse represents subject data returned to clients.
type SubjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListSubjectsResponse contains paginated subject results.
type ListSubjectsResponse struct {
	Subjects   []SubjectResponse `json:"subjects"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
