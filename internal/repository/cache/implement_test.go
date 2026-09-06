package cache

import (
	"context"
	"testing"
	"time"

	redis_helper "github.com/homework/lab/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupCache(t *testing.T) (*CacheInstance, *redis.Client) {
	t.Helper()
	client := redis_helper.InitMockRedis(t)
	return &CacheInstance{redis: client}, client
}

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

func TestCacheInstance_Delete(t *testing.T) {
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
			if !tc.closeClient {
				assert.NoError(t, client.HSet(ctx, "group", "field", "value").Err())
			} else {
				assert.NoError(t, client.Close())
			}

			err := cache.Delete(ctx, "group", "field")
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.False(t, client.HExists(ctx, "group", "field").Val())
			}
		})
	}
}

func TestCacheInstance_DeleteGroup(t *testing.T) {
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
			if !tc.closeClient {
				assert.NoError(t, client.HSet(ctx, "group", "field", "value").Err())
			} else {
				assert.NoError(t, client.Close())
			}

			err := cache.DeleteGroup(ctx, "group")
			assert.ErrorIs(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, int64(0), client.Exists(ctx, "group").Val())
			}
		})
	}
}
