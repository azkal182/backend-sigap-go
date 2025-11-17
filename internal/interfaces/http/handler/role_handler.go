package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// RoleHandler handles role management requests
type RoleHandler struct {
	roleUseCase *usecase.RoleUseCase
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleUseCase *usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{
		roleUseCase: roleUseCase,
	}
}

// CreateRole handles role creation
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.roleUseCase.CreateRole(c.Request.Context(), req)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleAlreadyExists:
			response.ErrorConflict(c, "Role already exists")
		default:
			response.ErrorInternalServer(c, "Failed to create role", err.Error())
		}
		return
	}

	response.SuccessCreated(c, resp, "Role created successfully")
}

// GetRole handles getting a role by ID
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	resp, err := h.roleUseCase.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		default:
			response.ErrorInternalServer(c, "Failed to get role", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Role retrieved successfully")
}

// UpdateRole handles role update
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.roleUseCase.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		case domainErrors.ErrRoleAlreadyExists:
			response.ErrorConflict(c, "Slug already taken")
		default:
			response.ErrorInternalServer(c, "Failed to update role", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Role updated successfully")
}

// DeleteRole handles role deletion
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	err = h.roleUseCase.DeleteRole(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		case domainErrors.ErrProtectedRole:
			response.ErrorForbidden(c, "Cannot delete protected role")
		default:
			response.ErrorInternalServer(c, "Failed to delete role", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// ListRoles handles listing roles with pagination
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.roleUseCase.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list roles", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Roles retrieved successfully")
}

// AssignPermission handles assigning a permission to a role
func (h *RoleHandler) AssignPermission(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid permission ID", err.Error())
		return
	}

	err = h.roleUseCase.AssignPermission(c.Request.Context(), roleID, permissionID)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		case domainErrors.ErrPermissionNotFound:
			response.ErrorNotFound(c, "Permission not found")
		case domainErrors.ErrProtectedRole:
			response.ErrorForbidden(c, "Cannot modify permissions of protected role")
		default:
			response.ErrorInternalServer(c, "Failed to assign permission", err.Error())
		}
		return
	}

	response.SuccessOK(c, nil, "Permission assigned successfully")
}

// RemovePermission handles removing a permission from a role
func (h *RoleHandler) RemovePermission(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid role ID", err.Error())
		return
	}

	var req dto.RemovePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid permission ID", err.Error())
		return
	}

	err = h.roleUseCase.RemovePermission(c.Request.Context(), roleID, permissionID)
	if err != nil {
		switch err {
		case domainErrors.ErrRoleNotFound:
			response.ErrorNotFound(c, "Role not found")
		case domainErrors.ErrPermissionNotFound:
			response.ErrorNotFound(c, "Permission not found")
		case domainErrors.ErrProtectedRole:
			response.ErrorForbidden(c, "Cannot modify permissions of protected role")
		default:
			response.ErrorInternalServer(c, "Failed to remove permission", err.Error())
		}
		return
	}

	response.SuccessOK(c, nil, "Permission removed successfully")
}
