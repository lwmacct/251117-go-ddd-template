package product

import "github.com/lwmacct/251117-go-ddd-template/internal/domain/product"

// ToProductDTO 将产品实体转换为 DTO
func ToProductDTO(p *product.Product) *ProductDTO {
	if p == nil {
		return nil
	}
	return &ProductDTO{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Type:        string(p.Type),
		Description: p.Description,
		Price:       p.Price,
		Status:      string(p.Status),
		LayoutRef:   p.LayoutRef,
		MaxSeats:    p.MaxSeats,
		TrialDays:   p.TrialDays,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToProductDTOs 将产品实体列表转换为 DTO 列表
func ToProductDTOs(products []*product.Product) []*ProductDTO {
	if len(products) == 0 {
		return nil
	}
	dtos := make([]*ProductDTO, 0, len(products))
	for _, p := range products {
		if dto := ToProductDTO(p); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToProductCatalogDTOs 将产品实体列表转换为目录 DTO 列表（公开信息）
func ToProductCatalogDTOs(products []*product.Product) []*ProductCatalogDTO {
	if len(products) == 0 {
		return nil
	}
	dtos := make([]*ProductCatalogDTO, 0, len(products))
	for _, p := range products {
		if dto := ToProductCatalogDTO(p); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToProductCatalogDTO 将产品实体转换为目录 DTO
func ToProductCatalogDTO(p *product.Product) *ProductCatalogDTO {
	if p == nil {
		return nil
	}
	return &ProductCatalogDTO{
		Code:        p.Code,
		Name:        p.Name,
		Type:        string(p.Type),
		Description: p.Description,
		Price:       p.Price,
		LayoutRef:   p.LayoutRef,
		MaxSeats:    p.MaxSeats,
		TrialDays:   p.TrialDays,
	}
}
