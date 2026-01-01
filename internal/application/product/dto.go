package product

import "time"

// ============================================================================
// Request DTOs
// ============================================================================

// CreateProductDTO 创建产品请求 DTO
type CreateProductDTO struct {
	Code        string  `json:"code" binding:"required,min=2,max=50,alphanum"`
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Type        string  `json:"type" binding:"required,oneof=personal team"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"gte=0"`
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive"`
	LayoutRef   string  `json:"layout_ref" binding:"omitempty,max=100"`
	MaxSeats    int     `json:"max_seats" binding:"gte=0"`
	TrialDays   int     `json:"trial_days" binding:"gte=0"`
}

// UpdateProductDTO 更新产品请求 DTO
type UpdateProductDTO struct {
	Code        *string  `json:"code" binding:"omitempty,min=2,max=50,alphanum"`
	Name        *string  `json:"name" binding:"omitempty,min=2,max=100"`
	Type        *string  `json:"type" binding:"omitempty,oneof=personal team"`
	Description *string  `json:"description" binding:"omitempty,max=500"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
	LayoutRef   *string  `json:"layout_ref" binding:"omitempty,max=100"`
	MaxSeats    *int     `json:"max_seats" binding:"omitempty,gte=0"`
	TrialDays   *int     `json:"trial_days" binding:"omitempty,gte=0"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// ProductDTO 产品响应 DTO
type ProductDTO struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	LayoutRef   string    `json:"layout_ref"`
	MaxSeats    int       `json:"max_seats"`
	TrialDays   int       `json:"trial_days"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductCatalogDTO 产品目录响应 DTO（公开信息）
type ProductCatalogDTO struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	LayoutRef   string  `json:"layout_ref"`
	MaxSeats    int     `json:"max_seats"`
	TrialDays   int     `json:"trial_days"`
}
