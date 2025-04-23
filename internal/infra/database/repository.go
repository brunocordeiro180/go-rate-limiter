package database

import "context"

type RateLimiterRepository interface {
	Increment(ctx context.Context, key string, expiration int) (int, error)
	Get(ctx context.Context, key string) (int, error)
	Reset(ctx context.Context, key string) error
}
