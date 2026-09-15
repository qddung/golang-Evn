package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCacheInstance_Write(t *testing.T) {
	testCases := []struct {
		name          string
		closeClient   bool
		expectedError error
	}{
		{name: "success"},
		{name: "connection error", closeClient: true, expectedError: redis.ErrClosed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache, client := setupCache(t)
			ctx := context.Background()
			if tc.closeClient {
				assert.NoError(t, client.Close())
			}

			err := cache.Write(ctx, "group", "field", []byte("value"), time.Minute)
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				value, readErr := client.HGet(ctx, "group", "field").Result()
				assert.NoError(t, readErr)
				assert.Equal(t, "value", value)
			}
		})
	}
}
