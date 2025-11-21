package usecase

import (
	"context"

	"github.com/google/uuid"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
)

// requireActorID ensures an authenticated actor ID is present in context.
func requireActorID(ctx context.Context) (uuid.UUID, error) {
	if id, ok := actorIDFromContext(ctx); ok {
		return *id, nil
	}
	return uuid.Nil, domainErrors.ErrUnauthorized
}

// actorIDFromContext extracts the actor UUID from context if any.
func actorIDFromContext(ctx context.Context) (*uuid.UUID, bool) {
	if v := ctx.Value(appService.CtxKeyActorID); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return &id, true
		}
	}
	return nil, false
}

// actorUsernameFromContext fetches actor username for logging/metadata.
func actorUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value(appService.CtxKeyActorUsername).(string); ok {
		return username
	}
	return ""
}
