package entity

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceSessionStatus represents lifecycle of an attendance session.
type AttendanceSessionStatus string

const (
	AttendanceSessionStatusOpen      AttendanceSessionStatus = "open"
	AttendanceSessionStatusSubmitted AttendanceSessionStatus = "submitted"
	AttendanceSessionStatusLocked    AttendanceSessionStatus = "locked"
)

// StudentAttendanceStatus enumerates student attendance outcomes.
type StudentAttendanceStatus string

const (
	StudentAttendancePresent StudentAttendanceStatus = "present"
	StudentAttendanceAbsent  StudentAttendanceStatus = "absent"
	StudentAttendancePermit  StudentAttendanceStatus = "permit"
	StudentAttendanceSick    StudentAttendanceStatus = "sick"
)

// TeacherAttendanceStatus enumerates teacher attendance outcomes.
type TeacherAttendanceStatus string

const (
	TeacherAttendancePresent TeacherAttendanceStatus = "present"
	TeacherAttendanceAbsent  TeacherAttendanceStatus = "absent"
)

// AttendanceSession captures a session generated from class schedules.
type AttendanceSession struct {
	ID              uuid.UUID               `json:"id"`
	ClassScheduleID uuid.UUID               `json:"class_schedule_id" gorm:"not null;index;uniqueIndex:uniq_schedule_date"`
	Date            time.Time               `json:"date" gorm:"type:date;not null;index;uniqueIndex:uniq_schedule_date"`
	StartTime       *time.Time              `json:"start_time"`
	EndTime         *time.Time              `json:"end_time"`
	TeacherID       uuid.UUID               `json:"teacher_id" gorm:"not null;index"`
	Status          AttendanceSessionStatus `json:"status" gorm:"size:20;not null;index"`
	LockedAt        *time.Time              `json:"locked_at"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`

	StudentAttendances []StudentAttendance `json:"student_attendances" gorm:"foreignKey:AttendanceSessionID"`
	TeacherAttendances []TeacherAttendance `json:"teacher_attendances" gorm:"foreignKey:AttendanceSessionID"`
}

// TableName overrides gorm table name.
func (AttendanceSession) TableName() string {
	return "attendance_sessions"
}

// StudentAttendance stores each student's attendance entry per session.
type StudentAttendance struct {
	ID                  uuid.UUID               `json:"id"`
	AttendanceSessionID uuid.UUID               `json:"attendance_session_id" gorm:"not null;index;uniqueIndex:uniq_student_session"`
	StudentID           uuid.UUID               `json:"student_id" gorm:"not null;index;uniqueIndex:uniq_student_session"`
	Status              StudentAttendanceStatus `json:"status" gorm:"size:20;not null"`
	Note                string                  `json:"note" gorm:"size:255"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

// TableName overrides gorm table name.
func (StudentAttendance) TableName() string {
	return "student_attendances"
}

// TeacherAttendance stores the teacher presence per session.
type TeacherAttendance struct {
	ID                  uuid.UUID               `json:"id"`
	AttendanceSessionID uuid.UUID               `json:"attendance_session_id" gorm:"not null;uniqueIndex:uniq_teacher_session"`
	TeacherID           uuid.UUID               `json:"teacher_id" gorm:"not null;index"`
	Status              TeacherAttendanceStatus `json:"status" gorm:"size:20;not null"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

// TableName overrides gorm table name.
func (TeacherAttendance) TableName() string {
	return "teacher_attendances"
}
