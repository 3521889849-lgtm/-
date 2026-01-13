package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON 自定义JSON类型，用于GORM处理JSON字段
type JSON json.RawMessage

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to unmarshal JSON value")
	}
	*j = JSON(bytes)
	return nil
}

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// MarshalJSON 实现 json.Marshaler 接口
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json: UnmarshalJSON on nil pointer")
	}
	*j = JSON(data)
	return nil
}

// String 返回JSON字符串表示
func (j JSON) String() string {
	return string(j)
}
