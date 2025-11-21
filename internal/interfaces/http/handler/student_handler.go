package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// StudentHandler manages student HTTP endpoints.
type StudentHandler struct {
	studentUseCase   *usecase.StudentUseCase
	sksResultUseCase *usecase.StudentSKSResultUseCase
}

// CreateStudentSKSResult handles POST /api/students/:id/sks-results.
func (h *StudentHandler) CreateStudentSKSResult(c *gin.Context) {
	studentID := c.Param("id")
	var req dto.CreateStudentSKSResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}
	req.StudentID = studentID
	result, err := h.sksResultUseCase.CreateStudentSKSResult(c.Request.Context(), req)
	if err != nil {
		h.handleSKSError(c, err, "create")
		return
	}
	response.SuccessCreated(c, result, "Student SKS result created successfully")
}

// UpdateStudentSKSResult handles PUT /api/students/:id/sks-results/:result_id.
func (h *StudentHandler) UpdateStudentSKSResult(c *gin.Context) {
	resultID, err := uuid.Parse(c.Param("result_id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS result ID", err.Error())
		return
	}
	var req dto.UpdateStudentSKSResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}
	result, err := h.sksResultUseCase.UpdateStudentSKSResult(c.Request.Context(), resultID, req)
	if err != nil {
		h.handleSKSError(c, err, "update")
		return
	}
	response.SuccessOK(c, result, "Student SKS result updated successfully")
}

// ListStudentSKSResults handles GET /api/students/:id/sks-results.
func (h *StudentHandler) ListStudentSKSResults(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid student ID", err.Error())
		return
	}
	fanID := c.Query("fan_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	results, err := h.sksResultUseCase.ListStudentSKSResults(c.Request.Context(), studentID, fanID, page, pageSize)
	if err != nil {
		h.handleSKSError(c, err, "list")
		return
	}

	response.SuccessOK(c, results, "Student SKS results retrieved successfully")
}

// ListFanCompletionStatuses handles GET /api/students/:id/fans.
func (h *StudentHandler) ListFanCompletionStatuses(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid student ID", err.Error())
		return
	}
	statuses, err := h.sksResultUseCase.ListFanCompletionStatuses(c.Request.Context(), studentID)
	if err != nil {
		h.handleSKSError(c, err, "list")
		return
	}
	response.SuccessOK(c, statuses, "Student FAN completion statuses retrieved successfully")
}

func (h *StudentHandler) handleSKSError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrStudentNotFound:
		response.ErrorNotFound(c, "Student not found")
	case domainErrors.ErrSKSDefinitionNotFound:
		response.ErrorNotFound(c, "SKS definition not found")
	case domainErrors.ErrStudentSKSResultNotFound:
		response.ErrorNotFound(c, "Student SKS result not found")
	case domainErrors.ErrTeacherNotFound:
		response.ErrorNotFound(c, "Teacher not found")
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid SKS result data", err.Error())
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" student SKS result", err.Error())
	}
}

// NewStudentHandler constructs StudentHandler.
func NewStudentHandler(studentUseCase *usecase.StudentUseCase, sksResultUseCase *usecase.StudentSKSResultUseCase) *StudentHandler {
	return &StudentHandler{
		studentUseCase:   studentUseCase,
		sksResultUseCase: sksResultUseCase,
	}
}

// CreateStudent handles POST /api/students
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	resp, err := h.studentUseCase.CreateStudent(c.Request.Context(), req)
	if err != nil {
		switch err {
		case domainErrors.ErrStudentAlreadyExists:
			response.ErrorConflict(c, "Student already exists")
		default:
			response.ErrorInternalServer(c, "Failed to create student", err.Error())
		}
		return
	}

	response.SuccessCreated(c, resp, "Student created successfully")
}

// GetStudent handles GET /api/students/:id
func (h *StudentHandler) GetStudent(c *gin.Context) {
	studentID, err := parseUUIDParam(c, "id")
	if err != nil {
		return
	}

	resp, err := h.studentUseCase.GetStudentByID(c.Request.Context(), studentID)
	if err != nil {
		switch err {
		case domainErrors.ErrStudentNotFound:
			response.ErrorNotFound(c, "Student not found")
		default:
			response.ErrorInternalServer(c, "Failed to get student", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Student retrieved successfully")
}

// ListStudents handles GET /api/students
func (h *StudentHandler) ListStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.studentUseCase.ListStudents(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list students", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Students retrieved successfully")
}

// UpdateStudent handles PUT /api/students/:id
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	studentID, err := parseUUIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	resp, err := h.studentUseCase.UpdateStudent(c.Request.Context(), studentID, req)
	if err != nil {
		switch err {
		case domainErrors.ErrStudentNotFound:
			response.ErrorNotFound(c, "Student not found")
		default:
			response.ErrorInternalServer(c, "Failed to update student", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Student updated successfully")
}

// UpdateStudentStatus handles PATCH /api/students/:id/status
func (h *StudentHandler) UpdateStudentStatus(c *gin.Context) {
	studentID, err := parseUUIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	resp, err := h.studentUseCase.UpdateStudentStatus(c.Request.Context(), studentID, req.Status)
	if err != nil {
		switch err {
		case domainErrors.ErrStudentNotFound:
			response.ErrorNotFound(c, "Student not found")
		default:
			response.ErrorInternalServer(c, "Failed to update student status", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Student status updated successfully")
}

// MutateStudentDormitory handles POST /api/students/:id/mutate-dormitory
func (h *StudentHandler) MutateStudentDormitory(c *gin.Context) {
	studentID, err := parseUUIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.MutateStudentDormitoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	dormitoryID, err := uuid.Parse(req.DormitoryID)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
		return
	}

	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	resp, err := h.studentUseCase.MutateStudentDormitory(c.Request.Context(), studentID, dormitoryID, startDate)
	if err != nil {
		switch err {
		case domainErrors.ErrStudentNotFound:
			response.ErrorNotFound(c, "Student not found")
		case domainErrors.ErrDormitoryNotFound:
			response.ErrorNotFound(c, "Dormitory not found")
		default:
			response.ErrorInternalServer(c, "Failed to mutate dormitory", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Student dormitory mutated successfully")
}

func parseUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid UUID", err.Error())
		return uuid.Nil, err
	}
	return id, nil
}
