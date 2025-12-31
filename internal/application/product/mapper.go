package product

import "github.com/lwmacct/251117-go-ddd-template/internal/domain/product"

// ToProductDTO 将产品实体转换为 DTO
func ToProductDTO(p *product.Product) *ProductDTO {
	if p == nil {
		return nil
	}
	return &ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      string(p.Status),
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
