package usecase

import (
	"context"

	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// PermissionUseCase handles read-only permission use cases
type PermissionUseCase struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionUseCase creates a new permission use case
func NewPermissionUseCase(permissionRepo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{permissionRepo: permissionRepo}
}

// ListPermissions retrieves a paginated list of permissions
func (uc *PermissionUseCase) ListPermissions(ctx context.Context, page, pageSize int) (*dto.ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	permissions, total, err := uc.permissionRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.PermissionResponse, 0, len(permissions))
	for _, p := range permissions {
		items = append(items, dto.PermissionResponse{
			ID:       p.ID.String(),
			Name:     p.Name,
			Slug:     p.Slug,
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListPermissionsResponse{
		Permissions: items,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}
