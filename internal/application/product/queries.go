package product

// GetProductQuery 获取产品查询
type GetProductQuery struct {
	ID uint
}

// ListProductsQuery 产品列表查询
type ListProductsQuery struct {
	Offset int
	Limit  int
}
