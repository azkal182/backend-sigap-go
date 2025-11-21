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

// AttendanceHandler manages attendance endpoints.
type AttendanceHandler struct {
	attendanceUseCase *usecase.AttendanceUseCase
}

// NewAttendanceHandler constructs AttendanceHandler.
func NewAttendanceHandler(attendanceUseCase *usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{attendanceUseCase: attendanceUseCase}
}

// OpenSessions handles POST /api/attendance-sessions/open.
func (h *AttendanceHandler) OpenSessions(c *gin.Context) {
	var req dto.OpenAttendanceSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.attendanceUseCase.OpenSessions(c.Request.Context(), req); err != nil {
		h.handleAttendanceError(c, err, "open attendance sessions")
		return
	}

	response.SuccessOK(c, gin.H{"status": "opened"}, "Attendance sessions opened")
}

// SubmitStudentAttendance handles POST /api/attendance-sessions/:id/students.
func (h *AttendanceHandler) SubmitStudentAttendance(c *gin.Context) {
	sessionID, ok := parseAttendanceSessionID(c)
	if !ok {
		return
	}

	var req dto.SubmitStudentAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.attendanceUseCase.SubmitStudentAttendance(c.Request.Context(), sessionID, req); err != nil {
		h.handleAttendanceError(c, err, "submit student attendance")
		return
	}

	response.SuccessOK(c, gin.H{"status": "submitted"}, "Student attendance submitted")
}

// SubmitTeacherAttendance handles POST /api/attendance-sessions/:id/teacher.
func (h *AttendanceHandler) SubmitTeacherAttendance(c *gin.Context) {
	sessionID, ok := parseAttendanceSessionID(c)
	if !ok {
		return
	}

	var req dto.SubmitTeacherAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.attendanceUseCase.SubmitTeacherAttendance(c.Request.Context(), sessionID, req); err != nil {
		h.handleAttendanceError(c, err, "submit teacher attendance")
		return
	}

	response.SuccessOK(c, gin.H{"status": "submitted"}, "Teacher attendance submitted")
}

// LockSessions handles POST /api/attendance-sessions/lock-day.
func (h *AttendanceHandler) LockSessions(c *gin.Context) {
	var req dto.LockAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	if err := h.attendanceUseCase.LockSessions(c.Request.Context(), req); err != nil {
		h.handleAttendanceError(c, err, "lock attendance sessions")
		return
	}

	response.SuccessOK(c, gin.H{"status": "locked"}, "Attendance sessions locked")
}

// ListAttendanceSessions handles GET /api/attendance-sessions.
func (h *AttendanceHandler) ListAttendanceSessions(c *gin.Context) {
	listReq := dto.ListAttendanceSessionsRequest{
		Page:     parseQueryInt(c, "page", 1),
		PageSize: parseQueryInt(c, "page_size", 10),
	}

	if val := c.Query("class_schedule_id"); val != "" {
		listReq.ClassScheduleID = &val
	}
	if val := c.Query("teacher_id"); val != "" {
		listReq.TeacherID = &val
	}
	if val := c.Query("date"); val != "" {
		listReq.Date = &val
	}
	if val := c.Query("status"); val != "" {
		listReq.Status = &val
	}

	resp, err := h.attendanceUseCase.ListAttendanceSessions(c.Request.Context(), listReq)
	if err != nil {
		h.handleAttendanceError(c, err, "list attendance sessions")
		return
	}

	response.SuccessOK(c, resp, "Attendance sessions retrieved")
}

func (h *AttendanceHandler) handleAttendanceError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrBadRequest, domainErrors.ErrAttendanceInvalidStatus:
		response.ErrorBadRequest(c, "Invalid attendance request", err.Error())
	case domainErrors.ErrClassScheduleNotFound:
		response.ErrorNotFound(c, "Class schedule not found", err.Error())
	case domainErrors.ErrAttendanceSessionNotFound:
		response.ErrorNotFound(c, "Attendance session not found", err.Error())
	case domainErrors.ErrAttendanceAlreadyLocked:
		response.ErrorConflict(c, "Attendance session already locked", err.Error())
	default:
		response.ErrorInternalServer(c, "Failed to "+action, err.Error())
	}
}

func parseAttendanceSessionID(c *gin.Context) (uuid.UUID, bool) {
	idStr := c.Param("id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid attendance session ID", err.Error())
		return uuid.Nil, false
	}
	return sessionID, true
}

func parseQueryInt(c *gin.Context, key string, defaultVal int) int {
	if val := c.Query(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
