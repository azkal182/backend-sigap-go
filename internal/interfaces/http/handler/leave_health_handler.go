package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// LeavePermitHandler manages leave permit endpoints.
type LeavePermitHandler struct {
	useCase *usecase.LeavePermitUseCase
}

// HealthStatusHandler manages health status endpoints.
type HealthStatusHandler struct {
	useCase *usecase.HealthStatusUseCase
}

// NewLeavePermitHandler constructs LeavePermitHandler.
func NewLeavePermitHandler(useCase *usecase.LeavePermitUseCase) *LeavePermitHandler {
	return &LeavePermitHandler{useCase: useCase}
}

// NewHealthStatusHandler constructs HealthStatusHandler.
func NewHealthStatusHandler(useCase *usecase.HealthStatusUseCase) *HealthStatusHandler {
	return &HealthStatusHandler{useCase: useCase}
}

// CreateLeavePermit handles POST /api/leave-permits.
func (h *LeavePermitHandler) CreateLeavePermit(c *gin.Context) {
	var req dto.CreateLeavePermitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	permit, err := h.useCase.CreateLeavePermit(c.Request.Context(), req)
	if err != nil {
		h.handleLeavePermitError(c, err, "create")
		return
	}

	response.SuccessCreated(c, permit, "Leave permit created successfully")
}

// ListLeavePermits handles GET /api/leave-permits.
func (h *LeavePermitHandler) ListLeavePermits(c *gin.Context) {
	var req dto.ListLeavePermitsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	permits, err := h.useCase.ListLeavePermits(c.Request.Context(), req)
	if err != nil {
		h.handleLeavePermitError(c, err, "list")
		return
	}

	response.SuccessOK(c, permits, "Leave permits retrieved successfully")
}

// ApproveLeavePermit handles PUT /api/leave-permits/:id/approve.
func (h *LeavePermitHandler) ApproveLeavePermit(c *gin.Context) {
	h.updateLeavePermitStatus(c, string(entity.LeavePermitStatusApproved))
}

// RejectLeavePermit handles PUT /api/leave-permits/:id/reject.
func (h *LeavePermitHandler) RejectLeavePermit(c *gin.Context) {
	h.updateLeavePermitStatus(c, string(entity.LeavePermitStatusRejected))
}

// CompleteLeavePermit handles PUT /api/leave-permits/:id/complete.
func (h *LeavePermitHandler) CompleteLeavePermit(c *gin.Context) {
	h.updateLeavePermitStatus(c, string(entity.LeavePermitStatusCompleted))
}

func (h *LeavePermitHandler) updateLeavePermitStatus(c *gin.Context, status string) {
	permitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid leave permit ID", err.Error())
		return
	}

	resp, err := h.useCase.UpdateLeavePermitStatus(c.Request.Context(), permitID, dto.UpdateLeavePermitStatusRequest{Status: status})
	if err != nil {
		h.handleLeavePermitError(c, err, "update leave permit status")
		return
	}

	response.SuccessOK(c, resp, "Leave permit status updated successfully")
}

func (h *LeavePermitHandler) handleLeavePermitError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrStudentNotFound:
		response.ErrorNotFound(c, "Student not found")
	case domainErrors.ErrLeavePermitNotFound:
		response.ErrorNotFound(c, "Leave permit not found")
	case domainErrors.ErrLeavePermitConflict:
		response.ErrorConflict(c, "Leave permit overlaps with existing permit")
	case domainErrors.ErrLeavePermitStatus:
		response.ErrorBadRequest(c, "Invalid leave permit status transition", err.Error())
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid leave permit data", err.Error())
	case domainErrors.ErrUnauthorized:
		response.ErrorUnauthorized(c, "Unauthorized")
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" leave permit", err.Error())
	}
}

// CreateHealthStatus handles POST /api/health-statuses.
func (h *HealthStatusHandler) CreateHealthStatus(c *gin.Context) {
	var req dto.CreateHealthStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	status, err := h.useCase.CreateHealthStatus(c.Request.Context(), req)
	if err != nil {
		h.handleHealthStatusError(c, err, "create")
		return
	}

	response.SuccessCreated(c, status, "Health status created successfully")
}

// ListHealthStatuses handles GET /api/health-statuses.
func (h *HealthStatusHandler) ListHealthStatuses(c *gin.Context) {
	var req dto.ListHealthStatusesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	statuses, err := h.useCase.ListHealthStatuses(c.Request.Context(), req)
	if err != nil {
		h.handleHealthStatusError(c, err, "list")
		return
	}

	response.SuccessOK(c, statuses, "Health statuses retrieved successfully")
}

// RevokeHealthStatus handles PUT /api/health-statuses/:id/revoke.
func (h *HealthStatusHandler) RevokeHealthStatus(c *gin.Context) {
	statusID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid health status ID", err.Error())
		return
	}

	var req dto.RevokeHealthStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	status, err := h.useCase.RevokeHealthStatus(c.Request.Context(), statusID, req)
	if err != nil {
		h.handleHealthStatusError(c, err, "revoke")
		return
	}

	response.SuccessOK(c, status, "Health status revoked successfully")
}

func (h *HealthStatusHandler) handleHealthStatusError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrStudentNotFound:
		response.ErrorNotFound(c, "Student not found")
	case domainErrors.ErrHealthStatusNotFound:
		response.ErrorNotFound(c, "Health status not found")
	case domainErrors.ErrHealthStatusActive:
		response.ErrorConflict(c, "Health status already active for the given period")
	case domainErrors.ErrHealthStatusForbidden:
		response.ErrorBadRequest(c, "Invalid health status operation", err.Error())
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid health status data", err.Error())
	case domainErrors.ErrUnauthorized:
		response.ErrorUnauthorized(c, "Unauthorized")
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" health status", err.Error())
	}
}
