package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
)

// Context keys used for audit logging
const (
	CtxKeyRequestPath   = "request_path"
	CtxKeyRequestMethod = "request_method"
	CtxKeyStatusCode    = "status_code"
	CtxKeyIPAddress     = "ip_address"
	CtxKeyUserAgent     = "user_agent"
	CtxKeyActorID       = "user_id"
	CtxKeyActorUsername = "user_username"
	CtxKeyActorRoles    = "user_roles"
)

// AuditLogger defines interface for writing audit logs
type AuditLogger interface {
	Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error
}

type auditLogger struct {
	repo domainRepo.AuditLogRepository
}

// NewAuditLogger creates a new AuditLogger
func NewAuditLogger(repo domainRepo.AuditLogRepository) AuditLogger {
	return &auditLogger{repo: repo}
}

func (l *auditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	// Marshal metadata to JSON (best-effort)
	var metadataStr string
	if len(metadata) > 0 {
		if b, err := json.Marshal(metadata); err == nil {
			metadataStr = string(b)
		}
	}

	// Extract actor info from context
	var actorIDPtr *uuid.UUID
	if v := ctx.Value(CtxKeyActorID); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			actorIDPtr = &id
		}
	}
	actorUsername, _ := ctx.Value(CtxKeyActorUsername).(string)
	actorRolesStr := ""
	if roles, ok := ctx.Value(CtxKeyActorRoles).([]string); ok {
		if b, err := json.Marshal(roles); err == nil {
			actorRolesStr = string(b)
		}
	}

	requestPath, _ := ctx.Value(CtxKeyRequestPath).(string)
	requestMethod, _ := ctx.Value(CtxKeyRequestMethod).(string)
	ipAddress, _ := ctx.Value(CtxKeyIPAddress).(string)
	userAgent, _ := ctx.Value(CtxKeyUserAgent).(string)
	statusCode := 0
	if sc, ok := ctx.Value(CtxKeyStatusCode).(int); ok {
		statusCode = sc
	}

	logEntry := &entity.AuditLog{
		ID:            uuid.New(),
		ActorID:       actorIDPtr,
		ActorUsername: actorUsername,
		ActorRoles:    actorRolesStr,
		Action:        action,
		Resource:      resource,
		TargetID:      targetID,
		RequestPath:   requestPath,
		RequestMethod: requestMethod,
		StatusCode:    statusCode,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Metadata:      metadataStr,
		CreatedAt:     time.Now(),
	}

	// Best-effort logging: if audit log fails, jangan block main flow
	_ = l.repo.Create(ctx, logEntry)
	return nil
}
