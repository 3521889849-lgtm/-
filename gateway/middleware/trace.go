// Package middleware 提供 Gateway 层 HTTP 中间件
// 链路追踪中间件：集成 pkg/trace 包，支持 Span 创建和上下文传播
package middleware

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"example_shop/pkg/logger"
	"example_shop/pkg/trace"

	"go.uber.org/zap"
)

// ============ 链路追踪常量 ============

const (
	// TraceIDHeader HTTP Header 中 TraceID 的键名
	// 用于在服务间传递链路追踪ID，便于日志关联和问题排查
	TraceIDHeader = "X-Trace-ID"
	// SpanIDHeader HTTP Header 中 SpanID 的键名
	SpanIDHeader = "X-Span-ID"
)

// TraceMiddleware 创建链路追踪中间件
// 功能：
// 1. 从 HTTP Header 中提取 TraceID/SpanID，如果没有则生成新的
// 2. 创建 Server Span，记录 HTTP 请求信息
// 3. 将追踪信息注入到 context 中
// 4. 在响应 Header 中返回 TraceID/SpanID
// 5. 记录请求日志（包含追踪信息）
func TraceMiddleware(next http.Handler) http.Handler {
	return TraceMiddlewareWithService("gateway", next)
}

// TraceMiddlewareWithService 创建带服务名的链路追踪中间件
func TraceMiddlewareWithService(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 1. 从 Header 提取追踪信息
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = trace.NewTraceID()
		}
		parentSpanID := r.Header.Get(SpanIDHeader)

		// 2. 创建 Server Span
		ctx := r.Context()
		ctx = trace.WithTraceID(ctx, traceID)
		if parentSpanID != "" {
			ctx = trace.WithParentSpanID(ctx, parentSpanID)
		}

		spanName := r.Method + " " + r.URL.Path
		ctx, span := trace.StartSpan(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithService(serviceName),
		)

		// 设置 HTTP 属性
		span.SetAttribute("http.method", r.Method)
		span.SetAttribute("http.path", r.URL.Path)
		span.SetAttribute("http.host", r.Host)
		span.SetAttribute("http.remote_addr", r.RemoteAddr)
		span.SetAttribute("http.user_agent", r.UserAgent())

		// 3. 同时注入到 logger 的 context 中（兼容旧代码）
		ctx = logger.WithTraceID(ctx, traceID)
		r = r.WithContext(ctx)

		// 4. 在响应 Header 中返回追踪信息
		w.Header().Set(TraceIDHeader, traceID)
		w.Header().Set(SpanIDHeader, span.SpanID)

		// 5. 创建响应记录器以捕获状态码
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 调用下一个处理器
		next.ServeHTTP(rec, r)

		// 6. 记录响应信息
		duration := time.Since(startTime)
		span.SetAttribute("http.status_code", strconv.Itoa(rec.statusCode))
		span.SetAttribute("http.duration_ms", strconv.FormatInt(duration.Milliseconds(), 10))
		span.SetAttribute("http.response_size", strconv.Itoa(rec.size))

		// 根据状态码设置 span 状态
		if rec.statusCode >= 400 {
			span.SetStatus(trace.SpanStatusError)
		}

		// 7. 结束 Span
		span.End()

		// 8. 记录请求日志（兼容旧代码）
		logger.InfoWithTrace(ctx, "HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", rec.statusCode),
			zap.Duration("duration", duration),
		)
	})
}

// responseRecorder 用于记录响应状态码和大小
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader 捕获状态码
func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Write 记录响应大小
func (rec *responseRecorder) Write(b []byte) (int, error) {
	size, err := rec.ResponseWriter.Write(b)
	rec.size += size
	return size, err
}

// Hijack 实现 http.Hijacker 接口，支持 WebSocket 升级
// WebSocket 需要接管底层 TCP 连接，必须实现此接口
func (rec *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rec.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// GetTraceIDFromRequest 从 HTTP 请求中获取 TraceID
func GetTraceIDFromRequest(r *http.Request) string {
	return trace.GetTraceID(r.Context())
}

// GetTraceIDFromContext 从 context 中获取 TraceID
func GetTraceIDFromContext(ctx context.Context) string {
	return trace.GetTraceID(ctx)
}

// GetSpanIDFromContext 从 context 中获取 SpanID
func GetSpanIDFromContext(ctx context.Context) string {
	return trace.GetSpanID(ctx)
}
