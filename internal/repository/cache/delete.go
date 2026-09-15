package cache

import "context"

// Implement the logic to delete from the cache using the provided key
func (c *CacheInstance) Delete(ctx context.Context, groupKey string, field string) error {
	// Return any error encountered
	err := c.redis.HDel(ctx, groupKey, field).Err() // Example of deleting a value from Redis
	if err != nil {
		return err
	}
	return nil
}

// Implement the logic to delete a group from the cache using the provided group key
func (c *CacheInstance) DeleteGroup(ctx context.Context, groupKey string) error {
	// Return any error encountered
	err := c.redis.Del(ctx, groupKey).Err()
	if err != nil {
		return err
	}
	return nil
}
