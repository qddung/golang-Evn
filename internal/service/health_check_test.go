package service

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Parallel()

	// List of test cases with different mock environment variables and expected results
	testCases := []struct {
		name           string
		mockEnvName    string
		mockEnvID      string
		setupPingRepo  func(ctx context.Context) *mocks.Ping
		expectedResult HealthStatusResult
		expectedErr    error
	}{
		{
			name:        "Normal case",
			mockEnvName: "bookmark_service",
			mockEnvID:   "fixed-uuid-1234",
			setupPingRepo: func(ctx context.Context) *mocks.Ping {
				pingRepo := mocks.NewPing(t)
				pingRepo.On("Ping", ctx).Return(nil)
				return pingRepo
			},
			expectedResult: HealthStatusResult{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "fixed-uuid-1234",
			},
			expectedErr: nil,
		},
		{
			name:        "Redis connection error",
			mockEnvName: "prod_bookmark",
			mockEnvID:   "prod-uuid-9999",
			setupPingRepo: func(ctx context.Context) *mocks.Ping {
				pingRepo := mocks.NewPing(t)
				pingRepo.On("Ping", ctx).Return(errRedisNotAvailable)
				return pingRepo
			},
			expectedResult: HealthStatusResult{
				Message:     "Error",
				ServiceName: "prod_bookmark",
				InstanceID:  "prod-uuid-9999",
			},
			expectedErr: errRedisNotAvailable,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			ctx := context.Background()
			pingRepo := tc.setupPingRepo(ctx)
			// Create a new instance of HealthCheck with the mock environment variables
			svc := NewHealthCheck(tc.mockEnvName, tc.mockEnvID, pingRepo)

			// Execute Ping và nhận kết quả
			result, err := svc.Ping(ctx)

			// Check if the result matches the expected result
			assert.ErrorIs(testItem, err, tc.expectedErr)
			assert.Equal(testItem, tc.expectedResult, result)
		})
	}
}
