package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type UnitRouter struct {
	handler *handler.UnitHandler
}

func NewUnitRouter(handler *handler.UnitHandler) *UnitRouter {
	return &UnitRouter{handler: handler}
}

func (r *UnitRouter) SetupUnitRouter(config *gin.RouterGroup) {
	routes := config.Group("/units")
	{
		routes.POST("", r.handler.CreateUnit)
		routes.GET("", r.handler.GetAllUnits)
		routes.GET("/:id", r.handler.GetUnitByID)
		routes.PUT("/:id", r.handler.UpdateUnit)
		routes.DELETE("/:id", r.handler.DeleteUnit)
	}
}
