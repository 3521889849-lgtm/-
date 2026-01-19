// Package config 提供系统配置管理
//
// 本文件定义了系统的所有配置结构，包括：
//   - MySQL数据库配置
//   - Redis缓存配置
//   - AES加密密钥配置
//
// 配置从 conf/config.yaml 文件中加载
package config

// Cfg 全局配置实例
// 在程序启动时由viperInit.go初始化
var Cfg = new(Config)

// Config 系统总配置结构
//
// 包含系统运行所需的所有配置项
// 通过viper从config.yaml文件加载
type Config struct {
	MysqlInit // MySQL数据库配置
	RedisInit // Redis缓存配置
	AesConfig // AES加密配置
}

// MysqlInit MySQL数据库配置
//
// 用于连接MySQL数据库
type MysqlInit struct {
	Host     string // MySQL服务器地址，如："localhost"、"127.0.0.1"
	Port     int    // MySQL端口号，默认：3306
	User     string // 数据库用户名
	Password string // 数据库密码
	Database string // 数据库名称
}

// RedisInit Redis缓存配置
//
// 用于连接Redis服务器
type RedisInit struct {
	Host     string // Redis服务器地址
	Port     int    // Redis端口号，默认：6379
	Password string // Redis密码，无密码则为空字符串
	Database int    // Redis数据库编号，默认：0
}

// AesConfig AES加密配置
//
// 用于敏感数据加密（身份证号、手机号等）
type AesConfig struct {
	AesKey string // AES加密密钥，必须是16、24或32字节
}
