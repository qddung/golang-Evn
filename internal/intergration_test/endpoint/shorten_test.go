package endpoint // Để _test để đảm bảo tính đóng gói độc lập

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestShorten_Integration(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string

		// setupTestHTTP là một hàm để thiết lập và gửi yêu cầu HTTP đến API Engine, trả về ResponseRecorder để kiểm tra kết quả
		setupTestHTTP func(router api.Engine) *httptest.ResponseRecorder

		expectedStatusCode         int
		getExpectedResponseContain func() string
		configTest                 *config.Config
	}{
		{
			name: "Create shorten link successfully",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{
					"url": "https://www.google.com",
					"exp": 3600,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusOK,
			getExpectedResponseContain: func() string {
				return "Shorten URL generated successfully!"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
		{
			name: "Create shorten link fail - invalid url format",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{
					"url": "invalid-url-not-a-link",
					"exp": 3600,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			getExpectedResponseContain: func() string {
				return "Field validation for 'Url' failed on the 'url' tag"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
		{
			name: "Create shorten link fail - missing required exp",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{
					"url": "https://www.google.com",
				}
				bodyBytes, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			getExpectedResponseContain: func() string {
				return "Field validation for 'Exp' failed on the 'required' tag"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
		{
			name: "Create shorten link fail - exp exceed maximum limit (604800)",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{
					"url": "https://www.google.com",
					"exp": 999999,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			getExpectedResponseContain: func() string {
				return "Field validation for 'Exp' failed on the 'lte' tag"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
		{
			name: "Create shorten link fail - invalid json body",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBufferString("{invalid-json}"))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			getExpectedResponseContain: func() string {
				return "invalid character"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
		{
			name: "Wrong shorten endpoint",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{
					"url": "https://www.google.com",
					"exp": 3600,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten_not_found", bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				router.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusNotFound,
			getExpectedResponseContain: func() string {
				return "404 page not found"
			},
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()

			fmt.Printf("Loaded config: %+v\n", tc.configTest)
			redisMock := redisPkg.InitMockRedis(testItem)
			apiEngine := api.NewEngine(tc.configTest, redisMock)

			rec := tc.setupTestHTTP(apiEngine)

			// Check the status code of the response
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code, "Expected status code does not match actual status code")
			// Check the response body content
			assert.Contains(testItem, rec.Body.String(), tc.getExpectedResponseContain(), "Expected response body does not match actual response body")

		})
	}
}
