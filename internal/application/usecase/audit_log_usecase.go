package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// AuditLogUseCase handles read-only audit log operations
type AuditLogUseCase struct {
	repo repository.AuditLogRepository
}

// NewAuditLogUseCase creates a new audit log use case
func NewAuditLogUseCase(repo repository.AuditLogRepository) *AuditLogUseCase {
	return &AuditLogUseCase{repo: repo}
}

// ListAuditLogs retrieves a paginated list of audit logs
func (uc *AuditLogUseCase) ListAuditLogs(ctx context.Context, page, pageSize int, resource, action, actorUsername string) (*dto.ListAuditLogsResponse, error) {
	filter := repository.AuditLogFilter{
		Page:          page,
		PageSize:      pageSize,
		Resource:      resource,
		Action:        action,
		ActorUsername: actorUsername,
	}

	logs, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.AuditLogResponse, 0, len(logs))
	for _, l := range logs {
		var actorIDStr string
		if l.ActorID != nil {
			actorIDStr = l.ActorID.String()
		}

		var roles []string
		if l.ActorRoles != "" {
			_ = json.Unmarshal([]byte(l.ActorRoles), &roles)
		}

		items = append(items, dto.AuditLogResponse{
			ID:            l.ID.String(),
			ActorID:       actorIDStr,
			ActorUsername: l.ActorUsername,
			ActorRoles:    roles,
			Action:        l.Action,
			Resource:      l.Resource,
			TargetID:      l.TargetID,
			RequestPath:   l.RequestPath,
			RequestMethod: l.RequestMethod,
			StatusCode:    l.StatusCode,
			IPAddress:     l.IPAddress,
			UserAgent:     l.UserAgent,
			Metadata:      l.Metadata,
			CreatedAt:     l.CreatedAt.Format(time.RFC3339),
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListAuditLogsResponse{
		Logs:       items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
