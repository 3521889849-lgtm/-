// Package model 提供数据模型的基础类型定义
//
// 本文件定义了通用的数据类型，主要用于数据库字段映射和JSON序列化
package model

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
//   - 兼容MySQL的JSON类型、PostgreSQL的JSONB类型等
//
// 使用示例：
//
//	type Config struct {
//	    ID       uint64 `gorm:"primaryKey"`
//	    Settings *JSON  `gorm:"type:json"`  // 数据库JSON字段
//	}
type JSON json.RawMessage

// ==================== 数据库接口实现 ====================

// Scan 实现 sql.Scanner 接口，用于从数据库读取JSON数据
//
// 当GORM从数据库读取数据时会自动调用此方法
// 支持的数据库返回类型：[]byte、string、nil
//
// 参数：
//   - value: 数据库返回的原始值
//
// 返回：
//   - error: 如果类型不支持则返回错误
func (j *JSON) Scan(value interface{}) error {
	// 处理NULL值
	if value == nil {
		*j = nil
		return nil
	}

	// 根据不同类型进行转换
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("JSON类型转换失败：不支持的数据类型 %T", value)
	}

	*j = JSON(bytes)
	return nil
}

// Value 实现 driver.Valuer 接口，用于向数据库写入JSON数据
//
// 当GORM向数据库写入数据时会自动调用此方法
//
// 返回：
//   - driver.Value: 序列化后的JSON字节数组
//   - error: 序列化失败时返回错误
func (j JSON) Value() (driver.Value, error) {
	// 空值返回NULL
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
//
// 返回：
//   - []byte: 序列化后的JSON字节数组
//   - error: 序列化失败时返回错误
func (j JSON) MarshalJSON() ([]byte, error) {
	// 空值返回null
	if len(j) == 0 {
		return []byte("null"), nil
	}

	// 直接返回原始JSON数据
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，用于JSON反序列化
//
// 当使用 json.Unmarshal() 反序列化时会自动调用此方法
//
// 参数：
//   - data: 要反序列化的JSON字节数组
//
// 返回：
//   - error: 如果指针为nil则返回错误
func (j *JSON) UnmarshalJSON(data []byte) error {
	// 检查指针是否为nil
	if j == nil {
		return errors.New("JSON反序列化失败：不能对nil指针进行反序列化")
	}

	// 直接存储原始JSON数据
	*j = JSON(data)
	return nil
}

// ==================== 工具方法 ====================

// String 返回JSON的字符串表示
//
// 用于打印和日志输出
//
// 返回：
//   - string: JSON字符串
func (j JSON) String() string {
	return string(j)
}
