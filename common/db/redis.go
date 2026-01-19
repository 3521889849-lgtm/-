// Package db 提供数据库连接管理
//
// 本文件实现Redis缓存的初始化和连接管理
package db

import (
	"context"
	"example_shop/common/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Rdb 全局Redis客户端实例
//
// 程序启动时初始化，全局共享使用
// 用于缓存、会话管理、分布式锁等场景
var Rdb *redis.Client

// Ctx 全局Context实例
//
// 用于Redis操作的上下文
// 可以用于控制超时、取消等
var Ctx = context.Background()

// RedisInit 初始化Redis连接
//
// 功能：
//   - 读取配置创建Redis客户端
//   - 测试连接是否正常（Ping）
//   - 初始化全局Context
//
// 返回：
//   - error: 初始化失败时返回错误信息
func RedisInit() error {
	// 读取Redis配置
	redisCfg := config.Cfg.RedisInit

	// 创建Redis客户端
	// Options配置说明：
	//   - Addr: Redis服务器地址（格式：host:port）
	//   - Password: Redis密码，无密码则为空字符串
	//   - DB: Redis数据库编号（0-15），默认0
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.Database,
	})

	// 测试Redis连接
	// Ping命令：检查Redis服务是否正常
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}
