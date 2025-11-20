package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// TeacherHandler handles teacher HTTP endpoints.
type TeacherHandler struct {
	teacherUseCase *usecase.TeacherUseCase
}

// NewTeacherHandler constructs TeacherHandler.
func NewTeacherHandler(teacherUseCase *usecase.TeacherUseCase) *TeacherHandler {
	return &TeacherHandler{teacherUseCase: teacherUseCase}
}

// CreateTeacher handles POST /api/teachers.
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req dto.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	teacher, err := h.teacherUseCase.CreateTeacher(c.Request.Context(), req)
	if err != nil {
		switch err {
		case domainErrors.ErrTeacherAlreadyExists:
			response.ErrorConflict(c, "Teacher already exists", err.Error())
		case domainErrors.ErrUserNotFound:
			response.ErrorNotFound(c, "Existing user not found", err.Error())
		case domainErrors.ErrTeacherUserAssigned:
			response.ErrorConflict(c, "User already linked to another teacher", err.Error())
		case domainErrors.ErrRoleNotFound:
			response.ErrorBadRequest(c, "Teacher role missing", err.Error())
		default:
			response.ErrorInternalServer(c, "Failed to create teacher", err.Error())
		}
		return
	}

	response.SuccessCreated(c, teacher, "Teacher created successfully")
}

// ListTeachers handles GET /api/teachers.
func (h *TeacherHandler) ListTeachers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	var isActive *bool
	if raw := strings.TrimSpace(c.Query("is_active")); raw != "" {
		val := strings.EqualFold(raw, "true")
		if strings.EqualFold(raw, "false") {
			val = false
		}
		isActive = &val
	}

	result, err := h.teacherUseCase.ListTeachers(c.Request.Context(), page, pageSize, keyword, isActive)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list teachers", err.Error())
		return
	}

	response.SuccessOK(c, result, "Teachers retrieved successfully")
}

// GetTeacher handles GET /api/teachers/:id.
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid teacher ID", err.Error())
		return
	}

	teacher, err := h.teacherUseCase.GetTeacher(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrTeacherNotFound:
			response.ErrorNotFound(c, "Teacher not found")
		default:
			response.ErrorInternalServer(c, "Failed to get teacher", err.Error())
		}
		return
	}

	response.SuccessOK(c, teacher, "Teacher retrieved successfully")
}

// UpdateTeacher handles PUT /api/teachers/:id.
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid teacher ID", err.Error())
		return
	}

	var req dto.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	teacher, err := h.teacherUseCase.UpdateTeacher(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrTeacherNotFound:
			response.ErrorNotFound(c, "Teacher not found")
		default:
			response.ErrorInternalServer(c, "Failed to update teacher", err.Error())
		}
		return
	}

	response.SuccessOK(c, teacher, "Teacher updated successfully")
}

// DeactivateTeacher handles DELETE /api/teachers/:id (soft delete).
func (h *TeacherHandler) DeactivateTeacher(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid teacher ID", err.Error())
		return
	}

	if err := h.teacherUseCase.DeactivateTeacher(c.Request.Context(), id); err != nil {
		switch err {
		case domainErrors.ErrTeacherNotFound:
			response.ErrorNotFound(c, "Teacher not found")
		default:
			response.ErrorInternalServer(c, "Failed to deactivate teacher", err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
