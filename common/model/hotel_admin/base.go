// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了通用的基础类型，主要用于处理数据库中的JSON类型字段
//
// 功能说明：
//   - 提供自定义JSON类型，支持数据库JSON字段的读写
//   - 实现GORM所需的Scanner和Valuer接口
//   - 实现标准库json的Marshaler和Unmarshaler接口
//   - 兼容MySQL的JSON类型、PostgreSQL的JSONB类型等
package hotel_admin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ==================== 自定义JSON类型 ====================

// JSON 自定义JSON类型，用于GORM处理数据库中的JSON字段
//
// 使用场景：
//   - 当数据库字段类型为JSON时，使用此类型进行映射
//   - 支持自动序列化和反序列化
//   - 保存原始JSON格式，避免重复解析
//
// 使用示例：
//
//	type Config struct {
//	    ID       uint64 `gorm:"primaryKey"`
//	    Settings *JSON  `gorm:"type:json;comment:配置信息"`  // 数据库JSON字段
//	}
type JSON json.RawMessage

// ==================== 数据库接口实现 ====================

// Scan 实现 sql.Scanner 接口，用于从数据库读取JSON数据
//
// 当GORM从数据库读取数据时会自动调用此方法，将数据库中的JSON字段转换为Go类型
//
// 支持的数据库返回类型：
//   - []byte: 字节数组（MySQL JSON类型返回）
//   - string: 字符串（部分数据库返回）
//   - nil: NULL值
//
// 参数：
//   - value: 数据库返回的原始值
//
// 返回：
//   - error: 如果类型不支持则返回错误
func (j *JSON) Scan(value interface{}) error {
	// 处理NULL值 - 将数据库NULL转换为Go的nil
	if value == nil {
		*j = nil
		return nil
	}

	// 根据不同类型进行转换
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		// MySQL JSON类型通常返回[]byte
		bytes = v
	case string:
		// 部分数据库驱动返回string
		bytes = []byte(v)
	default:
		// 不支持的类型，返回错误
		return fmt.Errorf("JSON类型转换失败：不支持的数据类型 %T", value)
	}

	// 保存原始JSON字节数组
	*j = JSON(bytes)
	return nil
}

// Value 实现 driver.Valuer 接口，用于向数据库写入JSON数据
//
// 当GORM向数据库写入数据时会自动调用此方法，将Go类型转换为数据库JSON字段
//
// 返回：
//   - driver.Value: 序列化后的JSON字节数组，如果为空则返回NULL
//   - error: 序列化失败时返回错误
func (j JSON) Value() (driver.Value, error) {
	// 空值返回NULL - 避免存储空字符串
	if len(j) == 0 {
		return nil, nil
	}

	// 序列化为JSON字节数组
	return json.RawMessage(j).MarshalJSON()
}

// ==================== JSON接口实现 ====================

// MarshalJSON 实现 json.Marshaler 接口，用于JSON序列化
//
// 当使用 json.Marshal() 序列化时会自动调用此方法
// 主要用于API响应时将模型转换为JSON格式
//
// 返回：
//   - []byte: 序列化后的JSON字节数组
//   - error: 序列化失败时返回错误
func (j JSON) MarshalJSON() ([]byte, error) {
	// 空值返回null - JSON标准null值
	if len(j) == 0 {
		return []byte("null"), nil
	}

	// 直接返回原始JSON数据 - 避免重复解析
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，用于JSON反序列化
//
// 当使用 json.Unmarshal() 反序列化时会自动调用此方法
// 主要用于接收API请求时将JSON转换为模型
//
// 参数：
//   - data: 要反序列化的JSON字节数组
//
// 返回：
//   - error: 如果指针为nil则返回错误
func (j *JSON) UnmarshalJSON(data []byte) error {
	// 检查指针是否为nil - 防止空指针异常
	if j == nil {
		return errors.New("JSON反序列化失败：不能对nil指针进行反序列化")
	}

	// 直接存储原始JSON数据 - 延迟解析，提升性能
	*j = JSON(data)
	return nil
}

// ==================== 工具方法 ====================

// String 返回JSON的字符串表示
//
// 用于打印和日志输出，方便调试
//
// 返回：
//   - string: JSON字符串
func (j JSON) String() string {
	return string(j)
}
