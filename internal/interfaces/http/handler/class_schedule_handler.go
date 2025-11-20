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

// ClassScheduleHandler exposes HTTP endpoints for class schedules.
type ClassScheduleHandler struct {
	classScheduleUseCase *usecase.ClassScheduleUseCase
}

// NewClassScheduleHandler constructs a new handler.
func NewClassScheduleHandler(classScheduleUseCase *usecase.ClassScheduleUseCase) *ClassScheduleHandler {
	return &ClassScheduleHandler{classScheduleUseCase: classScheduleUseCase}
}

// CreateClassSchedule handles POST /api/class-schedules.
func (h *ClassScheduleHandler) CreateClassSchedule(c *gin.Context) {
	var req dto.CreateClassScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	schedule, err := h.classScheduleUseCase.CreateClassSchedule(c.Request.Context(), req)
	if err != nil {
		h.handleScheduleError(c, err, "create")
		return
	}

	response.SuccessCreated(c, schedule, "Class schedule created successfully")
}

// GetClassSchedule handles GET /api/class-schedules/:id.
func (h *ClassScheduleHandler) GetClassSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class schedule ID", err.Error())
		return
	}

	schedule, err := h.classScheduleUseCase.GetClassSchedule(c.Request.Context(), id)
	if err != nil {
		if err == domainErrors.ErrClassScheduleNotFound {
			response.ErrorNotFound(c, "Class schedule not found")
		} else {
			response.ErrorInternalServer(c, "Failed to get class schedule", err.Error())
		}
		return
	}

	response.SuccessOK(c, schedule, "Class schedule retrieved successfully")
}

// ListClassSchedules handles GET /api/class-schedules.
func (h *ClassScheduleHandler) ListClassSchedules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	classID := c.Query("class_id")
	teacherID := c.Query("teacher_id")
	dormitoryID := c.Query("dormitory_id")
	dayOfWeek := c.Query("day_of_week")

	var isActive *bool
	if val := c.Query("is_active"); val != "" {
		parsed := val == "true" || val == "1"
		isActive = &parsed
	}

	schedules, err := h.classScheduleUseCase.ListClassSchedules(
		c.Request.Context(),
		classID,
		teacherID,
		dormitoryID,
		dayOfWeek,
		page,
		pageSize,
		isActive,
	)
	if err != nil {
		if err == domainErrors.ErrBadRequest {
			response.ErrorBadRequest(c, "Invalid filters", err.Error())
		} else {
			response.ErrorInternalServer(c, "Failed to list class schedules", err.Error())
		}
		return
	}

	response.SuccessOK(c, schedules, "Class schedules retrieved successfully")
}

// UpdateClassSchedule handles PUT /api/class-schedules/:id.
func (h *ClassScheduleHandler) UpdateClassSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class schedule ID", err.Error())
		return
	}

	var req dto.UpdateClassScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	schedule, err := h.classScheduleUseCase.UpdateClassSchedule(c.Request.Context(), id, req)
	if err != nil {
		h.handleScheduleError(c, err, "update")
		return
	}

	response.SuccessOK(c, schedule, "Class schedule updated successfully")
}

// DeleteClassSchedule handles DELETE /api/class-schedules/:id.
func (h *ClassScheduleHandler) DeleteClassSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class schedule ID", err.Error())
		return
	}

	if err := h.classScheduleUseCase.DeleteClassSchedule(c.Request.Context(), id); err != nil {
		if err == domainErrors.ErrClassScheduleNotFound {
			response.ErrorNotFound(c, "Class schedule not found")
		} else {
			response.ErrorInternalServer(c, "Failed to delete class schedule", err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ClassScheduleHandler) handleScheduleError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrClassNotFound:
		response.ErrorNotFound(c, "Class not found")
	case domainErrors.ErrTeacherNotFound:
		response.ErrorNotFound(c, "Teacher not found")
	case domainErrors.ErrDormitoryNotFound:
		response.ErrorNotFound(c, "Dormitory not found")
	case domainErrors.ErrSubjectNotFound:
		response.ErrorNotFound(c, "Subject not found")
	case domainErrors.ErrScheduleSlotNotFound:
		response.ErrorNotFound(c, "Schedule slot not found")
	case domainErrors.ErrScheduleSlotInactive:
		response.ErrorBadRequest(c, "Schedule slot inactive", err.Error())
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid class schedule data", err.Error())
	case domainErrors.ErrClassScheduleNotFound:
		response.ErrorNotFound(c, "Class schedule not found")
	case domainErrors.ErrClassScheduleConflict:
		response.ErrorConflict(c, "Class schedule conflict", err.Error())
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" class schedule", err.Error())
	}
}
