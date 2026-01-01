package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type UnitHandler struct {
	service service.UnitService
}

func NewUnitHandler(service service.UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

type CreateUnitRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	SiteID   string  `json:"site_id" binding:"required"`
	ClientID *string `json:"client_id"`
}

type UpdateUnitRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	SiteID   string  `json:"site_id" binding:"required"`
	ClientID *string `json:"client_id"`
}

func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var req CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name/type/site_id", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	unit, err := h.service.CreateUnit(c.Request.Context(), req.Name, req.Type, req.SiteID, req.ClientID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("Unit created successfully", unit))
}

func (h *UnitHandler) BatchCreate(c *gin.Context) {
	var req []service.BatchCreateUnitItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("body", "Invalid input list", dto.ErrorCodeValidationFailed))
		return
	}

	units, err := h.service.BatchCreateUnits(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("Units created successfully", units))
}

func (h *UnitHandler) GetAllUnits(c *gin.Context) {
	limit := 10
	offset := 0

	// Handle pagination if needed
	units, total, err := h.service.GetAllUnits(c.Request.Context(), limit, offset)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.Header("Content-Range", fmt.Sprintf("units %d-%d/%d", offset, offset+len(units)-1, total))
	c.JSON(http.StatusOK, dto.ResponseSuccess("Units retrieved successfully", units))
}

func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.service.GetUnitByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if unit == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Unit not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Unit retrieved successfully", unit))
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name/type/site_id", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	unit, err := h.service.UpdateUnit(c.Request.Context(), id, req.Name, req.Type, req.SiteID, req.ClientID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if unit == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Unit not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Unit updated successfully", unit))
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteUnit(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Unit deleted successfully", nil))
}
