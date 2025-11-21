package dto

import "github.com/google/uuid"

// Leave permit DTOs
type CreateLeavePermitRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid4"`
	Type      string `json:"type" binding:"required,oneof=home_leave official_duty"`
	Reason    string `json:"reason" binding:"required"`
	StartDate string `json:"start_date" binding:"required,len=10"`
	EndDate   string `json:"end_date" binding:"required,len=10"`
}

type UpdateLeavePermitStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected completed"`
}

type LeavePermitResponse struct {
	ID         string  `json:"id"`
	StudentID  string  `json:"student_id"`
	Type       string  `json:"type"`
	Reason     string  `json:"reason"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Status     string  `json:"status"`
	CreatedBy  string  `json:"created_by"`
	ApprovedBy *string `json:"approved_by,omitempty"`
	ApprovedAt *string `json:"approved_at,omitempty"`
}

type ListLeavePermitsRequest struct {
	StudentID *string `form:"student_id"`
	Status    *string `form:"status"`
	Type      *string `form:"type"`
	Date      *string `form:"date"`
	Page      int     `form:"page"`
	PageSize  int     `form:"page_size"`
}

type ListLeavePermitsResponse struct {
	Permits    []LeavePermitResponse `json:"permits"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// Health status DTOs
type CreateHealthStatusRequest struct {
	StudentID string  `json:"student_id" binding:"required,uuid4"`
	Diagnosis string  `json:"diagnosis" binding:"required"`
	Notes     string  `json:"notes"`
	StartDate string  `json:"start_date" binding:"required,len=10"`
	EndDate   *string `json:"end_date"`
}

type RevokeHealthStatusRequest struct {
	Reason string `json:"reason"`
}

type HealthStatusResponse struct {
	ID        string  `json:"id"`
	StudentID string  `json:"student_id"`
	Diagnosis string  `json:"diagnosis"`
	Notes     string  `json:"notes"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Status    string  `json:"status"`
	CreatedBy string  `json:"created_by"`
	RevokedBy *string `json:"revoked_by,omitempty"`
	RevokedAt *string `json:"revoked_at,omitempty"`
}

type ListHealthStatusesRequest struct {
	StudentID *string `form:"student_id"`
	Status    *string `form:"status"`
	Date      *string `form:"date"`
	Page      int     `form:"page"`
	PageSize  int     `form:"page_size"`
}

type ListHealthStatusesResponse struct {
	Statuses   []HealthStatusResponse `json:"statuses"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// OptionalUUIDToString converts uuid pointer to string pointer.
func OptionalUUIDToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	str := id.String()
	return &str
}
