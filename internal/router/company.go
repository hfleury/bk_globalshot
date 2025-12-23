package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type CompanyRouter struct {
	handler *handler.CompanyHandler
}

func NewCompanyRouter(handler *handler.CompanyHandler) *CompanyRouter {
	return &CompanyRouter{
		handler: handler,
	}
}

func (cr *CompanyRouter) SetupCompanyRouter(api *gin.RouterGroup) {
	companies := api.Group("/companies")
	{
		companies.POST("", cr.handler.CreateCompany)
		companies.GET("", cr.handler.GetAllCompanies)
		companies.GET("/:id", cr.handler.GetCompanyByID)
		companies.PUT("/:id", cr.handler.UpdateCompany)
		companies.DELETE("/:id", cr.handler.DeleteCompany)
	}
}
