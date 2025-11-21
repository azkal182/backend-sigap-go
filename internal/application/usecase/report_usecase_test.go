package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

func TestReportUseCase_GetStudentAttendanceReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	date := "2025-11-21"
	filter := repository.StudentAttendanceReportFilter{Date: mustParseDate(t, date)}
	rows := []repository.StudentAttendanceAggregation{{
		DormitoryID: uuidPtr(uuid.New()),
		ClassID:     uuidPtr(uuid.New()),
		FanID:       uuidPtr(uuid.New()),
		Total:       30,
		Present:     25,
		Absent:      3,
		Permit:      1,
		Sick:        1,
	}}
	repo.On("AggregateStudentAttendance", ctx, filter).Return(rows, nil)

	resp, err := uc.GetStudentAttendanceReport(ctx, dto.StudentAttendanceReportRequest{Date: date})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, 30, resp.Rows[0].Total)
	repo.AssertExpectations(t)
}

func TestReportUseCase_GetTeacherAttendanceReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	date := "2025-11-21"
	filter := repository.TeacherAttendanceReportFilter{Date: mustParseDate(t, date)}
	rows := []repository.TeacherAttendanceAggregation{{
		TeacherID: uuid.New(),
		Total:     5,
		Present:   4,
		Absent:    1,
	}}
	repo.On("AggregateTeacherAttendance", ctx, filter).Return(rows, nil)

	resp, err := uc.GetTeacherAttendanceReport(ctx, dto.TeacherAttendanceReportRequest{Date: date})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, 4, resp.Rows[0].Present)
	repo.AssertExpectations(t)
}

func TestReportUseCase_GetLeavePermitReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	filter := repository.LeavePermitReportFilter{}
	rows := []repository.LeavePermitAggregation{{Type: "home_leave", Status: "approved", Total: 3}}
	repo.On("AggregateLeavePermits", ctx, filter).Return(rows, nil)

	resp, err := uc.GetLeavePermitReport(ctx, dto.LeavePermitReportRequest{})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, "home_leave", resp.Rows[0].Type)
	repo.AssertExpectations(t)
}

func TestReportUseCase_GetHealthStatusReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	filter := repository.HealthStatusReportFilter{}
	rows := []repository.HealthStatusAggregation{{Status: "active", Total: 2, Consecutive: 4}}
	repo.On("AggregateHealthStatuses", ctx, filter).Return(rows, nil)

	resp, err := uc.GetHealthStatusReport(ctx, dto.HealthStatusReportRequest{})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, 4, resp.Rows[0].Consecutive)
	repo.AssertExpectations(t)
}

func TestReportUseCase_GetSKSReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	filter := repository.SKSReportFilter{}
	avg := 85
	rows := []repository.SKSAggregation{{Total: 10, Passed: 8, Failed: 2, AverageScore: &avg}}
	repo.On("AggregateSKSResults", ctx, filter).Return(rows, nil)

	resp, err := uc.GetSKSReport(ctx, dto.SKSReportRequest{})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, 8, resp.Rows[0].Passed)
	assert.NotNil(t, resp.Rows[0].Average)
	repo.AssertExpectations(t)
}

func TestReportUseCase_GetMutationReport(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)
	ctx := context.Background()
	filter := repository.MutationReportFilter{}
	now := time.Now()
	rows := []repository.MutationHistoryRow{{
		StudentID: uuid.New(),
		StartDate: now,
		EndDate:   &now,
	}}
	repo.On("ListMutationHistory", ctx, filter).Return(rows, nil)

	resp, err := uc.GetMutationReport(ctx, dto.MutationReportRequest{})
	assert.NoError(t, err)
	require.Len(t, resp.Rows, 1)
	assert.Equal(t, now.Format("2006-01-02"), resp.Rows[0].StartDate)
	repo.AssertExpectations(t)
}

func TestReportUseCase_InvalidDate(t *testing.T) {
	repo := new(mocks.ReportRepositoryMock)
	uc := NewReportUseCase(repo)

	resp, err := uc.GetStudentAttendanceReport(context.Background(), dto.StudentAttendanceReportRequest{Date: "invalid"})
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func mustParseDate(t *testing.T, val string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", val)
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}
	return parsed
}

func uuidPtr(id uuid.UUID) *uuid.UUID {
	return &id
}
