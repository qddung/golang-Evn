package cache

import (
	"context"
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
