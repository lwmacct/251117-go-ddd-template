package testutil

import "time"

// PtrUint 创建 uint 指针
func PtrUint(v uint) *uint {
	return &v
}

// PtrString 创建 string 指针
func PtrString(v string) *string {
	return &v
}

// PtrTime 创建 time.Time 指针
func PtrTime(v time.Time) *time.Time {
	return &v
}

// PtrInt 创建 int 指针
func PtrInt(v int) *int {
	return &v
}

// PtrBool 创建 bool 指针
func PtrBool(v bool) *bool {
	return &v
}
