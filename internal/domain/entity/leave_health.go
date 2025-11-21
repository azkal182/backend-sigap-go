package entity

import (
	"time"

	"github.com/google/uuid"
)

// LeavePermitType enumerates supported permit categories.
type LeavePermitType string

const (
	LeavePermitTypeHomeLeave    LeavePermitType = "home_leave"
	LeavePermitTypeOfficialDuty LeavePermitType = "official_duty"
)

// LeavePermitStatus captures workflow transitions for permits.
type LeavePermitStatus string

const (
	LeavePermitStatusPending   LeavePermitStatus = "pending"
	LeavePermitStatusApproved  LeavePermitStatus = "approved"
	LeavePermitStatusRejected  LeavePermitStatus = "rejected"
	LeavePermitStatusCompleted LeavePermitStatus = "completed"
)

// LeavePermit models the security leave permit lifecycle.
type LeavePermit struct {
	ID         uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID  uuid.UUID         `json:"student_id" gorm:"type:char(36);not null;index"`
	Student    *Student          `json:"-" gorm:"foreignKey:StudentID"`
	Type       LeavePermitType   `json:"type" gorm:"type:varchar(32);not null"`
	Reason     string            `json:"reason" gorm:"type:varchar(255)"`
	StartDate  time.Time         `json:"start_date" gorm:"type:date;not null;index"`
	EndDate    time.Time         `json:"end_date" gorm:"type:date;not null;index"`
	Status     LeavePermitStatus `json:"status" gorm:"type:varchar(32);not null;index"`
	CreatedBy  uuid.UUID         `json:"created_by" gorm:"type:char(36);not null"`
	ApprovedBy *uuid.UUID        `json:"approved_by" gorm:"type:char(36)"`
	ApprovedAt *time.Time        `json:"approved_at"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// TableName overrides GORM default.
func (LeavePermit) TableName() string {
	return "leave_permits"
}

// HealthStatusState enumerates sick status lifecycle values.
type HealthStatusState string

const (
	HealthStatusStateActive  HealthStatusState = "active"
	HealthStatusStateRevoked HealthStatusState = "revoked"
)

// HealthStatus tracks UKS sick statuses that impact attendance.
type HealthStatus struct {
	ID        uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID uuid.UUID         `json:"student_id" gorm:"type:char(36);not null;index"`
	Student   *Student          `json:"-" gorm:"foreignKey:StudentID"`
	Diagnosis string            `json:"diagnosis" gorm:"type:varchar(255);not null"`
	Notes     string            `json:"notes" gorm:"type:text"`
	StartDate time.Time         `json:"start_date" gorm:"type:date;not null;index"`
	EndDate   *time.Time        `json:"end_date" gorm:"type:date"`
	Status    HealthStatusState `json:"status" gorm:"type:varchar(32);not null;index"`
	CreatedBy uuid.UUID         `json:"created_by" gorm:"type:char(36);not null"`
	RevokedBy *uuid.UUID        `json:"revoked_by" gorm:"type:char(36)"`
	RevokedAt *time.Time        `json:"revoked_at"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// TableName overrides GORM default.
func (HealthStatus) TableName() string {
	return "health_statuses"
}
