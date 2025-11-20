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

// SKSExamScheduleHandler exposes SKS exam schedule endpoints.
type SKSExamScheduleHandler struct {
	examUseCase *usecase.SKSExamScheduleUseCase
}

// NewSKSExamScheduleHandler builds a new handler instance.
func NewSKSExamScheduleHandler(examUseCase *usecase.SKSExamScheduleUseCase) *SKSExamScheduleHandler {
	return &SKSExamScheduleHandler{examUseCase: examUseCase}
}

// CreateSKSExamSchedule handles POST /api/sks-exams.
func (h *SKSExamScheduleHandler) CreateSKSExamSchedule(c *gin.Context) {
	var req dto.CreateSKSExamScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	exam, err := h.examUseCase.CreateSKSExamSchedule(c.Request.Context(), req)
	if err != nil {
		h.handleExamError(c, err, "create")
		return
	}

	response.SuccessCreated(c, exam, "SKS exam schedule created successfully")
}

// GetSKSExamSchedule handles GET /api/sks-exams/:id.
func (h *SKSExamScheduleHandler) GetSKSExamSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS exam schedule ID", err.Error())
		return
	}

	exam, err := h.examUseCase.GetSKSExamSchedule(c.Request.Context(), id)
	if err != nil {
		if err == domainErrors.ErrSKSExamScheduleNotFound {
			response.ErrorNotFound(c, "SKS exam schedule not found")
		} else {
			response.ErrorInternalServer(c, "Failed to get SKS exam schedule", err.Error())
		}
		return
	}

	response.SuccessOK(c, exam, "SKS exam schedule retrieved successfully")
}

// ListSKSExamSchedules handles GET /api/sks-exams?sks_id=....
func (h *SKSExamScheduleHandler) ListSKSExamSchedules(c *gin.Context) {
	sksID := c.Query("sks_id")
	if sksID == "" {
		response.ErrorBadRequest(c, "sks_id is required", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	exams, err := h.examUseCase.ListSKSExamSchedules(c.Request.Context(), sksID, page, pageSize)
	if err != nil {
		if err == domainErrors.ErrBadRequest {
			response.ErrorBadRequest(c, "Invalid SKS exam filters", err.Error())
		} else {
			response.ErrorInternalServer(c, "Failed to list SKS exam schedules", err.Error())
		}
		return
	}

	response.SuccessOK(c, exams, "SKS exam schedules retrieved successfully")
}

// UpdateSKSExamSchedule handles PUT /api/sks-exams/:id.
func (h *SKSExamScheduleHandler) UpdateSKSExamSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS exam schedule ID", err.Error())
		return
	}

	var req dto.UpdateSKSExamScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	exam, err := h.examUseCase.UpdateSKSExamSchedule(c.Request.Context(), id, req)
	if err != nil {
		h.handleExamError(c, err, "update")
		return
	}

	response.SuccessOK(c, exam, "SKS exam schedule updated successfully")
}

// DeleteSKSExamSchedule handles DELETE /api/sks-exams/:id.
func (h *SKSExamScheduleHandler) DeleteSKSExamSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS exam schedule ID", err.Error())
		return
	}

	if err := h.examUseCase.DeleteSKSExamSchedule(c.Request.Context(), id); err != nil {
		if err == domainErrors.ErrSKSExamScheduleNotFound {
			response.ErrorNotFound(c, "SKS exam schedule not found")
		} else {
			response.ErrorInternalServer(c, "Failed to delete SKS exam schedule", err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SKSExamScheduleHandler) handleExamError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrSKSDefinitionNotFound:
		response.ErrorNotFound(c, "SKS definition not found")
	case domainErrors.ErrTeacherNotFound:
		response.ErrorNotFound(c, "Teacher not found")
	case domainErrors.ErrSKSExamScheduleNotFound:
		response.ErrorNotFound(c, "SKS exam schedule not found")
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid SKS exam schedule data", err.Error())
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" SKS exam schedule", err.Error())
	}
}
