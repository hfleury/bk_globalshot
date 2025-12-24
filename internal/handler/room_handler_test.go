package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hfleury/bk_globalshot/internal/model"
	mock_services "github.com/hfleury/bk_globalshot/mock/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockRoomService(ctrl)
	handler := NewRoomHandler(mockService)

	room := &model.Room{
		ID:        "room-123",
		Name:      "Living Room",
		UnitID:    "unit-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().CreateRoom(gomock.Any(), "Living Room", "unit-123").Return(room, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := CreateRoomRequest{
		Name:   "Living Room",
		UnitID: "unit-123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonBody))

	handler.CreateRoom(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetAllRooms(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockRoomService(ctrl)
	handler := NewRoomHandler(mockService)

	rooms := []*model.Room{
		{ID: "1", Name: "Room 1", UnitID: "u1"},
		{ID: "2", Name: "Room 2", UnitID: "u1"},
	}

	mockService.EXPECT().GetAllRooms(gomock.Any(), 10, 0).Return(rooms, int64(2), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/rooms", nil)

	handler.GetAllRooms(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Range"), "rooms 0-1/2")
}
