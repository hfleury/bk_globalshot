package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type Router struct {
	eng *gin.Engine
}

func NewRouter(eng *gin.Engine) *Router {
	return &Router{
		eng: eng,
	}
}

func (r *Router) SetupRouter(authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler, companyHandler *handler.CompanyHandler) {
	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"} // Frontend Origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "Range"}
	config.ExposeHeaders = []string{"Content-Range"}
	r.eng.Use(cors.New(config))

	api := r.eng.Group("/v1")
	{
		authRouter := NewAuthRouter(authHandler)
		authRouter.SetupAuthRouter(api)
		healthRouter := NewHealthRouter(healthHandler)
		healthRouter.SetupHealthRouter(api)
		companyRouter := NewCompanyRouter(companyHandler)
		companyRouter.SetupCompanyRouter(api)
	}
}
