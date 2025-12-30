package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required"`
	CompanyID string `json:"company_id"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Role      string `json:"role" binding:"required"`
	CompanyID string `json:"company_id"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("email/password/role", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req.Email, req.Password, req.Role, req.CompanyID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("User created successfully", user))
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	limit := 10
	offset := 0

	users, total, err := h.service.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.Header("Content-Range", fmt.Sprintf("users %d-%d/%d", offset, offset+len(users)-1, total))
	c.JSON(http.StatusOK, dto.ResponseSuccess("Users retrieved successfully", users))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("User not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("User retrieved successfully", user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("email/role", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, req.Email, req.Role, req.CompanyID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("User not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("User updated successfully", user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("User deleted successfully", nil))
}
