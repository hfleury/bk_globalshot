package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/handler"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/router/middleware"
)

type RoomRouter struct {
	handler *handler.RoomHandler
}

func NewRoomRouter(handler *handler.RoomHandler) *RoomRouter {
	return &RoomRouter{handler: handler}
}

func (r *RoomRouter) SetupRoomRouter(group *gin.RouterGroup) {
	router := group.Group("/rooms")
	router.Use(middleware.RequireRoles(model.RoleAdmin, model.RoleCompany))
	{
		router.POST("", r.handler.CreateRoom)
		router.GET("", r.handler.GetAllRooms)
		router.GET("/:id", r.handler.GetRoomByID)
		router.PUT("/:id", r.handler.UpdateRoom)
		router.DELETE("/:id", r.handler.DeleteRoom)
	}
}
