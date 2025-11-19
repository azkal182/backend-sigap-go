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

// ClassHandler handles class related HTTP requests.
type ClassHandler struct {
	classUseCase *usecase.ClassUseCase
}

// NewClassHandler creates a new ClassHandler instance.
func NewClassHandler(classUseCase *usecase.ClassUseCase) *ClassHandler {
	return &ClassHandler{classUseCase: classUseCase}
}

// CreateClass handles POST /api/classes.
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	result, err := h.classUseCase.CreateClass(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to create class", err.Error())
		return
	}

	response.SuccessCreated(c, result, "Class created successfully")
}

// GetClass handles GET /api/classes/:id.
func (h *ClassHandler) GetClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class ID", err.Error())
		return
	}

	result, err := h.classUseCase.GetClass(c.Request.Context(), classID)
	if err != nil {
		switch err {
		case domainErrors.ErrClassNotFound:
			response.ErrorNotFound(c, "Class not found")
		default:
			response.ErrorInternalServer(c, "Failed to get class", err.Error())
		}
		return
	}

	response.SuccessOK(c, result, "Class retrieved successfully")
}

// ListClasses handles GET /api/classes?fan_id=...
func (h *ClassHandler) ListClasses(c *gin.Context) {
	fanID, err := uuid.Parse(c.Query("fan_id"))
	if err != nil {
		response.ErrorBadRequest(c, "fan_id query parameter is required and must be UUID", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.classUseCase.ListClassesByFan(c.Request.Context(), fanID, page, pageSize)
	if err != nil {
		switch err {
		case domainErrors.ErrFanNotFound:
			response.ErrorNotFound(c, "Fan not found")
		default:
			response.ErrorInternalServer(c, "Failed to list classes", err.Error())
		}
		return
	}

	response.SuccessOK(c, result, "Classes retrieved successfully")
}

// UpdateClass handles PUT /api/classes/:id.
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class ID", err.Error())
		return
	}

	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	result, err := h.classUseCase.UpdateClass(c.Request.Context(), classID, req)
	if err != nil {
		switch err {
		case domainErrors.ErrClassNotFound:
			response.ErrorNotFound(c, "Class not found")
		default:
			response.ErrorInternalServer(c, "Failed to update class", err.Error())
		}
		return
	}

	response.SuccessOK(c, result, "Class updated successfully")
}

// DeleteClass handles DELETE /api/classes/:id.
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class ID", err.Error())
		return
	}

	if err := h.classUseCase.DeleteClass(c.Request.Context(), classID); err != nil {
		switch err {
		case domainErrors.ErrClassNotFound:
			response.ErrorNotFound(c, "Class not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete class", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// EnrollStudent handles POST /api/classes/:id/students.
func (h *ClassHandler) EnrollStudent(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class ID", err.Error())
		return
	}

	var req dto.EnrollStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.classUseCase.EnrollStudent(c.Request.Context(), classID, req); err != nil {
		switch err {
		case domainErrors.ErrClassNotFound:
			response.ErrorNotFound(c, "Class not found")
		case domainErrors.ErrStudentNotFound:
			response.ErrorNotFound(c, "Student not found")
		case domainErrors.ErrStudentAlreadyEnrolled:
			response.ErrorConflict(c, "Student already enrolled in class")
		default:
			response.ErrorInternalServer(c, "Failed to enroll student", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// AssignStaff handles POST /api/classes/:id/staff.
func (h *ClassHandler) AssignStaff(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid class ID", err.Error())
		return
	}

	var req dto.AssignClassStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.classUseCase.AssignStaff(c.Request.Context(), classID, req); err != nil {
		switch err {
		case domainErrors.ErrClassNotFound:
			response.ErrorNotFound(c, "Class not found")
		default:
			response.ErrorInternalServer(c, "Failed to assign staff", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}
