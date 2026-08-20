package health_check_endpoint // Để _test để đảm bảo tính đóng gói độc lập

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	"github.com/homework/lab/internal/handler/health_check"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func newHealthCheckIntegrationConfig() *config.Config {
	return &config.Config{
		AppPort:     "8080",
		ServiceName: "app_service",
		InstanceID:  "instance_01",
	}
}

func TestHealthCheck_Integration(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		// setupTestHTTP là một hàm để thiết lập và gửi yêu cầu HTTP đến API Engine, trả về ResponseRecorder để kiểm tra kết quả
		setupTestHTTP              func(router api.Engine) *httptest.ResponseRecorder
		expectedStatusCode         int
		getExpectedResponseContain func() string
		configTest                 *config.Config
	}{
		{
			name: "Health check successfully",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusOK,
			getExpectedResponseContain: func() string {
				// Serialize the expected response to JSON format for comparison
				resp, _ := json.Marshal(health_check.HealthResponse{
					Message:     "OK",
					InstanceID:  "instance_01",
					ServiceName: "app_service",
				})
				return string(resp)
			},
			configTest: newHealthCheckIntegrationConfig(),
		},
		{
			name: "Wrong health endpoint",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check_not_found", nil)
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusNotFound,
			getExpectedResponseContain: func() string {
				return "404 page not found"
			},
			configTest: newHealthCheckIntegrationConfig(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			fmt.Printf("Loaded config: %+v\n", tc.configTest)
			fix := fixture.NewUserTestCase(t)
			connectorMock, errConnector := connection.InitDBConnectorMock(testItem, fix)
			if errConnector != nil {
				testItem.Fatal(errConnector)
			}
			apiEngine := api.NewEngine(tc.configTest, connectorMock)
			rec := tc.setupTestHTTP(apiEngine)
			// Check the status code of the response
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code, "Expected status code does not match actual status code")
			// Check the response body content
			assert.Equal(testItem, tc.getExpectedResponseContain(), rec.Body.String(), "Expected response body does not match actual response body")
		})
	}
}
