package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry for sensitive operations
type AuditLog struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ActorID       *uuid.UUID `json:"actor_id" gorm:"type:uuid"`
	ActorEmail    string     `json:"actor_email" gorm:"size:255"`
	ActorRoles    string     `json:"actor_roles" gorm:"type:text"`
	Action        string     `json:"action" gorm:"size:100;index"`
	Resource      string     `json:"resource" gorm:"size:100;index"`
	TargetID      string     `json:"target_id" gorm:"size:255;index"`
	RequestPath   string     `json:"request_path" gorm:"size:255"`
	RequestMethod string     `json:"request_method" gorm:"size:10"`
	StatusCode    int        `json:"status_code"`
	IPAddress     string     `json:"ip_address" gorm:"size:100"`
	UserAgent     string     `json:"user_agent" gorm:"size:512"`
	Metadata      string     `json:"metadata" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
