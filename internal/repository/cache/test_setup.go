package cache

import (
	"testing"

	redis_helper "github.com/homework/lab/pkg/redis"
	"github.com/redis/go-redis/v9"
)

func setupCache(t *testing.T) (*CacheInstance, *redis.Client) {
	t.Helper()
	client := redis_helper.InitMockRedis(t)
	return &CacheInstance{redis: client}, client
}
