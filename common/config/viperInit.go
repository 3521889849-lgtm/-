// Package config 提供系统配置管理
//
// 本文件实现配置文件的加载和解析
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ViperInit 初始化配置文件加载
//
// 功能：
//   - 使用viper库加载YAML配置文件
//   - 支持多个配置路径，适配不同的运行目录
//   - 自动解析配置到Cfg全局变量
//   - 支持配置文件热更新（可选）
//
// 配置文件位置：
//   - 优先查找 conf/config.yaml
//   - 然后查找 ../conf/config.yaml
//   - 然后查找 ../../conf/config.yaml
//   - 最后查找 ../../../conf/config.yaml
//
// 使用场景：
//   - API服务启动时：api/main/main.go 调用
//   - RPC服务启动时：rpc/hotel/main/main.go 调用
//   - 根据不同的工作目录，自动找到配置文件
//
// 返回：
//   - error: 配置加载失败时返回错误信息
func ViperInit() error {
	// ========== 配置文件路径 ==========
	
	// 添加多个配置路径，支持不同的运行目录
	// viper会按顺序尝试这些路径，找到第一个存在的配置文件
	viper.AddConfigPath("conf")             // 当前目录下的conf文件夹
	viper.AddConfigPath("../conf")          // 上一级目录的conf文件夹
	viper.AddConfigPath("../../conf")       // 上两级目录的conf文件夹
	viper.AddConfigPath("../../../conf")    // 上三级目录的conf文件夹（对应项目根目录）
	
	// ========== 配置文件名称和类型 ==========
	
	// 设置配置文件名（不含扩展名）
	viper.SetConfigName("config")
	
	// 设置配置文件类型
	// 支持的类型：yaml, json, toml, ini等
	viper.SetConfigType("yaml")
	
	// ========== 读取配置文件 ==========
	
	// 读取配置文件
	// 如果文件不存在或格式错误，会返回错误
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("配置文件读取失败: %w", err)
	}
	
	// ========== 解析配置到结构体 ==========
	
	// 将配置文件内容解析到Cfg全局变量
	// Unmarshal会自动根据结构体字段名匹配配置项
	// 支持嵌套结构体、切片、Map等复杂类型
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return fmt.Errorf("配置文件解析失败: %w", err)
	}
	
	return nil
}

// 使用示例：
//
//	func main() {
//	    // 加载配置
//	    if err := config.ViperInit(); err != nil {
//	        log.Fatal("配置加载失败:", err)
//	    }
//	    
//	    // 使用配置
//	    fmt.Printf("MySQL地址: %s:%d\n", config.Cfg.Host, config.Cfg.Port)
//	}
