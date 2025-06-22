package router

import (
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

func (r *Router) SetupRouter(authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler) {
	api := r.eng.Group("/v1")
	{

		authRouter := NewAuthRouter(authHandler)
		authRouter.SetupAuthRouter(api)
		healthRouter := NewHealthRouter(healthHandler)
		healthRouter.SetupHealthRouter(api)
	}
}
