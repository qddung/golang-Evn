package health_check_repository

import (
	"context"
	"testing"

	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestPing_Ping(t *testing.T) {
	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		expectedErr error
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
		},
		{
			name: "err case - redis connection err",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedErr: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			redisMock := tc.setupMock(ctx)
			testRepo := NewPing(redisMock)
			err := testRepo.Ping(ctx)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
