package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// AuditLogHandler handles audit log read-only requests
type AuditLogHandler struct {
	useCase *usecase.AuditLogUseCase
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(useCase *usecase.AuditLogUseCase) *AuditLogHandler {
	return &AuditLogHandler{useCase: useCase}
}

// ListAuditLogs lists audit logs with pagination and simple filters
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	resource := c.Query("resource")
	action := c.Query("action")
	actorUsername := c.Query("actor_username")

	resp, err := h.useCase.ListAuditLogs(c.Request.Context(), page, pageSize, resource, action, actorUsername)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list audit logs", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Audit logs retrieved successfully")
}
