package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client    *redis.Client
	burst     int64
	windowDur time.Duration
}

func NewRedisRateLimiter(client *redis.Client, burst int64, windowDur time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:    client,
		burst:     burst,
		windowDur: windowDur,
	}
}

func (rl *RedisRateLimiter) Allow(key string) (bool, int64, int64) {
	ctx := context.Background()
	window := time.Now().Unix() / int64(rl.windowDur.Seconds())
	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, window)

	count, err := rl.client.Incr(ctx, redisKey).Result()
	if err != nil {
		return true, rl.burst, rl.burst
	}

	if count == 1 {
		rl.client.Expire(ctx, redisKey, rl.windowDur*2)
	}

	if count > rl.burst {
		return false, 0, rl.burst
	}

	remaining := rl.burst - count
	if remaining < 0 {
		remaining = 0
	}
	return true, remaining, rl.burst
}
