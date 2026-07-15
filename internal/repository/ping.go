package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name Ping --filename ping.go
type Ping interface {
	Ping(ctx context.Context) error
}

type ping struct {
	c *redis.Client
}

// NewPing creates a new Ping
func NewPing(c *redis.Client) Ping {
	return &ping{
		c: c,
	}
}

// Ping pings redis
func (c *ping) Ping(ctx context.Context) error {
	return c.c.Ping(ctx).Err()
}
