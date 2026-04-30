package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb *redis.Client
}

func NewLimiter(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb}
}

const (
	capacity   = 5               // max tokens
	refillRate = 1               // tokens per second
)

func (l *Limiter) Allow(ctx context.Context, userID string) (bool, int) {
	key := "rate:" + userID

	now := time.Now().Unix()

	data, _ := l.rdb.HGetAll(ctx, key).Result()

	tokens := capacity
	last := now

	if len(data) > 0 {
		tokens, _ = strconv.Atoi(data["tokens"])
		last, _ = strconv.ParseInt(data["last"], 10, 64)
	}

	// calculate refill
	elapsed := int(now - last)
	if elapsed > 0 {
		tokens = min(capacity, tokens+(elapsed*refillRate))
		last = now 
	}

	if tokens == 0 {
		return false, 0
	}

	// consume token
	tokens--

	// store updated values
	l.rdb.HSet(ctx, key, map[string]interface{}{
		"tokens": tokens,
		"last":   last,
	})

	l.rdb.Expire(ctx, key, time.Minute)

	return true, tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}