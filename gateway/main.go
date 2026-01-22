// Package main 是 Gateway 网关服务的入口
// 该服务作为系统的API网关，负责：
// 1. 接收前端HTTP请求并转发到后端RPC服务
// 2. 处理WebSocket实时通信
// 3. JWT身份认证和权限校验
// 4. 跨域(CORS)处理
// 5. 请求链路追踪(TraceID)
package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"example_shop/gateway/config"     // 网关配置管理
	"example_shop/gateway/handler"    // HTTP请求处理器
	"example_shop/gateway/middleware" // 中间件（认证、链路追踪等）
	"example_shop/gateway/router"     // 路由配置
	"example_shop/gateway/rpc"        // RPC客户端
	"example_shop/gateway/ws"         // WebSocket模块
	"example_shop/pkg/logger"         // 日志工具
	"example_shop/pkg/trace"          // 链路追踪

	"go.uber.org/zap"
)

func main() {
	// 初始化日志系统（输出到控制台和文件）
	if err := logger.InitLoggerWithFile("gateway", "./logs"); err != nil {
		panic("Failed to init logger: " + err.Error())
	}
	defer logger.Sync()

	// 初始化配置
	if err := config.InitConfig(mustResolveConfigPath(
		"gateway/config/config.yaml",
		"config/config.yaml",
	)); err != nil {
		logger.Fatal("Failed to init config", zap.Error(err))
	}

	// 同步 JWT 配置到 middleware（确保登录和WebSocket验证使用同一密钥）
	middleware.SetJWTSecret(config.GetJWTSecret())
	middleware.SetJWTExpireTime(config.GetJWTExpireHours())

	// 初始化链路追踪（Jaeger）
	initTracing()
	defer trace.ShutdownJaeger()

	// 初始化 RPC 客户端
	customerService := config.GlobalConfig.Services["customer"]
	customerClient, err := rpc.NewCustomerClient(customerService.Name, customerService.Address)
	if err != nil {
		logger.Fatal("Failed to create customer client", zap.Error(err))
	}
	defer customerClient.Close()

	// 初始化 Handler
	customerHandler := handler.NewCustomerHandler(customerClient)

	// 初始化 WebSocket Hub
	hub := ws.NewHub(customerClient)
	go hub.Run()

	// 设置路由
	mux := router.SetupRoutes(customerHandler, hub)

	// 应用中间件：TraceID -> CORS
	httpHandler := middleware.TraceMiddleware(mux)
	httpHandler = withCORS(httpHandler)

	// 启动 HTTP 服务器
	logger.Info("Gateway started", zap.String("address", config.GlobalConfig.Server.Address))
	if err := http.ListenAndServe(config.GlobalConfig.Server.Address, httpHandler); err != nil {
		logger.Fatal("Gateway stopped with error", zap.Error(err))
	}
}

// withCORS 跨域处理中间件
// 功能说明：
// 1. 检查请求来源是否在允许列表中（localhost/127.0.0.1）
// 2. 设置CORS响应头，允许跨域请求
// 3. 处理OPTIONS预检请求
// 参数:
//   - next: 下一个HTTP处理器
//
// 返回: 包装后的HTTP处理器
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := isAllowedOrigin(origin)

		// 设置CORS响应头
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		}

		// 处理OPTIONS预检请求
		if r.Method == http.MethodOptions {
			if allowed {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				reqHeaders := r.Header.Get("Access-Control-Request-Headers")
				if strings.TrimSpace(reqHeaders) == "" {
					reqHeaders = "Content-Type, Authorization"
				}
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				w.Header().Set("Access-Control-Max-Age", "600") // 预检结果缓存10分钟
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isAllowedOrigin 检查请求来源是否在允许的白名单中
// 目前仅允许本地开发环境（localhost/127.0.0.1）
// 生产环境应根据实际域名配置
// 参数:
//   - origin: HTTP请求的Origin头
//
// 返回: true=允许跨域, false=拒绝
func isAllowedOrigin(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1"
}

// mustResolveConfigPath 解析配置文件路径
// 按顺序尝试多个候选路径，返回第一个存在的文件路径
// 支持从当前目录或项目根目录查找
// 参数:
//   - candidates: 候选配置文件路径列表
//
// 返回: 找到的配置文件路径，如果都不存在则返回第一个候选路径
func mustResolveConfigPath(candidates ...string) string {
	// 首先尝试候选路径
	for _, p := range candidates {
		if fileExists(p) {
			return p
		}
	}

	// 如果候选路径都不存在，尝试从项目根目录查找
	root, ok := findProjectRoot()
	if ok {
		for _, p := range candidates {
			if fileExists(filepath.Join(root, p)) {
				return filepath.Join(root, p)
			}
		}
	}

	return candidates[0]
}

// fileExists 检查文件是否存在
// 参数:
//   - path: 文件路径
//
// 返回: true=文件存在, false=文件不存在或路径为空
func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// findProjectRoot 查找项目根目录
// 从当前工作目录向上遍历，找到包含go.mod的目录即为项目根目录
// 返回:
//   - string: 项目根目录路径
//   - bool: true=找到, false=未找到
func findProjectRoot() (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	dir := wd
	for {
		// 找到go.mod文件，说明是项目根目录
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir, true
		}
		// 向上遍历
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

// initTracing 初始化链路追踪
func initTracing() {
	cfg := config.GlobalConfig.Trace
	if !cfg.Enabled {
		logger.Info("Tracing is disabled")
		return
	}

	// 配置 Jaeger Exporter
	jaegerCfg := &trace.JaegerConfig{
		Enabled:       cfg.JaegerEnabled,
		ServiceName:   cfg.ServiceName,
		Endpoint:      cfg.Endpoint,
		SampleRate:    cfg.SampleRate,
		BatchSize:     cfg.BatchSize,
		FlushInterval: cfg.FlushInterval,
	}

	// 设置默认值
	if jaegerCfg.ServiceName == "" {
		jaegerCfg.ServiceName = "gateway"
	}
	if jaegerCfg.Endpoint == "" {
		jaegerCfg.Endpoint = "http://localhost:4318/v1/traces"
	}
	if jaegerCfg.SampleRate <= 0 {
		jaegerCfg.SampleRate = 1.0 // Gateway 默认全量采样
	}

	trace.InitJaeger(jaegerCfg)
	logger.Info("Tracing initialized",
		zap.String("service", jaegerCfg.ServiceName),
		zap.Bool("jaeger_enabled", jaegerCfg.Enabled),
		zap.Float64("sample_rate", jaegerCfg.SampleRate),
	)
}
