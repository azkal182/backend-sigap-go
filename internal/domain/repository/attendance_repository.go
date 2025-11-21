package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// AttendanceSessionFilter captures optional filters for listing sessions.
type AttendanceSessionFilter struct {
	ClassScheduleID *uuid.UUID
	TeacherID       *uuid.UUID
	Date            *time.Time
	Status          *entity.AttendanceSessionStatus
	Limit           int
	Offset          int
}

// AttendanceSessionRepository defines persistence for sessions.
type AttendanceSessionRepository interface {
	Create(ctx context.Context, session *entity.AttendanceSession) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.AttendanceSessionStatus, lockedAt *time.Time) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AttendanceSession, error)
	GetOpenByScheduleAndDate(ctx context.Context, scheduleID uuid.UUID, date time.Time) (*entity.AttendanceSession, error)
	List(ctx context.Context, filter AttendanceSessionFilter) ([]*entity.AttendanceSession, int64, error)
	LockSessionsByDate(ctx context.Context, date time.Time) error
}

// StudentAttendanceRepository defines persistence for student attendance rows.
type StudentAttendanceRepository interface {
	BulkUpsert(ctx context.Context, attendances []*entity.StudentAttendance) error
	ListBySession(ctx context.Context, sessionID uuid.UUID) ([]*entity.StudentAttendance, error)
}

// TeacherAttendanceRepository defines persistence for teacher attendance rows.
type TeacherAttendanceRepository interface {
	Upsert(ctx context.Context, attendance *entity.TeacherAttendance) error
	GetBySession(ctx context.Context, sessionID uuid.UUID) (*entity.TeacherAttendance, error)
}
