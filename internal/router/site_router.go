package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type SiteRouter struct {
	handler *handler.SiteHandler
}

func NewSiteRouter(handler *handler.SiteHandler) *SiteRouter {
	return &SiteRouter{handler: handler}
}

func (r *SiteRouter) SetupSiteRouter(config *gin.RouterGroup) {
	routes := config.Group("/sites")
	{
		routes.POST("", r.handler.CreateSite)
		routes.GET("", r.handler.GetAllSites)
		routes.GET("/:id", r.handler.GetSiteByID)
		routes.PUT("/:id", r.handler.UpdateSite)
		routes.DELETE("/:id", r.handler.DeleteSite)
	}
}
