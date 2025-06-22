package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type HealthRouter struct {
	handler *handler.HealthHandler
}

func NewHealthRouter(handler *handler.HealthHandler) *HealthRouter {
	return &HealthRouter{
		handler: handler,
	}
}

func (hr *HealthRouter) SetupHealthRouter(api *gin.RouterGroup) {
	health := api.Group("/health")
	{
		health.GET("/health", hr.handler.Check)
	}
}
