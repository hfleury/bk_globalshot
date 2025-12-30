package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
)

type UserRouter struct {
	handler *handler.UserHandler
}

func NewUserRouter(handler *handler.UserHandler) *UserRouter {
	return &UserRouter{handler: handler}
}

func (r *UserRouter) SetupUserRouter(config *gin.RouterGroup) {
	routes := config.Group("/users")
	{
		routes.POST("", r.handler.CreateUser)
		routes.GET("", r.handler.GetAllUsers)
		routes.GET("/:id", r.handler.GetUserByID)
		routes.PUT("/:id", r.handler.UpdateUser)
		routes.DELETE("/:id", r.handler.DeleteUser)
	}
}
