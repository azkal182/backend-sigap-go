package dto

// CreateSKSDefinitionRequest represents payload for creating SKS definition entries.
type CreateSKSDefinitionRequest struct {
	FanID       string  `json:"fan_id" binding:"required,uuid4"`
	SubjectID   *string `json:"subject_id" binding:"omitempty,uuid4"`
	Code        string  `json:"code" binding:"required,min=2,max=50"`
	Name        string  `json:"name" binding:"required,min=3,max=150"`
	KKM         float64 `json:"kkm" binding:"gte=0"`
	Description string  `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
}

// UpdateSKSDefinitionRequest updates SKS definitions.
type UpdateSKSDefinitionRequest struct {
	SubjectID   *string  `json:"subject_id" binding:"omitempty,uuid4"`
	Name        *string  `json:"name" binding:"omitempty,min=3,max=150"`
	KKM         *float64 `json:"kkm" binding:"omitempty,gte=0"`
	Description *string  `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool    `json:"is_active" binding:"omitempty"`
}

// SKSDefinitionResponse represents SKS definition payloads.
type SKSDefinitionResponse struct {
	ID          string  `json:"id"`
	FanID       string  `json:"fan_id"`
	SubjectID   *string `json:"subject_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	KKM         float64 `json:"kkm"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ListSKSDefinitionsResponse wraps paginated SKS definition results.
type ListSKSDefinitionsResponse struct {
	Definitions []SKSDefinitionResponse `json:"definitions"`
	Total       int64                   `json:"total"`
	Page        int                     `json:"page"`
	PageSize    int                     `json:"page_size"`
	TotalPages  int                     `json:"total_pages"`
}

// CreateSKSExamScheduleRequest payload to create SKS exam schedule entries.
type CreateSKSExamScheduleRequest struct {
	SKSID      string  `json:"sks_id" binding:"required,uuid4"`
	ExaminerID *string `json:"examiner_id" binding:"omitempty,uuid4"`
	ExamDate   string  `json:"exam_date" binding:"required"`
	ExamTime   string  `json:"exam_time" binding:"required"`
	Location   string  `json:"location" binding:"omitempty,max=150"`
	Notes      string  `json:"notes" binding:"omitempty,max=255"`
}

// UpdateSKSExamScheduleRequest updates SKS exam schedule entries.
type UpdateSKSExamScheduleRequest struct {
	ExaminerID *string `json:"examiner_id" binding:"omitempty,uuid4"`
	ExamDate   *string `json:"exam_date" binding:"omitempty"`
	ExamTime   *string `json:"exam_time" binding:"omitempty"`
	Location   *string `json:"location" binding:"omitempty,max=150"`
	Notes      *string `json:"notes" binding:"omitempty,max=255"`
}

// SKSExamScheduleResponse represents exam schedule payloads.
type SKSExamScheduleResponse struct {
	ID         string  `json:"id"`
	SKSID      string  `json:"sks_id"`
	ExaminerID *string `json:"examiner_id"`
	ExamDate   string  `json:"exam_date"`
	ExamTime   string  `json:"exam_time"`
	Location   string  `json:"location"`
	Notes      string  `json:"notes"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// ListSKSExamSchedulesResponse wraps paginated exam schedule results.
type ListSKSExamSchedulesResponse struct {
	Exams      []SKSExamScheduleResponse `json:"exams"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalPages int                       `json:"total_pages"`
}
