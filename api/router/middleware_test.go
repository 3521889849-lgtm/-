package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example_shop/api/router"

	"github.com/gin-gonic/gin"
)

func TestHealth_HasRequestIDHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if got := w.Header().Get("X-Request-Id"); got == "" {
		t.Fatalf("expected X-Request-Id header to be set")
	}
}

func TestCORSOptions_ReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://example.com" {
		t.Fatalf("expected allow-origin to echo origin, got %q", got)
	}
}

