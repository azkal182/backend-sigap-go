package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

// RoleUseCase handles role management use cases
type RoleUseCase struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	auditLogger    appService.AuditLogger
}

// NewRoleUseCase creates a new role use case
func NewRoleUseCase(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	auditLogger appService.AuditLogger,
) *RoleUseCase {
	return &RoleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		auditLogger:    auditLogger,
	}
}

// CreateRole creates a new role
func (uc *RoleUseCase) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Check if role with same slug already exists
	existingRole, _ := uc.roleRepo.GetBySlug(ctx, req.Slug)
	if existingRole != nil {
		return nil, domainErrors.ErrRoleAlreadyExists
	}

	// Create role
	role := &entity.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        strings.ToLower(req.Slug),
		IsActive:    req.IsActive,
		IsProtected: req.IsProtected,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		permissions := make([]entity.Permission, 0)
		for _, permIDStr := range req.PermissionIDs {
			permID, err := uuid.Parse(permIDStr)
			if err != nil {
				continue
			}
			permission, err := uc.permissionRepo.GetByID(ctx, permID)
			if err != nil {
				continue
			}
			permissions = append(permissions, *permission)
		}
		role.Permissions = permissions
	}

	// Save role
	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Get role with permissions
	roleWithPerms, err := uc.roleRepo.GetWithPermissions(ctx, role.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "role", "role:create", role.ID.String(), map[string]string{
		"name": role.Name,
		"slug": role.Slug,
	})

	return uc.toRoleResponse(roleWithPerms), nil
}

// GetRoleByID retrieves a role by ID
func (uc *RoleUseCase) GetRoleByID(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.GetWithPermissions(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrRoleNotFound
	}

	return uc.toRoleResponse(role), nil
}

// UpdateRole updates a role
func (uc *RoleUseCase) UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	// Get existing role
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrRoleNotFound
	}

	// Update fields
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Slug != "" {
		// Check if slug is already taken by another role
		existingRole, _ := uc.roleRepo.GetBySlug(ctx, req.Slug)
		if existingRole != nil && existingRole.ID != id {
			return nil, domainErrors.ErrRoleAlreadyExists
		}
		role.Slug = strings.ToLower(req.Slug)
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	role.UpdatedAt = time.Now()

	// Save updated role
	if err := uc.roleRepo.Update(ctx, role); err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Get updated role with permissions
	roleWithPerms, err := uc.roleRepo.GetWithPermissions(ctx, role.ID)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "role", "role:update", role.ID.String(), map[string]string{
		"name": role.Name,
		"slug": role.Slug,
	})

	return uc.toRoleResponse(roleWithPerms), nil
}

// DeleteRole deletes a role
func (uc *RoleUseCase) DeleteRole(ctx context.Context, id uuid.UUID) error {
	// Check if role exists
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return domainErrors.ErrRoleNotFound
	}

	// Prevent deletion of protected roles
	if role.IsProtected {
		return domainErrors.ErrProtectedRole
	}

	// Delete role
	if err := uc.roleRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log (best-effort)
	_ = uc.auditLogger.Log(ctx, "role", "role:delete", id.String(), map[string]string{
		"name": role.Name,
		"slug": role.Slug,
	})

	return nil
}

// ListRoles retrieves a paginated list of roles
func (uc *RoleUseCase) ListRoles(ctx context.Context, page, pageSize int) (*dto.ListRolesResponse, error) {
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

	roles, total, err := uc.roleRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	roleResponses := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		roleWithPerms, _ := uc.roleRepo.GetWithPermissions(ctx, role.ID)
		roleResponses = append(roleResponses, *uc.toRoleResponse(roleWithPerms))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListRolesResponse{
		Roles:      roleResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// AssignPermission assigns a permission to a role
func (uc *RoleUseCase) AssignPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	// Check if role exists
	role, err := uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return domainErrors.ErrRoleNotFound
	}

	// Check if role is protected
	if role.IsProtected {
		return domainErrors.ErrProtectedRole
	}

	// Check if permission exists
	_, err = uc.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return domainErrors.ErrPermissionNotFound
	}

	// Assign permission
	return uc.roleRepo.AssignPermission(ctx, roleID, permissionID)
}

// RemovePermission removes a permission from a role
func (uc *RoleUseCase) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	// Check if role exists
	role, err := uc.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return domainErrors.ErrRoleNotFound
	}

	// Check if role is protected
	if role.IsProtected {
		return domainErrors.ErrProtectedRole
	}

	// Remove permission
	return uc.roleRepo.RemovePermission(ctx, roleID, permissionID)
}

// toRoleResponse converts entity.Role to dto.RoleResponse
func (uc *RoleUseCase) toRoleResponse(role *entity.Role) *dto.RoleResponse {
	permissions := make([]string, 0, len(role.Permissions))
	for _, perm := range role.Permissions {
		permissions = append(permissions, perm.Name)
	}

	return &dto.RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		Slug:        role.Slug,
		IsActive:    role.IsActive,
		IsProtected: role.IsProtected,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}
}
