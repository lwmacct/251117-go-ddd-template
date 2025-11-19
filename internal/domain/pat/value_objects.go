package pat

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PermissionList is a custom type for JSON array of permissions
type PermissionList []string

// Scan implements sql.Scanner interface
func (p *PermissionList) Scan(value interface{}) error {
	if value == nil {
		*p = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PermissionList")
	}

	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer interface
func (p PermissionList) Value() (driver.Value, error) {
	if p == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(p)
}

// StringList is a custom type for JSON array of strings
type StringList []string

// Scan implements sql.Scanner interface
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringList")
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(s)
}
