package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.Error(fmt.Errorf("this is a test error from Login handler"))
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, dto.ValidationError(
			"email or password",
			"Invalid request format.",
			dto.ErrorCodeValidationFailed,
		))
		return
	}

	ctx := c.Request.Context()

	token, success, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	if !success {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse("Invalid email or password"))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Login successful", gin.H{
		"token": token,
	}))
}
