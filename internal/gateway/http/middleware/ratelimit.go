package middleware

import (
	"context"
	"sync"
	"time"

	"example_shop/internal/gateway/http/dto"

	"github.com/cloudwego/hertz/pkg/app"
)

type counter struct {
	sec   int64
	count int
}

type rateLimiter struct {
	mu    sync.Mutex
	items map[string]counter
	limit int
}

func newRateLimiter(limit int) *rateLimiter {
	return &rateLimiter{
		items: make(map[string]counter),
		limit: limit,
	}
}

func (r *rateLimiter) allow(key string) bool {
	nowSec := time.Now().Unix()
	r.mu.Lock()
	v := r.items[key]
	if v.sec != nowSec {
		v = counter{sec: nowSec, count: 0}
	}
	v.count++
	r.items[key] = v
	r.mu.Unlock()
	return v.count <= r.limit
}

func RateLimit(maxPerSec int) app.HandlerFunc {
	if maxPerSec <= 0 {
		maxPerSec = 10
	}
	rl := newRateLimiter(maxPerSec)
	return func(ctx context.Context, c *app.RequestContext) {
		key := string(c.Request.Header.Peek("X-User-ID"))
		if key == "" {
			key = c.ClientIP()
		}
		if !rl.allow(key) {
			c.JSON(429, dto.BaseHTTPResp{Code: 429, Msg: "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
