package product

import "errors"

// 产品相关错误
var (
	// ErrProductNotFound 产品不存在
	ErrProductNotFound = errors.New("product not found")
	// ErrProductNameExists 产品名称已存在
	ErrProductNameExists = errors.New("product name already exists")
	// ErrProductCodeExists 产品代码已存在
	ErrProductCodeExists = errors.New("product code already exists")
	// ErrInvalidProductCode 无效的产品代码
	ErrInvalidProductCode = errors.New("invalid product code")
	// ErrInvalidProductType 无效的产品类型
	ErrInvalidProductType = errors.New("invalid product type")
)

// 产品业务约束错误
var (
	// ErrCannotDeleteProduct 产品有订阅时不能删除
	ErrCannotDeleteProduct = errors.New("cannot delete product with active subscriptions")
)
