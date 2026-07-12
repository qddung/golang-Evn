package redis

import "github.com/redis/go-redis/v9"

func NewRedisClient(cfg *config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
