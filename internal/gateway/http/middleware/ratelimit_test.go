package middleware

import (
	"bytes"
	"context"
	"testing"
	"time"

	"example_shop/internal/gateway/http/dto"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func waitNextSecond() {
	now := time.Now().Unix()
	for time.Now().Unix() == now {
		time.Sleep(2 * time.Millisecond)
	}
}

func TestRateLimit(t *testing.T) {
	waitNextSecond()

	h := server.Default()
	h.Use(RateLimit(1))
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, dto.BaseHTTPResp{Code: 200, Msg: "ok"})
	})

	w1 := ut.PerformRequest(h.Engine, "GET", "/ping", &ut.Body{Body: bytes.NewBuffer(nil), Len: 0})
	if w1.Result().StatusCode() != 200 {
		t.Fatalf("unexpected status: %d", w1.Result().StatusCode())
	}

	w2 := ut.PerformRequest(h.Engine, "GET", "/ping", &ut.Body{Body: bytes.NewBuffer(nil), Len: 0})
	if w2.Result().StatusCode() != 429 {
		t.Fatalf("unexpected status: %d body=%s", w2.Result().StatusCode(), string(w2.Result().Body()))
	}
}

