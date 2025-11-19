package twofa

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// RecoveryCodes 恢复码数组类型（值对象）
type RecoveryCodes []string

// Scan 实现 sql.Scanner 接口，从数据库读取时自动处理空值
func (r *RecoveryCodes) Scan(value interface{}) error {
	if value == nil {
		*r = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to unmarshal RecoveryCodes value")
	}

	// 如果是空JSON，使用空数组
	if len(bytes) == 0 || string(bytes) == "{}" || string(bytes) == "[]" {
		*r = []string{}
		return nil
	}

	return json.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口，写入数据库
func (r RecoveryCodes) Value() (driver.Value, error) {
	if len(r) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(r)
}
