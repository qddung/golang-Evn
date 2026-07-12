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

func NewPing(c *redis.Client) Ping {
	return &ping{
		c: c,
	}
}

func (c *ping) Ping(ctx context.Context) error {
	return c.c.Ping(ctx).Err()
}
