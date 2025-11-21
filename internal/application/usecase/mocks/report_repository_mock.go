package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
)

// ReportRepositoryMock implements repository.ReportRepository for testing.
type ReportRepositoryMock struct {
	mock.Mock
}

var _ domainRepo.ReportRepository = (*ReportRepositoryMock)(nil)

func (m *ReportRepositoryMock) AggregateStudentAttendance(ctx context.Context, filter domainRepo.StudentAttendanceReportFilter) ([]domainRepo.StudentAttendanceAggregation, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.StudentAttendanceAggregation)
	return rows, args.Error(1)
}

func (m *ReportRepositoryMock) AggregateTeacherAttendance(ctx context.Context, filter domainRepo.TeacherAttendanceReportFilter) ([]domainRepo.TeacherAttendanceAggregation, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.TeacherAttendanceAggregation)
	return rows, args.Error(1)
}

func (m *ReportRepositoryMock) AggregateLeavePermits(ctx context.Context, filter domainRepo.LeavePermitReportFilter) ([]domainRepo.LeavePermitAggregation, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.LeavePermitAggregation)
	return rows, args.Error(1)
}

func (m *ReportRepositoryMock) AggregateHealthStatuses(ctx context.Context, filter domainRepo.HealthStatusReportFilter) ([]domainRepo.HealthStatusAggregation, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.HealthStatusAggregation)
	return rows, args.Error(1)
}

func (m *ReportRepositoryMock) AggregateSKSResults(ctx context.Context, filter domainRepo.SKSReportFilter) ([]domainRepo.SKSAggregation, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.SKSAggregation)
	return rows, args.Error(1)
}

func (m *ReportRepositoryMock) ListMutationHistory(ctx context.Context, filter domainRepo.MutationReportFilter) ([]domainRepo.MutationHistoryRow, error) {
	args := m.Called(ctx, filter)
	rows, _ := args.Get(0).([]domainRepo.MutationHistoryRow)
	return rows, args.Error(1)
}
