package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// inMemoryAuditLogRepo is a simple in-memory implementation of AuditLogRepository for testing
type inMemoryAuditLogRepo struct {
	logs []*entity.AuditLog
}

func (r *inMemoryAuditLogRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	r.logs = append(r.logs, log)
	return nil
}

func (r *inMemoryAuditLogRepo) List(ctx context.Context, filter repository.AuditLogFilter) ([]*entity.AuditLog, int64, error) {
	// very simple filter implementation for testing
	filtered := make([]*entity.AuditLog, 0)
	for _, l := range r.logs {
		if filter.Resource != "" && l.Resource != filter.Resource {
			continue
		}
		if filter.Action != "" && l.Action != filter.Action {
			continue
		}
		if filter.ActorUsername != "" && l.ActorUsername != filter.ActorUsername {
			continue
		}
		filtered = append(filtered, l)
	}

	total := int64(len(filtered))
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	offset := (filter.Page - 1) * filter.PageSize
	end := offset + filter.PageSize
	if offset > len(filtered) {
		return []*entity.AuditLog{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], total, nil
}

func TestAuditLogUseCase_ListAuditLogs(t *testing.T) {
	repo := &inMemoryAuditLogRepo{}
	uc := NewAuditLogUseCase(repo)

	now := time.Now()
	// seed some logs
	repo.logs = []*entity.AuditLog{
		{
			ID:            uuid.New(),
			ActorUsername: "admin",
			ActorRoles:    "[\"admin\"]",
			Action:        "user:create",
			Resource:      "user",
			TargetID:      uuid.New().String(),
			RequestPath:   "/api/users",
			CreatedAt:     now,
		},
		{
			ID:            uuid.New(),
			ActorUsername: "admin",
			ActorRoles:    "[\"admin\"]",
			Action:        "role:create",
			Resource:      "role",
			TargetID:      uuid.New().String(),
			RequestPath:   "/api/roles",
			CreatedAt:     now,
		},
	}

	ctx := context.Background()
	resp, err := uc.ListAuditLogs(ctx, 1, 10, "user", "user:create", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, len(resp.Logs))

	logResp := resp.Logs[0]
	assert.Equal(t, "user", logResp.Resource)
	assert.Equal(t, "user:create", logResp.Action)
	assert.Equal(t, "admin", logResp.ActorUsername)
}
