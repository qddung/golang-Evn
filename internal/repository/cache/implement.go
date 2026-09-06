package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// read cache
func (c *CacheInstance) Read(ctx context.Context, groupKey string, field string) (string, error) {
	// Implement the logic to read from the cache using the provided key
	// Return the value and any error encountered
	value, err := c.redis.HGet(ctx, groupKey, field).Result() // Example of getting a value from Redis
	if err != nil {
		return "", err
	}
	return value, nil
}

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

// delete cache
func (c *CacheInstance) Delete(ctx context.Context, groupKey string, field string) error {
	// Implement the logic to delete from the cache using the provided key
	// Return any error encountered
	err := c.redis.HDel(ctx, groupKey, field).Err() // Example of deleting a value from Redis
	if err != nil {
		return err
	}
	return nil
}

// delete group
func (c *CacheInstance) DeleteGroup(ctx context.Context, groupKey string) error {
	// Implement the logic to delete a group from the cache using the provided group key
	// Return any error encountered
	err := c.redis.Del(ctx, groupKey).Err()
	if err != nil {
		return err
	}
	return nil
}
