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

// DormitoryHandler handles dormitory management requests
type DormitoryHandler struct {
	dormitoryUseCase *usecase.DormitoryUseCase
}

// NewDormitoryHandler creates a new dormitory handler
func NewDormitoryHandler(dormitoryUseCase *usecase.DormitoryUseCase) *DormitoryHandler {
	return &DormitoryHandler{
		dormitoryUseCase: dormitoryUseCase,
	}
}

// CreateDormitory handles dormitory creation
// @Summary Create a new dormitory
// @Description Create a new dormitory
// @Tags dormitories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDormitoryRequest true "Create dormitory request"
// @Success 201 {object} dto.DormitoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/dormitories [post]
func (h *DormitoryHandler) CreateDormitory(c *gin.Context) {
	var req dto.CreateDormitoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.dormitoryUseCase.CreateDormitory(c.Request.Context(), req)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to create dormitory", err.Error())
		return
	}

	response.SuccessCreated(c, resp, "Dormitory created successfully")
}

// GetDormitory handles getting a dormitory by ID
// @Summary Get dormitory by ID
// @Description Get dormitory details by ID
// @Tags dormitories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dormitory ID"
// @Success 200 {object} dto.DormitoryResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dormitories/{id} [get]
func (h *DormitoryHandler) GetDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
		return
	}

	resp, err := h.dormitoryUseCase.GetDormitoryByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrDormitoryNotFound:
			response.ErrorNotFound(c, "Dormitory not found")
		default:
			response.ErrorInternalServer(c, "Failed to get dormitory", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Dormitory retrieved successfully")
}

// UpdateDormitory handles dormitory update
// @Summary Update dormitory
// @Description Update dormitory details
// @Tags dormitories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dormitory ID"
// @Param request body dto.UpdateDormitoryRequest true "Update dormitory request"
// @Success 200 {object} dto.DormitoryResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dormitories/{id} [put]
func (h *DormitoryHandler) UpdateDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
		return
	}

	var req dto.UpdateDormitoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c, "Invalid request body", err.Error())
		return
	}

	resp, err := h.dormitoryUseCase.UpdateDormitory(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainErrors.ErrDormitoryNotFound:
			response.ErrorNotFound(c, "Dormitory not found")
		default:
			response.ErrorInternalServer(c, "Failed to update dormitory", err.Error())
		}
		return
	}

	response.SuccessOK(c, resp, "Dormitory updated successfully")
}

// DeleteDormitory handles dormitory deletion
// @Summary Delete dormitory
// @Description Delete a dormitory (soft delete)
// @Tags dormitories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dormitory ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dormitories/{id} [delete]
func (h *DormitoryHandler) DeleteDormitory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ErrorBadRequest(c, "Invalid dormitory ID", err.Error())
		return
	}

	err = h.dormitoryUseCase.DeleteDormitory(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domainErrors.ErrDormitoryNotFound:
			response.ErrorNotFound(c, "Dormitory not found")
		default:
			response.ErrorInternalServer(c, "Failed to delete dormitory", err.Error())
		}
		return
	}

	response.SuccessNoContent(c)
}

// ListDormitories handles listing dormitories with pagination
// @Summary List dormitories
// @Description Get paginated list of dormitories
// @Tags dormitories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.ListDormitoriesResponse
// @Failure 400 {object} map[string]string
// @Router /api/dormitories [get]
func (h *DormitoryHandler) ListDormitories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := h.dormitoryUseCase.ListDormitories(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list dormitories", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Dormitories retrieved successfully")
}
