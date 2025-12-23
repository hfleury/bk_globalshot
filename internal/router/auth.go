package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type AuthRouter struct {
	handler *handler.AuthHandler
}

func NewAuthRouter(handler *handler.AuthHandler) *AuthRouter {
	return &AuthRouter{
		handler: handler,
	}
}

func (ar *AuthRouter) SetupAuthRouter(api *gin.RouterGroup) {

	auth := api.Group("/auth")
	{
		auth.POST("/login", ar.handler.Login)
		auth.POST("/reset-password", ar.handler.ResetPassword)
	}
}
