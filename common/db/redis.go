/*
 * Redis缓存连接和初始化模块
 * 
 * 功能说明：
 * - 建立Redis连接
 * - 测试连接有效性
 * 
 * Redis用途：
 * - 缓存热点数据（车次信息、余票数量等）
 * - 分布式锁（防止并发问题）
 * - 会话存储（Token、用户状态等）
 * - 限流计数器
 */
package db

import (
	"context"
	"example_shop/common/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Rdb 全局Redis客户端对象
// 初始化后供所有模块使用
var Rdb *redis.Client

// Ctx 全局上下文对象
// 用于Redis操作，后续可扩展超时、取消等控制
var Ctx = context.Background()

// RedisInit 初始化Redis连接
// 
// 执行流程：
// 1. 读取Redis配置
// 2. 创建Redis客户端
// 3. 测试连接有效性（Ping）
// 
// 返回值：
// - error: 如果初始化失败，返回错误信息
func RedisInit() error {
	redisCfg := config.Cfg.Redis
	if redisCfg.Port == 0 {
		redisCfg.Port = 6379
	}

	// 创建Redis客户端
	// NewClient：创建标准的Redis客户端（单机模式）
	// 如需集群模式，可使用NewClusterClient
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port), // Redis地址
		Password: redisCfg.Password,                                  // Redis密码（如果未设置则为空）
		DB:       redisCfg.Database,                                  // Redis数据库编号（0-15，默认使用0）
	})

	// 测试连接有效性
	// Ping：发送PING命令测试Redis是否可用
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	
	fmt.Println("Redis连接成功")
	return nil
}
