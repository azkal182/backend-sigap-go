package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// ReportHandler handles read-only reporting endpoints.
type ReportHandler struct {
	reportUseCase *usecase.ReportUseCase
}

// NewReportHandler creates a new ReportHandler instance.
func NewReportHandler(reportUseCase *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUseCase: reportUseCase}
}

// GetStudentAttendanceReport handles GET /api/reports/attendance/students.
func (h *ReportHandler) GetStudentAttendanceReport(c *gin.Context) {
	var req dto.StudentAttendanceReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetStudentAttendanceReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch student attendance report", err.Error())
		return
	}

	response.SuccessOK(c, report, "Student attendance report retrieved successfully")
}

// GetTeacherAttendanceReport handles GET /api/reports/attendance/teachers.
func (h *ReportHandler) GetTeacherAttendanceReport(c *gin.Context) {
	var req dto.TeacherAttendanceReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetTeacherAttendanceReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch teacher attendance report", err.Error())
		return
	}

	response.SuccessOK(c, report, "Teacher attendance report retrieved successfully")
}

// GetLeavePermitReport handles GET /api/reports/leave-permits.
func (h *ReportHandler) GetLeavePermitReport(c *gin.Context) {
	var req dto.LeavePermitReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetLeavePermitReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch leave permit report", err.Error())
		return
	}

	response.SuccessOK(c, report, "Leave permit report retrieved successfully")
}

// GetHealthStatusReport handles GET /api/reports/health-statuses.
func (h *ReportHandler) GetHealthStatusReport(c *gin.Context) {
	var req dto.HealthStatusReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetHealthStatusReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch health status report", err.Error())
		return
	}

	response.SuccessOK(c, report, "Health status report retrieved successfully")
}

// GetSKSReport handles GET /api/reports/sks.
func (h *ReportHandler) GetSKSReport(c *gin.Context) {
	var req dto.SKSReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetSKSReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch SKS report", err.Error())
		return
	}

	response.SuccessOK(c, report, "SKS report retrieved successfully")
}

// GetMutationReport handles GET /api/reports/mutations.
func (h *ReportHandler) GetMutationReport(c *gin.Context) {
	var req dto.MutationReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	report, err := h.reportUseCase.GetMutationReport(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to fetch mutation report", err.Error())
		return
	}

	response.SuccessOK(c, report, "Mutation report retrieved successfully")
}
