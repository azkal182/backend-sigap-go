package repository

import (
	"context"

	"github.com/your-org/go-backend-starter/internal/domain/entity"
)

// AuditLogRepository defines the interface for audit log data operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	List(ctx context.Context, filter AuditLogFilter) ([]*entity.AuditLog, int64, error)
}

// AuditLogFilter represents filtering and pagination options for listing audit logs
type AuditLogFilter struct {
	Page       int
	PageSize   int
	Resource   string
	Action     string
	ActorEmail string
}
