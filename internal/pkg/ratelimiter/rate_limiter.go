package ratelimiter

import (
	"context"

	"github.com/brunocordeiro180/go-rate-limiter/internal/infra/database"
)

type RateLimiter struct {
	repository database.RateLimiterRepository
	rateIp     int
	rateToken  int
	duration   int
}

func NewRateLimiter(repository database.RateLimiterRepository, ip, token, duration int) *RateLimiter {
	rateIp := ip
	rateToken := token
	rateDuration := duration
	return &RateLimiter{repository: repository, rateIp: rateIp, rateToken: rateToken, duration: rateDuration}
}

func (l *RateLimiter) Check(ctx context.Context, key string, isToken bool) bool {
	limit := l.rateIp
	if isToken {
		limit = l.rateToken
	}

	count, err := l.repository.Increment(ctx, key, 1)
	if err != nil {
		return false
	}

	if count > limit {
		l.repository.Increment(ctx, key+":blocked", l.duration)
		return false
	}

	return true
}

func (l *RateLimiter) IsBlocked(ctx context.Context, key string) bool {
	blocked, err := l.repository.Get(ctx, key+":blocked")
	return err == nil && blocked > 0
}
