package dto

// OpenAttendanceSessionRequest triggers automatic session generation for schedules on a date.
type OpenAttendanceSessionRequest struct {
	ClassScheduleIDs []string `json:"class_schedule_ids" binding:"required,dive,uuid4"`
	Date             string   `json:"date" binding:"required"`
}

// SubmitStudentAttendanceRequest bulk submits student attendance for a session.
type SubmitStudentAttendanceRequest struct {
	Records []StudentAttendanceRecord `json:"records" binding:"required,dive"`
}

// StudentAttendanceRecord represents a single student's attendance payload.
type StudentAttendanceRecord struct {
	StudentID string `json:"student_id" binding:"required,uuid4"`
	Status    string `json:"status" binding:"required,oneof=present absent permit sick"`
	Note      string `json:"note" binding:"omitempty,max=255"`
}

// SubmitTeacherAttendanceRequest records teacher presence for a session.
type SubmitTeacherAttendanceRequest struct {
	TeacherID string `json:"teacher_id" binding:"required,uuid4"`
	Status    string `json:"status" binding:"required,oneof=present absent"`
}

// LockAttendanceRequest locks sessions for a specific date.
type LockAttendanceRequest struct {
	Date string `json:"date" binding:"required"`
}

// ListAttendanceSessionsRequest defines filters for listing attendance sessions.
type ListAttendanceSessionsRequest struct {
	ClassScheduleID *string `json:"class_schedule_id"`
	TeacherID       *string `json:"teacher_id"`
	Date            *string `json:"date"`
	Status          *string `json:"status"`
	Page            int     `json:"page"`
	PageSize        int     `json:"page_size"`
}

// AttendanceSessionResponse represents session data returned to clients.
type AttendanceSessionResponse struct {
	ID              string                            `json:"id"`
	ClassScheduleID string                            `json:"class_schedule_id"`
	Date            string                            `json:"date"`
	StartTime       *string                           `json:"start_time"`
	EndTime         *string                           `json:"end_time"`
	TeacherID       string                            `json:"teacher_id"`
	Status          string                            `json:"status"`
	LockedAt        *string                           `json:"locked_at"`
	StudentRecords  []StudentAttendanceRecordResponse `json:"student_records"`
	TeacherRecord   *TeacherAttendanceRecordResponse  `json:"teacher_record"`
}

// StudentAttendanceRecordResponse represents a student's attendance record.
type StudentAttendanceRecordResponse struct {
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

// TeacherAttendanceRecordResponse represents teacher attendance record.
type TeacherAttendanceRecordResponse struct {
	TeacherID string `json:"teacher_id"`
	Status    string `json:"status"`
}

// ListAttendanceSessionsResponse wraps paginated attendance sessions.
type ListAttendanceSessionsResponse struct {
	Sessions   []AttendanceSessionResponse `json:"sessions"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}
