// Package init 提供系统初始化功能
//
// 本文件在程序启动时自动执行初始化流程
package init

import (
	"example_shop/common/config"
	"example_shop/common/db"
	"log"
)

// init 包级初始化函数
//
// 功能：
//   - Go语言自动调用，在main函数之前执行
//   - 按顺序初始化系统各个组件
//   - 任何一步失败都会终止程序
//
// 初始化顺序：
//   1. 配置文件加载（viperInit）
//   2. MySQL数据库连接
//   3. Redis缓存连接
//
// 使用方式：
//   在需要初始化的包中导入此包即可：
//   import _ "example_shop/common/init"
//
// 注意事项：
//   - 此函数只会执行一次
//   - 初始化失败会直接终止程序（log.Fatalf）
//   - 建议在API和RPC服务的main包中导入
func init() {
	// ========== 第一步：加载配置文件 ==========
	// 必须最先执行，因为后续步骤依赖配置
	if err := config.ViperInit(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// ========== 第二步：初始化MySQL数据库 ==========
	// 建立数据库连接并自动迁移表结构
	if err := db.MysqlInit(); err != nil {
		log.Fatalf("MySQL连接初始化失败: %v", err)
	}

	// ========== 第三步：初始化Redis缓存 ==========
	// 建立Redis连接用于缓存和会话管理
	if err := db.RedisInit(); err != nil {
		log.Fatalf("Redis连接初始化失败: %v", err)
	}

	// ========== 初始化完成 ==========
	// 如果程序能执行到这里，说明所有初始化都成功了
	// 后续可以添加其他初始化步骤，如：
	// - 初始化日志系统
	// - 加载定时任务
	// - 初始化消息队列
	// - 加载缓存数据
}
