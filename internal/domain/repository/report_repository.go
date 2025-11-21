package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DateRange captures optional start/end filters reused by report queries.
type DateRange struct {
	Start *time.Time
	End   *time.Time
}

// StudentAttendanceReportFilter defines inputs for aggregating student attendance.
type StudentAttendanceReportFilter struct {
	Date        time.Time
	DormitoryID *uuid.UUID
	ClassID     *uuid.UUID
	FanID       *uuid.UUID
}

// StudentAttendanceAggregation represents counts grouped by dorm/class/FAN.
type StudentAttendanceAggregation struct {
	DormitoryID *uuid.UUID
	ClassID     *uuid.UUID
	FanID       *uuid.UUID
	Total       int
	Present     int
	Absent      int
	Permit      int
	Sick        int
}

// TeacherAttendanceReportFilter defines inputs for aggregating teacher attendance.
type TeacherAttendanceReportFilter struct {
	Date      time.Time
	SlotID    *uuid.UUID
	TeacherID *uuid.UUID
}

// TeacherAttendanceAggregation summarizes teacher attendance metrics.
type TeacherAttendanceAggregation struct {
	TeacherID uuid.UUID
	Total     int
	Present   int
	Absent    int
}

// LeavePermitReportFilter defines filters for leave permit aggregations.
type LeavePermitReportFilter struct {
	Status      *string
	Type        *string
	DormitoryID *uuid.UUID
	DateRange   DateRange
}

// LeavePermitAggregation represents aggregated leave permit counts.
type LeavePermitAggregation struct {
	DormitoryID *uuid.UUID
	Type        string
	Status      string
	Total       int
}

// HealthStatusReportFilter defines filters for health status aggregation.
type HealthStatusReportFilter struct {
	Status      *string
	DormitoryID *uuid.UUID
	DateRange   DateRange
}

// HealthStatusAggregation represents active/revoked counts per dormitory.
type HealthStatusAggregation struct {
	DormitoryID *uuid.UUID
	Status      string
	Total       int
	Consecutive int
}

// SKSReportFilter defines filters for SKS pass-rate aggregation.
type SKSReportFilter struct {
	FanID     *uuid.UUID
	SKSID     *uuid.UUID
	IsPassed  *bool
	DateRange DateRange
}

// SKSAggregation holds pass/fail totals and optional average score.
type SKSAggregation struct {
	FanID        *uuid.UUID
	SKSID        *uuid.UUID
	Total        int
	Passed       int
	Failed       int
	AverageScore *int
}

// MutationReportFilter defines filters for student mutation history reporting.
type MutationReportFilter struct {
	StudentID   *uuid.UUID
	FanID       *uuid.UUID
	DormitoryID *uuid.UUID
	DateRange   DateRange
}

// MutationHistoryRow represents a dorm/class transition window.
type MutationHistoryRow struct {
	StudentID   uuid.UUID
	FromDormID  *uuid.UUID
	ToDormID    *uuid.UUID
	FromClassID *uuid.UUID
	ToClassID   *uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

// ReportRepository declares read-only aggregation queries for Phase 8 reports.
type ReportRepository interface {
	AggregateStudentAttendance(ctx context.Context, filter StudentAttendanceReportFilter) ([]StudentAttendanceAggregation, error)
	AggregateTeacherAttendance(ctx context.Context, filter TeacherAttendanceReportFilter) ([]TeacherAttendanceAggregation, error)
	AggregateLeavePermits(ctx context.Context, filter LeavePermitReportFilter) ([]LeavePermitAggregation, error)
	AggregateHealthStatuses(ctx context.Context, filter HealthStatusReportFilter) ([]HealthStatusAggregation, error)
	AggregateSKSResults(ctx context.Context, filter SKSReportFilter) ([]SKSAggregation, error)
	ListMutationHistory(ctx context.Context, filter MutationReportFilter) ([]MutationHistoryRow, error)
}
