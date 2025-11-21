package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// ReportUseCase orchestrates read-only report aggregation flows.
type ReportUseCase struct {
	reportRepo repository.ReportRepository
}

// NewReportUseCase creates a new ReportUseCase instance.
func NewReportUseCase(reportRepo repository.ReportRepository) *ReportUseCase {
	return &ReportUseCase{reportRepo: reportRepo}
}

// GetStudentAttendanceReport aggregates student attendance per filters.
func (uc *ReportUseCase) GetStudentAttendanceReport(ctx context.Context, req dto.StudentAttendanceReportRequest) (*dto.StudentAttendanceReportResponse, error) {
	date, err := parseISODate(req.Date)
	if err != nil {
		return nil, err
	}

	filter := repository.StudentAttendanceReportFilter{
		Date:        date,
		DormitoryID: parseUUIDPtr(req.DormitoryID),
		ClassID:     parseUUIDPtr(req.ClassID),
		FanID:       parseUUIDPtr(req.FanID),
	}

	aggregations, err := uc.reportRepo.AggregateStudentAttendance(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.StudentAttendanceReportRow, 0, len(aggregations))
	for _, agg := range aggregations {
		rows = append(rows, dto.StudentAttendanceReportRow{
			DormitoryID: uuidPtrToString(agg.DormitoryID),
			ClassID:     uuidPtrToString(agg.ClassID),
			FanID:       uuidPtrToString(agg.FanID),
			Total:       agg.Total,
			Present:     agg.Present,
			Absent:      agg.Absent,
			Permit:      agg.Permit,
			Sick:        agg.Sick,
		})
	}

	return &dto.StudentAttendanceReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

// GetTeacherAttendanceReport aggregates teacher punctuality metrics.
func (uc *ReportUseCase) GetTeacherAttendanceReport(ctx context.Context, req dto.TeacherAttendanceReportRequest) (*dto.TeacherAttendanceReportResponse, error) {
	date, err := parseISODate(req.Date)
	if err != nil {
		return nil, err
	}

	filter := repository.TeacherAttendanceReportFilter{
		Date:      date,
		SlotID:    parseUUIDPtr(req.SlotID),
		TeacherID: parseUUIDPtr(req.TeacherID),
	}

	aggregations, err := uc.reportRepo.AggregateTeacherAttendance(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.TeacherAttendanceReportRow, 0, len(aggregations))
	for _, agg := range aggregations {
		rows = append(rows, dto.TeacherAttendanceReportRow{
			TeacherID: agg.TeacherID.String(),
			Total:     agg.Total,
			Present:   agg.Present,
			Absent:    agg.Absent,
		})
	}

	return &dto.TeacherAttendanceReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

// GetLeavePermitReport aggregates leave permit stats per filters.
func (uc *ReportUseCase) GetLeavePermitReport(ctx context.Context, req dto.LeavePermitReportRequest) (*dto.LeavePermitReportResponse, error) {
	dateRange, err := buildDateRange(req.DateRangeFilter)
	if err != nil {
		return nil, err
	}

	filter := repository.LeavePermitReportFilter{
		Status:      req.Status,
		Type:        req.Type,
		DormitoryID: parseUUIDPtr(req.DormitoryID),
		DateRange:   dateRange,
	}

	aggregations, err := uc.reportRepo.AggregateLeavePermits(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.LeavePermitReportRow, 0, len(aggregations))
	for _, agg := range aggregations {
		rows = append(rows, dto.LeavePermitReportRow{
			DormitoryID: uuidPtrToString(agg.DormitoryID),
			Type:        agg.Type,
			Status:      agg.Status,
			Total:       agg.Total,
		})
	}

	return &dto.LeavePermitReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

// GetHealthStatusReport aggregates health status stats per filters.
func (uc *ReportUseCase) GetHealthStatusReport(ctx context.Context, req dto.HealthStatusReportRequest) (*dto.HealthStatusReportResponse, error) {
	dateRange, err := buildDateRange(req.DateRangeFilter)
	if err != nil {
		return nil, err
	}

	filter := repository.HealthStatusReportFilter{
		Status:      req.Status,
		DormitoryID: parseUUIDPtr(req.DormitoryID),
		DateRange:   dateRange,
	}

	aggregations, err := uc.reportRepo.AggregateHealthStatuses(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.HealthStatusReportRow, 0, len(aggregations))
	for _, agg := range aggregations {
		rows = append(rows, dto.HealthStatusReportRow{
			DormitoryID: uuidPtrToString(agg.DormitoryID),
			Status:      agg.Status,
			Total:       agg.Total,
			Consecutive: agg.Consecutive,
		})
	}

	return &dto.HealthStatusReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

// GetSKSReport aggregates SKS pass/fail summaries per filters.
func (uc *ReportUseCase) GetSKSReport(ctx context.Context, req dto.SKSReportRequest) (*dto.SKSReportResponse, error) {
	dateRange, err := buildDateRange(req.DateRangeFilter)
	if err != nil {
		return nil, err
	}

	filter := repository.SKSReportFilter{
		FanID:     parseUUIDPtr(req.FanID),
		SKSID:     parseUUIDPtr(req.SKSID),
		IsPassed:  req.IsPassed,
		DateRange: dateRange,
	}

	aggregations, err := uc.reportRepo.AggregateSKSResults(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.SKSReportRow, 0, len(aggregations))
	for _, agg := range aggregations {
		rows = append(rows, dto.SKSReportRow{
			FanID:   uuidPtrToString(agg.FanID),
			SKSID:   uuidPtrToString(agg.SKSID),
			Total:   agg.Total,
			Passed:  agg.Passed,
			Failed:  agg.Failed,
			Average: agg.AverageScore,
		})
	}

	return &dto.SKSReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

// GetMutationReport lists dorm/class mutation histories per filters.
func (uc *ReportUseCase) GetMutationReport(ctx context.Context, req dto.MutationReportRequest) (*dto.MutationReportResponse, error) {
	dateRange, err := buildDateRange(req.DateRangeFilter)
	if err != nil {
		return nil, err
	}

	filter := repository.MutationReportFilter{
		StudentID:   parseUUIDPtr(req.StudentID),
		FanID:       parseUUIDPtr(req.FanID),
		DormitoryID: parseUUIDPtr(req.DormitoryID),
		DateRange:   dateRange,
	}

	rowsData, err := uc.reportRepo.ListMutationHistory(ctx, filter)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	rows := make([]dto.MutationReportRow, 0, len(rowsData))
	for _, row := range rowsData {
		rows = append(rows, dto.MutationReportRow{
			StudentID:   row.StudentID.String(),
			FromDormID:  uuidPtrToString(row.FromDormID),
			ToDormID:    uuidPtrToString(row.ToDormID),
			FromClassID: uuidPtrToString(row.FromClassID),
			ToClassID:   uuidPtrToString(row.ToClassID),
			StartDate:   row.StartDate.Format("2006-01-02"),
			EndDate:     timePtrToString(row.EndDate),
		})
	}

	return &dto.MutationReportResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Rows:        rows,
		Filters:     req,
	}, nil
}

func parseISODate(value string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, domainErrors.ErrBadRequest
	}
	return date, nil
}

func parseUUIDPtr(val *string) *uuid.UUID {
	if val == nil || *val == "" {
		return nil
	}
	parsed, err := uuid.Parse(*val)
	if err != nil {
		return nil
	}
	return &parsed
}

func buildDateRange(filter dto.DateRangeFilter) (repository.DateRange, error) {
	var dr repository.DateRange
	if filter.StartDate != nil {
		start, err := parseISODate(*filter.StartDate)
		if err != nil {
			return dr, err
		}
		dr.Start = &start
	}
	if filter.EndDate != nil {
		end, err := parseISODate(*filter.EndDate)
		if err != nil {
			return dr, err
		}
		dr.End = &end
	}
	return dr, nil
}
