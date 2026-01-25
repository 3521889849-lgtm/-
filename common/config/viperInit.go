/*
 * Viper配置加载模块
 *
 * 功能说明：
 * - 使用Viper库加载YAML配置文件
 * - 支持多路径搜索配置文件
 * - 将配置解析到Config结构体
 *
 * Viper优势：
 * - 支持多种配置格式（YAML、JSON、TOML等）
 * - 支持环境变量覆盖
 * - 支持配置文件热重载
 */
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ViperInit 初始化Viper配置加载器
//
// 功能：
// 1. 设置配置文件搜索路径（支持不同运行目录）
// 2. 设置配置文件名和类型
// 3. 读取配置文件
// 4. 解析配置文件到Cfg全局对象
//
// 返回值：
// - error: 如果配置加载失败，返回错误信息
func ViperInit() error {
	var err error

	// 添加多个配置路径，支持不同运行目录运行程序
	// 例如：从项目根目录运行、从cmd/gateway目录运行等
	viper.AddConfigPath("conf")          // 当前目录下的conf文件夹
	viper.AddConfigPath("../conf")       // 上一级目录的conf文件夹
	viper.AddConfigPath("../../conf")    // 上两级目录的conf文件夹
	viper.AddConfigPath("../../../conf") // 上三级目录的conf文件夹（项目根目录）

	// 设置配置文件名（不含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件类型为YAML
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("未找到配置文件 conf/config.yaml（可参考 conf/config.yaml 并复制为 config.yaml），也可使用环境变量覆盖：%w", err)
		}
		return fmt.Errorf("配置文件读取失败: %w", err)
	}

	// 将配置文件内容解析到Cfg全局对象
	// Unmarshal：将Viper读取的配置映射到Config结构体
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return fmt.Errorf("配置文件解析失败: %w", err)
	}
	if err := ValidateAndNormalize(); err != nil {
		return fmt.Errorf("配置校验失败: %w", err)
	}
	// 配置加载成功提示
	fmt.Println("配置动态加载成功:", viper.ConfigFileUsed())
	return nil
}
