package dto

// DateRangeFilter represents optional start/end filters shared by reports.
type DateRangeFilter struct {
	StartDate *string `json:"start_date" form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string `json:"end_date" form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

// StudentAttendanceReportRequest aggregates student attendance per day/class/dormitory.
type StudentAttendanceReportRequest struct {
	Date        string  `json:"date" form:"date" binding:"required,datetime=2006-01-02"`
	DormitoryID *string `json:"dormitory_id" form:"dormitory_id" binding:"omitempty,uuid4"`
	ClassID     *string `json:"class_id" form:"class_id" binding:"omitempty,uuid4"`
	FanID       *string `json:"fan_id" form:"fan_id" binding:"omitempty,uuid4"`
}

// StudentAttendanceReportRow represents aggregated counts for a cohort.
type StudentAttendanceReportRow struct {
	DormitoryID *string `json:"dormitory_id,omitempty"`
	ClassID     *string `json:"class_id,omitempty"`
	FanID       *string `json:"fan_id,omitempty"`
	Total       int     `json:"total"`
	Present     int     `json:"present"`
	Absent      int     `json:"absent"`
	Permit      int     `json:"permit"`
	Sick        int     `json:"sick"`
}

// StudentAttendanceReportResponse wraps rows plus metadata for the student attendance report.
type StudentAttendanceReportResponse struct {
	GeneratedAt string                         `json:"generated_at"`
	Rows        []StudentAttendanceReportRow   `json:"rows"`
	Filters     StudentAttendanceReportRequest `json:"filters"`
}

// TeacherAttendanceReportRequest aggregates teacher punctuality by date/slot.
type TeacherAttendanceReportRequest struct {
	Date      string  `json:"date" form:"date" binding:"required,datetime=2006-01-02"`
	SlotID    *string `json:"slot_id" form:"slot_id" binding:"omitempty,uuid4"`
	TeacherID *string `json:"teacher_id" form:"teacher_id" binding:"omitempty,uuid4"`
}

// TeacherAttendanceReportRow summarizes teacher attendance metrics.
type TeacherAttendanceReportRow struct {
	TeacherID string `json:"teacher_id"`
	Total     int    `json:"total"`
	Present   int    `json:"present"`
	Absent    int    `json:"absent"`
}

// TeacherAttendanceReportResponse represents teacher attendance aggregation output.
type TeacherAttendanceReportResponse struct {
	GeneratedAt string                         `json:"generated_at"`
	Rows        []TeacherAttendanceReportRow   `json:"rows"`
	Filters     TeacherAttendanceReportRequest `json:"filters"`
}

// LeavePermitReportRequest filters leave permit aggregations.
type LeavePermitReportRequest struct {
	Status      *string `json:"status" form:"status" binding:"omitempty,oneof=pending approved rejected completed"`
	Type        *string `json:"type" form:"type" binding:"omitempty,oneof=home_leave official_duty"`
	DormitoryID *string `json:"dormitory_id" form:"dormitory_id" binding:"omitempty,uuid4"`
	DateRangeFilter
}

// LeavePermitReportRow summarises leave permit counts per grouping.
type LeavePermitReportRow struct {
	DormitoryID *string `json:"dormitory_id,omitempty"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Total       int     `json:"total"`
}

// LeavePermitReportResponse contains aggregated leave permit data.
type LeavePermitReportResponse struct {
	GeneratedAt string                   `json:"generated_at"`
	Rows        []LeavePermitReportRow   `json:"rows"`
	Filters     LeavePermitReportRequest `json:"filters"`
}

// HealthStatusReportRequest filters health status reporting data.
type HealthStatusReportRequest struct {
	Status      *string `json:"status" form:"status" binding:"omitempty,oneof=active revoked"`
	DormitoryID *string `json:"dormitory_id" form:"dormitory_id" binding:"omitempty,uuid4"`
	DateRangeFilter
}

// HealthStatusReportRow summarises health status counts.
type HealthStatusReportRow struct {
	DormitoryID *string `json:"dormitory_id,omitempty"`
	Status      string  `json:"status"`
	Total       int     `json:"total"`
	Consecutive int     `json:"consecutive"`
}

// HealthStatusReportResponse contains aggregated health status data.
type HealthStatusReportResponse struct {
	GeneratedAt string                    `json:"generated_at"`
	Rows        []HealthStatusReportRow   `json:"rows"`
	Filters     HealthStatusReportRequest `json:"filters"`
}

// SKSReportRequest filters SKS result reporting data.
type SKSReportRequest struct {
	FanID    *string `json:"fan_id" form:"fan_id" binding:"omitempty,uuid4"`
	SKSID    *string `json:"sks_id" form:"sks_id" binding:"omitempty,uuid4"`
	IsPassed *bool   `json:"is_passed" form:"is_passed"`
	DateRangeFilter
}

// SKSReportRow summarises SKS pass/fail counts.
type SKSReportRow struct {
	FanID   *string `json:"fan_id,omitempty"`
	SKSID   *string `json:"sks_id,omitempty"`
	Total   int     `json:"total"`
	Passed  int     `json:"passed"`
	Failed  int     `json:"failed"`
	Average *int    `json:"average_score,omitempty"`
}

// SKSReportResponse contains aggregated SKS result data.
type SKSReportResponse struct {
	GeneratedAt string           `json:"generated_at"`
	Rows        []SKSReportRow   `json:"rows"`
	Filters     SKSReportRequest `json:"filters"`
}

// MutationReportRequest filters student mutation histories.
type MutationReportRequest struct {
	StudentID   *string `json:"student_id" form:"student_id" binding:"omitempty,uuid4"`
	FanID       *string `json:"fan_id" form:"fan_id" binding:"omitempty,uuid4"`
	DormitoryID *string `json:"dormitory_id" form:"dormitory_id" binding:"omitempty,uuid4"`
	DateRangeFilter
}

// MutationReportRow represents a change event between dorm/class assignments.
type MutationReportRow struct {
	StudentID   string  `json:"student_id"`
	FromDormID  *string `json:"from_dormitory_id"`
	ToDormID    *string `json:"to_dormitory_id"`
	FromClassID *string `json:"from_class_id"`
	ToClassID   *string `json:"to_class_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

// MutationReportResponse contains mutation history data.
type MutationReportResponse struct {
	GeneratedAt string                `json:"generated_at"`
	Rows        []MutationReportRow   `json:"rows"`
	Filters     MutationReportRequest `json:"filters"`
}
