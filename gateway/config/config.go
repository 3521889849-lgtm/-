// Package config 提供 Gateway 网关服务的配置管理功能
// 主要功能：
// 1. 加载和解析YAML配置文件
// 2. 管理服务地址、JWT等配置项
// 3. 提供全局配置访问入口
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper" // 配置管理库
)

// Config 网关服务的根配置结构
type Config struct {
	Server   ServerConfig             `mapstructure:"server"`   // 服务器配置
	Services map[string]ServiceConfig `mapstructure:"services"` // 后端服务配置映射
	JWT      JWTConfig                `mapstructure:"jwt"`      // JWT认证配置
	Trace    TraceConfig              `mapstructure:"trace"`    // 链路追踪配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address string `mapstructure:"address"` // 监听地址，如 ":8080"
}

// ServiceConfig 后端服务配置项
type ServiceConfig struct {
	Name    string `mapstructure:"name"`    // 服务名称
	Address string `mapstructure:"address"` // 服务地址（如 "127.0.0.1:9999"）
}

// JWTConfig JWT配置结构
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`       // JWT签名密钥，生产环境应使用强密码
	ExpireHours int    `mapstructure:"expire_hours"` // Token过期时间（小时）
}

// TraceConfig 链路追踪配置
// 支持 Jaeger 导出，可通过 enabled 字段控制开关
type TraceConfig struct {
	Enabled       bool          `mapstructure:"enabled"`        // 是否启用链路追踪
	ServiceName   string        `mapstructure:"service_name"`   // 服务名称（显示在 Jaeger UI）
	JaegerEnabled bool          `mapstructure:"jaeger_enabled"` // 是否启用 Jaeger 导出
	Endpoint      string        `mapstructure:"endpoint"`       // Jaeger Collector 地址
	SampleRate    float64       `mapstructure:"sample_rate"`    // 采样率 (0.0~1.0)
	BatchSize     int           `mapstructure:"batch_size"`     // 批量发送大小
	FlushInterval time.Duration `mapstructure:"flush_interval"` // 刷新间隔
}

// GlobalConfig 全局配置实例，由InitConfig初始化
var GlobalConfig *Config

// InitConfig 初始化配置
// 从指定YAML文件加载配置，并设置默认值
// 参数:
//   - configPath: 配置文件路径
//
// 返回: 配置加载错误
func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置JWT配置默认值
	viper.SetDefault("jwt.secret", "piaowu-secret-key-2026")
	viper.SetDefault("jwt.expire_hours", 24)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// GetJWTSecret 获取JWT签名密钥
// 优先从配置文件读取，如果未配置则返回默认值
// 返回: JWT密钥字符串
func GetJWTSecret() string {
	if GlobalConfig != nil && GlobalConfig.JWT.Secret != "" {
		return GlobalConfig.JWT.Secret
	}
	return "piaowu-secret-key-2026" // 默认密钥
}

// GetJWTExpireHours 获取Token过期时间
// 优先从配置文件读取，如果未配置则返回默认24小时
// 返回: 过期时间（小时数）
func GetJWTExpireHours() int {
	if GlobalConfig != nil && GlobalConfig.JWT.ExpireHours > 0 {
		return GlobalConfig.JWT.ExpireHours
	}
	return 24 // 默认24小时
}
