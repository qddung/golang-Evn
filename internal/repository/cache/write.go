package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// write cache
func (c *CacheInstance) Write(ctx context.Context, groupKey string, field string, value []byte, exp time.Duration) error {
	// Implement the logic to write to the cache using the provided key and value
	// Return any error encountered
	_, err := c.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, groupKey, field, value)
		pipe.Expire(ctx, groupKey, exp)
		return nil
	})
	return err
}
