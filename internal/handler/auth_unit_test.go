package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_services "github.com/hfleury/bk_globalshot/mock/services"
	"github.com/stretchr/testify/assert"
)

func TestLoging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedAuthService := mock_services.NewMockAuthService(ctrl)
	authHandler := NewAuthHandler(mockedAuthService)
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 26)

	// Using table-driven test for easy HHTP status check
	tests := []struct {
		name               string
		requestBody        LoginRequest
		onAuthService      func()
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name: "SUCCESS - Valid login",
			requestBody: LoginRequest{
				Email:    email,
				Password: password,
			},
			onAuthService: func() {
				mockedAuthService.EXPECT().
					Login(gomock.Any(), email, password).
					Return("valid-token", true, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: gin.H{
				"message": "Login successful",
				"token":   "valid-token",
			},
		},
		{
			name: "FAIL - Invalid email format",
			requestBody: LoginRequest{
				Email:    "not-an-email",
				Password: password,
			},
			onAuthService: func() {
				// No expectation needed it fails before reaching service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: gin.H{
				"error": "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
			},
		},
		{
			name: "FAIL - Invalid credentials",
			requestBody: LoginRequest{
				Email:    email,
				Password: password,
			},
			onAuthService: func() {
				mockedAuthService.EXPECT().
					Login(gomock.Any(), email, password).
					Return("", false, nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: gin.H{
				"error": "Invalid credentials",
			},
		},
		{
			name: "FAIL - Internal server error",
			requestBody: LoginRequest{
				Email:    email,
				Password: password,
			},
			onAuthService: func() {
				mockedAuthService.EXPECT().
					Login(gomock.Any(), email, password).
					Return("", false, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: gin.H{
				"error": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.onAuthService()

			recorder := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(reqBody))
			req = req.WithContext(context.Background())
			r := gin.Default()
			r.POST("/v1/auth/login", authHandler.Login)

			r.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, resp)

		})
	}
}
