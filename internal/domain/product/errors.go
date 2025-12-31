package product

import "errors"

var (
	// ErrProductNotFound 产品不存在。
	ErrProductNotFound = errors.New("product not found")
	// ErrProductNameExists 产品名称已存在。
	ErrProductNameExists = errors.New("product name already exists")
)
