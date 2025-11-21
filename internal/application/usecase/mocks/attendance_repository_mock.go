package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// AttendanceSessionRepositoryMock mocks attendance session persistence.
type AttendanceSessionRepositoryMock struct {
	mock.Mock
}

var _ repository.AttendanceSessionRepository = (*AttendanceSessionRepositoryMock)(nil)

func (m *AttendanceSessionRepositoryMock) Create(ctx context.Context, session *entity.AttendanceSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *AttendanceSessionRepositoryMock) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.AttendanceSessionStatus, lockedAt *time.Time) error {
	args := m.Called(ctx, id, status, lockedAt)
	return args.Error(0)
}

func (m *AttendanceSessionRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.AttendanceSession, error) {
	args := m.Called(ctx, id)
	if session, ok := args.Get(0).(*entity.AttendanceSession); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttendanceSessionRepositoryMock) GetOpenByScheduleAndDate(ctx context.Context, scheduleID uuid.UUID, date time.Time) (*entity.AttendanceSession, error) {
	args := m.Called(ctx, scheduleID, date)
	if session, ok := args.Get(0).(*entity.AttendanceSession); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttendanceSessionRepositoryMock) List(ctx context.Context, filter repository.AttendanceSessionFilter) ([]*entity.AttendanceSession, int64, error) {
	args := m.Called(ctx, filter)
	sessions, _ := args.Get(0).([]*entity.AttendanceSession)
	var total int64
	if val, ok := args.Get(1).(int64); ok {
		total = val
	}
	return sessions, total, args.Error(2)
}

func (m *AttendanceSessionRepositoryMock) LockSessionsByDate(ctx context.Context, date time.Time) error {
	args := m.Called(ctx, date)
	return args.Error(0)
}

// StudentAttendanceRepositoryMock mocks student attendance persistence.
type StudentAttendanceRepositoryMock struct {
	mock.Mock
}

var _ repository.StudentAttendanceRepository = (*StudentAttendanceRepositoryMock)(nil)

func (m *StudentAttendanceRepositoryMock) BulkUpsert(ctx context.Context, attendances []*entity.StudentAttendance) error {
	args := m.Called(ctx, attendances)
	return args.Error(0)
}

func (m *StudentAttendanceRepositoryMock) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]*entity.StudentAttendance, error) {
	args := m.Called(ctx, sessionID)
	records, _ := args.Get(0).([]*entity.StudentAttendance)
	return records, args.Error(1)
}

// TeacherAttendanceRepositoryMock mocks teacher attendance persistence.
type TeacherAttendanceRepositoryMock struct {
	mock.Mock
}

var _ repository.TeacherAttendanceRepository = (*TeacherAttendanceRepositoryMock)(nil)

func (m *TeacherAttendanceRepositoryMock) Upsert(ctx context.Context, attendance *entity.TeacherAttendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *TeacherAttendanceRepositoryMock) GetBySession(ctx context.Context, sessionID uuid.UUID) (*entity.TeacherAttendance, error) {
	args := m.Called(ctx, sessionID)
	if record, ok := args.Get(0).(*entity.TeacherAttendance); ok {
		return record, args.Error(1)
	}
	return nil, args.Error(1)
}
