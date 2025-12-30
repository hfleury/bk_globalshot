package router

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
	"github.com/hfleury/bk_globalshot/internal/router/middleware"
	"github.com/hfleury/bk_globalshot/pkg/token"
)

type Router struct {
	eng *gin.Engine
}

func NewRouter(eng *gin.Engine) *Router {
	return &Router{
		eng: eng,
	}
}

func (r *Router) SetupRouter(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	companyHandler *handler.CompanyHandler,
	roomHandler *handler.RoomHandler, // Added
	siteHandler *handler.SiteHandler, // Added
	unitHandler *handler.UnitHandler, // Added
	userHandler *handler.UserHandler, // Added
	tokenMaker token.Maker,
) {

	// CORS Configuration
	config := cors.DefaultConfig()
	// Allow localhost by default
	allowOrigins := []string{"http://localhost:5173", "http://127.0.0.1:5173"}

	// Add ALLOWED_ORIGIN from env (e.g., Render Frontend URL)
	if envOrigin := os.Getenv("ALLOWED_ORIGIN"); envOrigin != "" {
		allowOrigins = append(allowOrigins, envOrigin)
	}

	config.AllowOrigins = allowOrigins
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

		// Private routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(tokenMaker))
		{
			companyRouter := NewCompanyRouter(companyHandler)
			companyRouter.SetupCompanyRouter(protected)

			roomRouter := NewRoomRouter(roomHandler)
			roomRouter.SetupRoomRouter(protected)

			siteRouter := NewSiteRouter(siteHandler)
			siteRouter.SetupSiteRouter(protected)

			unitRouter := NewUnitRouter(unitHandler)
			unitRouter.SetupUnitRouter(protected)

			userRouter := NewUserRouter(userHandler)
			userRouter.SetupUserRouter(protected)
		}
	}
}
