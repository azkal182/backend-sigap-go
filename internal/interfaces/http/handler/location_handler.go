package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/response"
)

// LocationHandler handles read-only location endpoints (public)
type LocationHandler struct {
	useCase *usecase.LocationUseCase
}

func NewLocationHandler(useCase *usecase.LocationUseCase) *LocationHandler {
	return &LocationHandler{useCase: useCase}
}

// helper to parse page & page_size with defaults
func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return page, pageSize
}

// GET /api/provinces
func (h *LocationHandler) ListProvinces(c *gin.Context) {
	page, pageSize := parsePagination(c)
	search := c.Query("search")

	resp, err := h.useCase.ListProvinces(c.Request.Context(), page, pageSize, search)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list provinces", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Provinces retrieved successfully")
}

// GET /api/provinces/:id
func (h *LocationHandler) GetProvince(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid province ID", err.Error())
		return
	}

	resp, err := h.useCase.GetProvinceByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, "Province not found", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Province retrieved successfully")
}

// GET /api/regencies
func (h *LocationHandler) ListRegencies(c *gin.Context) {
	page, pageSize := parsePagination(c)
	search := c.Query("search")

	var provinceID *int
	if v := c.Query("province_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			provinceID = &id
		}
	}

	resp, err := h.useCase.ListRegencies(c.Request.Context(), page, pageSize, provinceID, search)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list regencies", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Regencies retrieved successfully")
}

// GET /api/regencies/:id
func (h *LocationHandler) GetRegency(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid regency ID", err.Error())
		return
	}

	resp, err := h.useCase.GetRegencyByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, "Regency not found", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Regency retrieved successfully")
}

// GET /api/districts
func (h *LocationHandler) ListDistricts(c *gin.Context) {
	page, pageSize := parsePagination(c)
	search := c.Query("search")

	var regencyID *int
	if v := c.Query("regency_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			regencyID = &id
		}
	}

	resp, err := h.useCase.ListDistricts(c.Request.Context(), page, pageSize, regencyID, search)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list districts", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Districts retrieved successfully")
}

// GET /api/districts/:id
func (h *LocationHandler) GetDistrict(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid district ID", err.Error())
		return
	}

	resp, err := h.useCase.GetDistrictByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, "District not found", err.Error())
		return
	}

	response.SuccessOK(c, resp, "District retrieved successfully")
}

// GET /api/villages
func (h *LocationHandler) ListVillages(c *gin.Context) {
	page, pageSize := parsePagination(c)
	search := c.Query("search")

	var districtID *int
	if v := c.Query("district_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			districtID = &id
		}
	}

	resp, err := h.useCase.ListVillages(c.Request.Context(), page, pageSize, districtID, search)
	if err != nil {
		response.ErrorInternalServer(c, "Failed to list villages", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Villages retrieved successfully")
}

// GET /api/villages/:id
func (h *LocationHandler) GetVillage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid village ID", err.Error())
		return
	}

	resp, err := h.useCase.GetVillageByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, "Village not found", err.Error())
		return
	}

	response.SuccessOK(c, resp, "Village retrieved successfully")
}

// Optional: simple health check for location routes
func (h *LocationHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "location service ok"})
}
