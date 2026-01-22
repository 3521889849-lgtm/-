package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JSON 自定义 JSON 类型，兼容 MySQL JSON 字段与 GORM
type JSON json.RawMessage

// 实现 sql.Scanner 接口，用于数据库读取时解析 JSON
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("{}")
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("unsupported type: %T, expected []byte", value))
	}

	// 验证 JSON 格式合法性
	if !json.Valid(bytes) {
		return errors.New(fmt.Sprintf("invalid JSON: %s", string(bytes)))
	}

	*j = JSON(bytes)
	return nil
}

// 实现 driver.Valuer 接口，用于数据库写入时序列化 JSON
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return string(j), nil
}

// 实现 gorm.SchemaInterface 接口，让 GORM 识别字段类型
func (j JSON) GormDataType() string {
	return "json"
}

// 实现 gorm.SchemaInterface 接口，自定义字段数据库类型
func (j JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "mariadb":
		return "JSON" // MySQL 数据库中字段类型为 JSON
	case "postgres":
		return "JSONB" // 兼容 PostgreSQL（可选）
	default:
		return "TEXT" // 其他数据库降级为 TEXT
	}
}

// 辅助方法：将结构体序列化为 JSON 类型（业务中常用）
func ToJSON(v interface{}) (JSON, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return JSON("{}"), err
	}
	return JSON(bytes), nil
}

// 辅助方法：将 JSON 类型反序列化为结构体（业务中常用）
func (j JSON) ToStruct(v interface{}) error {
	return json.Unmarshal(j, v)
}
