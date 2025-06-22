package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type HealthHandler struct {
	service service.HealthService
}

func NewHealthHandler(service service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.AbortWithStatusJSON(503, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "healthy"})
}
