package product

// CreateProductCommand 创建产品命令
type CreateProductCommand struct {
	Code        string
	Name        string
	Type        string // personal, team
	Description string
	Price       float64
	Status      string // active, inactive
	LayoutRef   string
	MaxSeats    int
	TrialDays   int
}

// UpdateProductCommand 更新产品命令
type UpdateProductCommand struct {
	ID          uint
	Code        *string
	Name        *string
	Type        *string
	Description *string
	Price       *float64
	Status      *string
	LayoutRef   *string
	MaxSeats    *int
	TrialDays   *int
}

// DeleteProductCommand 删除产品命令
type DeleteProductCommand struct {
	ID uint
}
