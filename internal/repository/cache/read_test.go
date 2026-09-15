package cache

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCacheInstance_Read(t *testing.T) {
	testCases := []struct {
		name          string
		closeClient   bool
		expectedValue string
		expectedError error
	}{
		{
			name:          "success",
			expectedValue: "value",
		},
		{
			name:          "connection error",
			closeClient:   true,
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache, client := setupCache(t)
			ctx := context.Background()
			if !tc.closeClient {
				assert.NoError(t, client.HSet(ctx, "group", "field", "value").Err())
			} else {
				assert.NoError(t, client.Close())
			}

			value, err := cache.Read(ctx, "group", "field")
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}
