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

// SKSDefinitionHandler handles HTTP requests for SKS definitions.
type SKSDefinitionHandler struct {
	definitionUseCase *usecase.SKSDefinitionUseCase
}

// NewSKSDefinitionHandler constructs a handler.
func NewSKSDefinitionHandler(definitionUseCase *usecase.SKSDefinitionUseCase) *SKSDefinitionHandler {
	return &SKSDefinitionHandler{definitionUseCase: definitionUseCase}
}

// CreateSKSDefinition handles POST /api/sks.
func (h *SKSDefinitionHandler) CreateSKSDefinition(c *gin.Context) {
	var req dto.CreateSKSDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	definition, err := h.definitionUseCase.CreateSKSDefinition(c.Request.Context(), req)
	if err != nil {
		h.handleDefinitionError(c, err, "create")
		return
	}

	response.SuccessCreated(c, definition, "SKS definition created successfully")
}

// GetSKSDefinition handles GET /api/sks/:id.
func (h *SKSDefinitionHandler) GetSKSDefinition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS definition ID", err.Error())
		return
	}

	definition, err := h.definitionUseCase.GetSKSDefinition(c.Request.Context(), id)
	if err != nil {
		if err == domainErrors.ErrSKSDefinitionNotFound {
			response.ErrorNotFound(c, "SKS definition not found")
		} else {
			response.ErrorInternalServer(c, "Failed to get SKS definition", err.Error())
		}
		return
	}

	response.SuccessOK(c, definition, "SKS definition retrieved successfully")
}

// ListSKSDefinitions handles GET /api/sks.
func (h *SKSDefinitionHandler) ListSKSDefinitions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	fanID := c.Query("fan_id")

	definitions, err := h.definitionUseCase.ListSKSDefinitions(c.Request.Context(), fanID, page, pageSize)
	if err != nil {
		if err == domainErrors.ErrBadRequest {
			response.ErrorBadRequest(c, "Invalid filters", err.Error())
		} else {
			response.ErrorInternalServer(c, "Failed to list SKS definitions", err.Error())
		}
		return
	}

	response.SuccessOK(c, definitions, "SKS definitions retrieved successfully")
}

// UpdateSKSDefinition handles PUT /api/sks/:id.
func (h *SKSDefinitionHandler) UpdateSKSDefinition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS definition ID", err.Error())
		return
	}

	var req dto.UpdateSKSDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidation(c, err)
		return
	}

	definition, err := h.definitionUseCase.UpdateSKSDefinition(c.Request.Context(), id, req)
	if err != nil {
		h.handleDefinitionError(c, err, "update")
		return
	}

	response.SuccessOK(c, definition, "SKS definition updated successfully")
}

// DeleteSKSDefinition handles DELETE /api/sks/:id.
func (h *SKSDefinitionHandler) DeleteSKSDefinition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid SKS definition ID", err.Error())
		return
	}

	if err := h.definitionUseCase.DeleteSKSDefinition(c.Request.Context(), id); err != nil {
		if err == domainErrors.ErrSKSDefinitionNotFound {
			response.ErrorNotFound(c, "SKS definition not found")
		} else {
			response.ErrorInternalServer(c, "Failed to delete SKS definition", err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SKSDefinitionHandler) handleDefinitionError(c *gin.Context, err error, action string) {
	switch err {
	case domainErrors.ErrFanNotFound:
		response.ErrorNotFound(c, "Fan not found")
	case domainErrors.ErrSubjectNotFound:
		response.ErrorNotFound(c, "Subject not found")
	case domainErrors.ErrSKSDefinitionAlreadyExist:
		response.ErrorConflict(c, "SKS definition already exists", err.Error())
	case domainErrors.ErrSKSDefinitionNotFound:
		response.ErrorNotFound(c, "SKS definition not found")
	case domainErrors.ErrBadRequest:
		response.ErrorBadRequest(c, "Invalid SKS definition data", err.Error())
	default:
		response.ErrorInternalServer(c, "Failed to "+action+" SKS definition", err.Error())
	}
}
