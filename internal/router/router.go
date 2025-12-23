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
	// Allow localhost by default
	allowOrigins := []string{"http://localhost:5173", "http://127.0.0.1:5173"}

	// Add ALLOWED_ORIGIN from env (e.g., Render Frontend URL)
	/*
	   Note: We need os package to check environment variable.
	   Since we cannot easily add imports here without potentially breaking if 'os' is not imported,
	   I will use a hardcoded check or ensure os is imported.
	   Actually, I should add "os" to imports first.
	   Let's assume I will add "os" import in a separate step or just do it here if possible?
	   Wait, I can't add imports easily with replace_file_content if they are far away.
	   Let's stick to adding the logic, but I need 'os'.
	   I will verify imports first.
	*/
	// Reverting to simple replacement that assumes 'os' is imported or I will use a helper.
	// Let's modify imports first.
	config.AllowOrigins = allowOrigins // Placeholder until I add 'os'
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
