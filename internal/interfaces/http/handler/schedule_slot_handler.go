package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// ScheduleSlotHandler handles schedule slot endpoints.
type ScheduleSlotHandler struct {
	slotUseCase *usecase.ScheduleSlotUseCase
}

// NewScheduleSlotHandler constructs handler.
func NewScheduleSlotHandler(slotUseCase *usecase.ScheduleSlotUseCase) *ScheduleSlotHandler {
	return &ScheduleSlotHandler{slotUseCase: slotUseCase}
}

// CreateScheduleSlot handles POST /api/schedule-slots.
func (h *ScheduleSlotHandler) CreateScheduleSlot(c *gin.Context) {
	var req dto.CreateScheduleSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	slot, err := h.slotUseCase.CreateScheduleSlot(c.Request.Context(), req)
	if err != nil {
		switch err {
		case domainErrors.ErrDormitoryNotFound:
			response.ErrorNotFound(c, "Dormitory not found")
		case domainErrors.ErrScheduleSlotConflict:
			response.ErrorConflict(c, "Schedule slot conflict", err.Error())
		case domainErrors.ErrBadRequest:
			response.ErrorBadRequest(c, "Invalid schedule slot data", err.Error())
		default:
			response.ErrorInternalServer(c, "Failed to create schedule slot", err.Error())
		}
		return
	}

	response.SuccessCreated(c, slot, "Schedule slot created successfully")
}

// ListScheduleSlots handles GET /api/schedule-slots.
func (h *ScheduleSlotHandler) ListScheduleSlots(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	dormitoryID := c.Query("dormitory_id")

	var isActive *bool
	if val := c.Query("is_active"); val != "" {
		parsed := val == "true" || val == "1"
		isActive = &parsed
	}

	result, err := h.slotUseCase.ListScheduleSlots(c.Request.Context(), dormitoryID, page, pageSize, isActive)
	if err != nil {
		if err == domainErrors.ErrBadRequest {
			response.ErrorBadRequest(c, "Invalid filters", err.Error())
			return
		}
		response.ErrorInternalServer(c, "Failed to list schedule slots", err.Error())
		return
	}

	response.SuccessOK(c, result, "Schedule slots retrieved successfully")
}

// GetScheduleSlot handles GET /api/schedule-slots/:id.
func (h *ScheduleSlotHandler) GetScheduleSlot(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid schedule slot ID", err.Error())
		return
	}

	slot, err := h.slotUseCase.GetScheduleSlot(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrScheduleSlotNotFound:
			response.ErrorNotFound(c, "Schedule slot not found")
		default:
			response.ErrorInternalServer(c, "Failed to get schedule slot", err.Error())
		}
		return
	}

	response.SuccessOK(c, slot, "Schedule slot retrieved successfully")
}

// UpdateScheduleSlot handles PUT /api/schedule-slots/:id.
func (h *ScheduleSlotHandler) UpdateScheduleSlot(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid schedule slot ID", err.Error())
		return
	}

	var req dto.UpdateScheduleSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	slot, err := h.slotUseCase.UpdateScheduleSlot(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrScheduleSlotNotFound:
			response.ErrorNotFound(c, "Schedule slot not found")
		case domainErrors.ErrScheduleSlotConflict:
			response.ErrorConflict(c, "Schedule slot conflict", err.Error())
		case domainErrors.ErrBadRequest:
			response.ErrorBadRequest(c, "Invalid schedule slot data", err.Error())
		default:
			response.ErrorInternalServer(c, "Failed to update schedule slot", err.Error())
		}
		return
	}

	response.SuccessOK(c, slot, "Schedule slot updated successfully")
}

// DeleteScheduleSlot handles DELETE /api/schedule-slots/:id.
func (h *ScheduleSlotHandler) DeleteScheduleSlot(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid schedule slot ID", err.Error())
		return
	}

	if err := h.slotUseCase.DeleteScheduleSlot(c.Request.Context(), id); err != nil {
		switch err {
		case domainErrors.ErrScheduleSlotNotFound:
			response.ErrorNotFound(c, "Schedule slot not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete schedule slot", err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
