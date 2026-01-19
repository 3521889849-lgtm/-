package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example_shop/api/handler"

	"github.com/gin-gonic/gin"
)

func TestCreateMember_BadJSON_ReturnsCodeMsg(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewMemberHandler()
	r.POST("/api/v1/members", h.CreateMember)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/members", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if _, ok := body["error"]; ok {
		t.Fatalf("unexpected error field in response: %v", body)
	}
	if code, ok := body["code"].(float64); !ok || int(code) != 400 {
		t.Fatalf("expected code=400, got %v", body["code"])
	}
	if msg, ok := body["msg"].(string); !ok || msg == "" {
		t.Fatalf("expected msg string, got %v", body["msg"])
	}
}

func TestGetMemberPointsBalance_InvalidID_ReturnsCodeMsg(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewMemberPointsHandler()
	r.GET("/api/v1/members/:id/points-balance", h.GetMemberPointsBalance)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members/abc/points-balance", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if _, ok := body["error"]; ok {
		t.Fatalf("unexpected error field in response: %v", body)
	}
	if code, ok := body["code"].(float64); !ok || int(code) != 400 {
		t.Fatalf("expected code=400, got %v", body["code"])
	}
	if msg, ok := body["msg"].(string); !ok || msg == "" {
		t.Fatalf("expected msg string, got %v", body["msg"])
	}
}

