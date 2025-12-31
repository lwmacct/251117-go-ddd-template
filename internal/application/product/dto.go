package product

import "time"

// ============================================================================
// Request DTOs
// ============================================================================

// CreateProductDTO 创建产品请求 DTO
type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"gte=0"`
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateProductDTO 更新产品请求 DTO
type UpdateProductDTO struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string  `json:"description" binding:"omitempty,max=500"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// ProductDTO 产品响应 DTO
type ProductDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
