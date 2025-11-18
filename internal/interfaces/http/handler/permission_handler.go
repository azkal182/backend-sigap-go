package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// PermissionHandler handles read-only permission requests
type PermissionHandler struct {
	permissionUseCase *usecase.PermissionUseCase
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler(permissionUseCase *usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{
		permissionUseCase: permissionUseCase,
	}
}

// ListPermissions handles listing permissions with pagination
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.permissionUseCase.ListPermissions(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list permissions", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Permissions retrieved successfully")
}
