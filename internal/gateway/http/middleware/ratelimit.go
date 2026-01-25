package middleware

import (
	"context"
	"sync"
	"time"

	"example_shop/internal/gateway/http/dto"

	"github.com/cloudwego/hertz/pkg/app"
)

type counter struct {
	sec      int64
	count    int
	lastSeen int64
}

type rateLimiter struct {
	mu           sync.Mutex
	items        map[string]counter
	limit        int
	maxIdleSec   int64
	lastSweepSec int64
}

func newRateLimiter(limit int) *rateLimiter {
	return &rateLimiter{
		items: make(map[string]counter),
		limit: limit,
		// key 为 user_id 或 IP，在真实线上会不断出现新 key；需要做闲置淘汰，避免 map 无限增长
		maxIdleSec: 60,
	}
}

func (r *rateLimiter) allow(key string) bool {
	nowSec := time.Now().Unix()
	r.mu.Lock()
	if r.lastSweepSec == 0 {
		r.lastSweepSec = nowSec
	}
	if nowSec-r.lastSweepSec >= 60 {
		r.sweepLocked(nowSec)
		r.lastSweepSec = nowSec
	}
	v := r.items[key]
	if v.sec != nowSec {
		v.sec = nowSec
		v.count = 0
	}
	v.count++
	v.lastSeen = nowSec
	r.items[key] = v
	r.mu.Unlock()
	return v.count <= r.limit
}

func (r *rateLimiter) sweepLocked(nowSec int64) {
	if r.maxIdleSec <= 0 {
		return
	}
	cutoff := nowSec - r.maxIdleSec
	for k, v := range r.items {
		if v.lastSeen <= cutoff {
			delete(r.items, k)
		}
	}
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
