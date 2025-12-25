package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/router/middleware"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type CompanyHandler struct {
	service service.CompanyService
}

func NewCompanyHandler(service service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

type CreateCompanyRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name", "Name is required", dto.ErrorCodeValidationFailed))
		return
	}

	company, err := h.service.CreateCompany(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("Company created successfully", company))
}

func (h *CompanyHandler) GetAllCompanies(c *gin.Context) {
	limit := 10
	offset := 0

	rangeParam := c.Query("range")
	if rangeParam != "" {
		var rangeSlice []int
		if err := json.Unmarshal([]byte(rangeParam), &rangeSlice); err == nil && len(rangeSlice) == 2 {
			offset = rangeSlice[0]
			limit = rangeSlice[1] - rangeSlice[0] + 1
		}
	}

	payload := middleware.GetAuthPayload(c)
	if payload != nil && model.Role(payload.Role) == model.RoleCompany {
		company, err := h.service.GetCompanyByID(c.Request.Context(), payload.CompanyID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
			return
		}

		var companies []*model.Company
		var total int64 = 0
		if company != nil {
			companies = append(companies, company)
			total = 1
		}

		c.Header("Content-Range", fmt.Sprintf("companies 0-0/%d", total))
		c.JSON(http.StatusOK, dto.ResponseSuccess("Companies retrieved successfully", companies))
		return
	}

	companies, total, err := h.service.GetAllCompanies(c.Request.Context(), limit, offset)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	// Set Content-Range header for React Admin
	// Format: companies start-end/total
	end := offset + len(companies) - 1
	if len(companies) == 0 {
		end = offset // or 0? usually if empty, start-start/total where start usually > total, or */total
		// React Admin handles empty content info gracefully usually if we say 0-0/0 or similar.
		// Let's set it to: companies offset-(offset + length)/total
	}
	// Content-Range: companies 0-9/100
	c.Header("Content-Range", fmt.Sprintf("companies %d-%d/%d", offset, end, total))

	c.JSON(http.StatusOK, dto.ResponseSuccess("Companies retrieved successfully", companies))
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	id := c.Param("id")

	payload := middleware.GetAuthPayload(c)
	if payload != nil && model.Role(payload.Role) == model.RoleCompany {
		if payload.CompanyID != id {
			c.JSON(http.StatusForbidden, dto.ForbiddenResponse("Access denied"))
			return
		}
	}

	company, err := h.service.GetCompanyByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if company == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Company not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Company retrieved successfully", company))
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id := c.Param("id")

	payload := middleware.GetAuthPayload(c)
	if payload != nil && model.Role(payload.Role) == model.RoleCompany {
		if payload.CompanyID != id {
			c.JSON(http.StatusForbidden, dto.ForbiddenResponse("Access denied"))
			return
		}
	}

	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name", "Name is required", dto.ErrorCodeValidationFailed))
		return
	}

	company, err := h.service.UpdateCompany(c.Request.Context(), id, req.Name)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if company == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Company not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Company updated successfully", company))
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id := c.Param("id")

	payload := middleware.GetAuthPayload(c)
	if payload != nil && model.Role(payload.Role) == model.RoleCompany {
		// Company cannot delete itself? Or allows it?
		// Usually deletion is admin only or restricted.
		// Let's allow strictly verification but arguably companies shouldn't delete themselves easily.
		if payload.CompanyID != id {
			c.JSON(http.StatusForbidden, dto.ForbiddenResponse("Access denied"))
			return
		}
	}

	err := h.service.DeleteCompany(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Company deleted successfully", nil))
}
