package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=Cache --output=./mocks --filename=cache_mock.go
type Cache interface {
	// read cache
	Read(ctx context.Context, groupKey string, field string) (string, error)

	// write cache
	Write(ctx context.Context, groupKey string, field string, value []byte, exp time.Duration) error

	// delete cache
	Delete(ctx context.Context, groupKey string, field string) error

	// delete group
	DeleteGroup(ctx context.Context, groupKey string) error
}
type CacheInstance struct {
	redis *redis.Client
}

func NewCache(redis *redis.Client) Cache {
	return &CacheInstance{redis}
}
