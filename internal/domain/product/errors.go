package product

import "errors"

var (
	// ErrProductNotFound 产品不存在。
	ErrProductNotFound = errors.New("产品不存在")
	// ErrProductNameExists 产品名称已存在。
	ErrProductNameExists = errors.New("产品名称已存在")
)
