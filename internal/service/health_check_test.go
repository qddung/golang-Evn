package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Parallel()

	// List of test cases with different mock environment variables and expected results
	testCases := []struct {
		name        string
		mockEnvName string
		mockEnvID   string

		expectedResult HealthStatusResult
	}{
		{
			name:        "Case 1: Cấu hình mặc định",
			mockEnvName: "bookmark_service",
			mockEnvID:   "fixed-uuid-1234",
			expectedResult: HealthStatusResult{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "fixed-uuid-1234",
			},
		},
		{
			name:        "Case 2: Ứng dụng production",
			mockEnvName: "prod_bookmark",
			mockEnvID:   "prod-uuid-9999",
			expectedResult: HealthStatusResult{
				Message:     "OK",
				ServiceName: "prod_bookmark",
				InstanceID:  "prod-uuid-9999",
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()

			// Create a new instance of HealthCheck with the mock environment variables
			svc := NewHealthCheck(tc.mockEnvName, tc.mockEnvID)

			// Execute Ping và nhận kết quả
			result, err := svc.Ping()

			// Check if the result matches the expected result
			assert.NoError(testItem, err)
			assert.Equal(testItem, tc.expectedResult, result)
		})
	}
}
