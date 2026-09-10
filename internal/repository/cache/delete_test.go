package cache

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type testCache struct {
	name          string
	closeClient   bool
	expectedError error
}

func assertConnectionClose(tc *testCache, ctx context.Context, t *testing.T, client *redis.Client, groupKey, field, value string) {
	if !tc.closeClient {
		assert.NoError(t, client.HSet(ctx, groupKey, field, value).Err())
	} else {
		assert.NoError(t, client.Close())
	}
}

func TestCacheInstance_Delete(t *testing.T) {
	testCases := []testCache{
		{name: "success", closeClient: false, expectedError: nil},
		{name: "connection error", closeClient: true, expectedError: redis.ErrClosed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			groupKey, field, value := "group", "field", "value"
			cache, client := setupCache(t)
			ctx := context.Background()
			assertConnectionClose(&tc, ctx, t, client, groupKey, field, value)

			err := cache.Delete(ctx, groupKey, field)
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.False(t, client.HExists(ctx, groupKey, field).Val())
			}
		})
	}
}

func TestCacheInstance_DeleteGroup(t *testing.T) {
	testCases := []testCache{
		{name: "success", closeClient: false, expectedError: nil},
		{name: "connection error", closeClient: true, expectedError: redis.ErrClosed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			groupKey, field, value := "group", "field", "value"
			cache, client := setupCache(t)
			ctx := context.Background()
			assertConnectionClose(&tc, ctx, t, client, groupKey, field, value)

			err := cache.DeleteGroup(ctx, groupKey)
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, int64(0), client.Exists(ctx, groupKey).Val())
			}
		})
	}
}
