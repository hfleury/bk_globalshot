package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type SiteHandler struct {
	service service.SiteService
}

func NewSiteHandler(service service.SiteService) *SiteHandler {
	return &SiteHandler{service: service}
}

type CreateSiteRequest struct {
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address"`
	CompanyID string `json:"company_id" binding:"required"`
}

type UpdateSiteRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name/company_id", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	site, err := h.service.CreateSite(c.Request.Context(), req.Name, req.Address, req.CompanyID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("Site created successfully", site))
}

func (h *SiteHandler) GetAllSites(c *gin.Context) {
	limit := 10
	offset := 0

	// Handle pagination if needed, similar to other handlers
	// Skipping explicit Range header parsing for brevity unless required by frontend framework (Ra-React Admin often uses Range)

	// Retrieve user from context to ensure authentication passed
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ResponseError("Unauthorized", nil))
		return
	}
	// Add user to context for service layer
	ctx := context.WithValue(c.Request.Context(), "user", user.(*model.User))

	sites, total, err := h.service.GetAllSites(ctx, limit, offset)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.Header("Content-Range", fmt.Sprintf("sites %d-%d/%d", offset, offset+len(sites)-1, total))
	c.JSON(http.StatusOK, dto.ResponseSuccess("Sites retrieved successfully", sites))
}

func (h *SiteHandler) GetSiteByID(c *gin.Context) {
	id := c.Param("id")
	site, err := h.service.GetSiteByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Site not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Site retrieved successfully", site))
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	site, err := h.service.UpdateSite(c.Request.Context(), id, req.Name, req.Address)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if site == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Site not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Site updated successfully", site))
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteSite(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Site deleted successfully", nil))
}
