package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/service"
)

type RoomHandler struct {
	service service.RoomService
}

func NewRoomHandler(service service.RoomService) *RoomHandler {
	return &RoomHandler{service: service}
}

type CreateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
	UnitID string `json:"unit_id" binding:"required"`
}

type UpdateRoomRequest struct {
	Name   string `json:"name" binding:"required"`
	UnitID string `json:"unit_id" binding:"required"`
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name/unit_id", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	room, err := h.service.CreateRoom(c.Request.Context(), req.Name, req.UnitID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess("Room created successfully", room))
}

func (h *RoomHandler) GetAllRooms(c *gin.Context) {
	limit := 10
	offset := 0

	rangeParam := c.Query("range")
	if rangeParam != "" {
		var rangeSlice []int
		if err := json.Unmarshal([]byte(rangeParam), &rangeSlice); err == nil && len(rangeSlice) == 2 {
			offset = rangeSlice[0]
			limit = rangeSlice[1] - rangeSlice[0] + 1
		}
	}

	filterParam := c.Query("filter")
	var unitID string
	if filterParam != "" {
		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(filterParam), &filter); err == nil {
			if val, ok := filter["unit_id"].(string); ok {
				unitID = val
			}
		}
	}

	rooms, total, err := h.service.GetAllRooms(c.Request.Context(), limit, offset, unitID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	end := offset + len(rooms) - 1
	if len(rooms) == 0 {
		end = offset
	}
	c.Header("Content-Range", fmt.Sprintf("rooms %d-%d/%d", offset, end, total))

	c.JSON(http.StatusOK, dto.ResponseSuccess("Rooms retrieved successfully", rooms))
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Room not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Room retrieved successfully", room))
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationError("name/unit_id", "Invalid input", dto.ErrorCodeValidationFailed))
		return
	}

	room, err := h.service.UpdateRoom(c.Request.Context(), id, req.Name, req.UnitID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, dto.ResponseError("Room not found", nil))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Room updated successfully", room))
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteRoom(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess("Room deleted successfully", nil))
}
