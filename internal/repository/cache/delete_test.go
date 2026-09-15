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

func setupInputRunDelete(t *testing.T, tc *testCache, ctx context.Context, groupKey, field, value string) (cacheInstance *CacheInstance, client *redis.Client) {
	t.Parallel()
	cache, client := setupCache(t)
	if !tc.closeClient {
		assert.NoError(t, client.HSet(ctx, groupKey, field, value).Err())
	} else {
		assert.NoError(t, client.Close())
	}
	return cache, client
}

var testCases = []testCache{
	{name: "success", closeClient: false, expectedError: nil},
	{name: "connection error", closeClient: true, expectedError: redis.ErrClosed},
}

func TestCacheInstance_Delete(t *testing.T) {

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			groupKey, field, value := "group", "field", "value"
			cache, client := setupInputRunDelete(t, &tc, ctx, groupKey, field, value)

			err := cache.Delete(ctx, groupKey, field)
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.False(t, client.HExists(ctx, groupKey, field).Val())
			}
		})
	}
}

func TestCacheInstance_DeleteGroup(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			groupKey, field, value := "group", "field", "value"
			cache, client := setupInputRunDelete(t, &tc, ctx, groupKey, field, value)

			err := cache.DeleteGroup(ctx, groupKey)
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, int64(0), client.Exists(ctx, groupKey).Val())
			}
		})
	}
}
