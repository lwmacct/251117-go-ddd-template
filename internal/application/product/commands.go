package product

// CreateProductCommand 创建产品命令
type CreateProductCommand struct {
	Name        string
	Description string
	Price       float64
	Status      string // active, inactive
}

// UpdateProductCommand 更新产品命令
type UpdateProductCommand struct {
	ID          uint
	Name        *string
	Description *string
	Price       *float64
	Status      *string
}

// DeleteProductCommand 删除产品命令
type DeleteProductCommand struct {
	ID uint
}
