package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hfleury/bk_globalshot/internal/model"
	mock_services "github.com/hfleury/bk_globalshot/mock/services"
	"github.com/hfleury/bk_globalshot/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestCompanyHandler_Isolation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		userRole        string
		userCompanyID   string
		targetCompanyID string
		setupService    func(*mock_services.MockCompanyService)
		expectedStatus  int
	}{
		{
			name:            "Company user accessing own company - Allowed",
			userRole:        string(model.RoleCompany),
			userCompanyID:   "company-123",
			targetCompanyID: "company-123",
			setupService: func(s *mock_services.MockCompanyService) {
				s.EXPECT().GetCompanyByID(gomock.Any(), "company-123").Return(&model.Company{ID: "company-123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "Company user accessing other company - Forbidden",
			userRole:        string(model.RoleCompany),
			userCompanyID:   "company-123",
			targetCompanyID: "company-456",
			setupService: func(s *mock_services.MockCompanyService) {
				// Service should NOT be called
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:            "Admin user accessing any company - Allowed",
			userRole:        string(model.RoleAdmin),
			userCompanyID:   "",
			targetCompanyID: "company-456",
			setupService: func(s *mock_services.MockCompanyService) {
				s.EXPECT().GetCompanyByID(gomock.Any(), "company-456").Return(&model.Company{ID: "company-456"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_services.NewMockCompanyService(ctrl)
			tt.setupService(mockService)
			handler := NewCompanyHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Simulate middleware setting payload
			payload := &token.Payload{
				Role:      tt.userRole,
				CompanyID: tt.userCompanyID,
			}
			c.Set("authorization_payload", payload)

			// Set params
			c.Params = []gin.Param{{Key: "id", Value: tt.targetCompanyID}}
			c.Request = httptest.NewRequest("GET", "/companies/"+tt.targetCompanyID, nil)

			handler.GetCompanyByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
