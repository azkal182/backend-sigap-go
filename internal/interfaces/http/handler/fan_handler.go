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

// FanHandler handles FAN related HTTP requests.
type FanHandler struct {
	fanUseCase *usecase.FanUseCase
}

// NewFanHandler creates a new FanHandler instance.
func NewFanHandler(fanUseCase *usecase.FanUseCase) *FanHandler {
	return &FanHandler{fanUseCase: fanUseCase}
}

// CreateFan handles POST /api/fans.
func (h *FanHandler) CreateFan(c *gin.Context) {
	var req dto.CreateFanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	result, err := h.fanUseCase.CreateFan(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to create fan", err.Error())
		return
	}

	response.SuccessCreated(c, result, "Fan created successfully")
}

// GetFan handles GET /api/fans/:id.
func (h *FanHandler) GetFan(c *gin.Context) {
	fanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid fan ID", err.Error())
		return
	}

	result, err := h.fanUseCase.GetFan(c.Request.Context(), fanID)
	if err != nil {
		switch err {
		case domainErrors.ErrFanNotFound:
			response.ErrorNotFound(c, "Fan not found")
		default:
			response.ErrorInternalServer(c, "Failed to get fan", err.Error())
		}
		return
	}

	response.SuccessOK(c, result, "Fan retrieved successfully")
}

// ListFans handles GET /api/fans.
func (h *FanHandler) ListFans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}

	if dormIDStr := c.Query("dormitory_id"); dormIDStr != "" {
		dormitoryID, err := uuid.Parse(dormIDStr)
		if err != nil {
			response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
			return
		}
		result, err := h.fanUseCase.ListFansByDormitory(c.Request.Context(), dormitoryID, page, pageSize)
		if err != nil {
			response.ErrorInternalServer(c, "Failed to list fans", err.Error())
			return
		}
		response.SuccessOK(c, result, "Fans retrieved successfully")
		return
	}

	result, err := h.fanUseCase.ListFans(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list fans", err.Error())
		return
	}

	response.SuccessOK(c, result, "Fans retrieved successfully")
}

// UpdateFan handles PUT /api/fans/:id.
func (h *FanHandler) UpdateFan(c *gin.Context) {
	fanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid fan ID", err.Error())
		return
	}

	var req dto.UpdateFanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	result, err := h.fanUseCase.UpdateFan(c.Request.Context(), fanID, req)
	if err != nil {
		switch err {
		case domainErrors.ErrFanNotFound:
			response.ErrorNotFound(c, "Fan not found")
		default:
			response.ErrorInternalServer(c, "Failed to update fan", err.Error())
		}
		return
	}

	response.SuccessOK(c, result, "Fan updated successfully")
}

// DeleteFan handles DELETE /api/fans/:id.
func (h *FanHandler) DeleteFan(c *gin.Context) {
	fanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid fan ID", err.Error())
		return
	}

	if err := h.fanUseCase.DeleteFan(c.Request.Context(), fanID); err != nil {
		switch err {
		case domainErrors.ErrFanNotFound:
			response.ErrorNotFound(c, "Fan not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete fan", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}
