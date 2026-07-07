package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/service"
	mocks_health_check "github.com/homework/lab/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		mockEnvName        string
		mockEnvID          string
		expectedStatusCode int
		expectedResult     HealthResponse
	}{
		{
			name:               "Case 1: Cấu hình mặc định",
			mockEnvName:        "bookmark_service",
			mockEnvID:          "fixed-uuid-1234",
			expectedStatusCode: http.StatusOK,
			expectedResult: HealthResponse{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "fixed-uuid-1234",
			},
		},
		{
			name:               "Case 2: Ứng dụng production",
			mockEnvName:        "prod_bookmark",
			mockEnvID:          "prod-uuid-9999",
			expectedStatusCode: http.StatusOK,
			expectedResult: HealthResponse{
				Message:     "OK",
				ServiceName: "prod_bookmark",
				InstanceID:  "prod-uuid-9999",
			},
		},
	}

	for _, tc := range testCases {
		// Preserve the test context for the parallel loop
		tc := tc

		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()

			// Initialize a mocked Gin HTTP context
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/health_check", nil)

			// Create a mock service
			mockSvc := mocks_health_check.NewHealthCheck(testItem) // Use testItem for the correct context

			// Configure the mock behavior
			// When the handler calls mockSvc.Ping(), the mock returns the data for the current test case
			mockSvc.On("Ping").Return(service.HealthStatusResult{
				Message:     tc.expectedResult.Message,
				ServiceName: tc.mockEnvName,
				InstanceID:  tc.mockEnvID,
			}, nil)

			// Initialize the handler and inject the mock service (DI)
			healthCheckHandler := NewHealthCheck(mockSvc)

			// Execute the function under test
			healthCheckHandler.Ping(c)

			// Compare the returned JSON response accurately
			var actualResponse HealthResponse
			// Decode the JSON response body into actualResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			assert.NoError(testItem, err)
			// Check the HTTP status code
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code)
			// Compare the two structs directly
			assert.Equal(testItem, tc.expectedResult, actualResponse)
		})
	}
}
