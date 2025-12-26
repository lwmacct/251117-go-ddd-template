package response

// PaginationQueryDTO 通用分页查询参数
// 所有列表查询 DTO 应嵌入此结构体以获得统一的分页能力
type PaginationQueryDTO struct {
	// Page 页码，从 1 开始
	Page int `form:"page" json:"page" binding:"omitempty,min=1" minimum:"1" default:"1"`
	// Limit 每页数量，默认 20，最大 1000
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=1000" minimum:"1" maximum:"1000" default:"20"`
}

// GetPage 获取页码，确保最小值为 1
func (p *PaginationQueryDTO) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetLimit 获取每页数量，确保在有效范围内
func (p *PaginationQueryDTO) GetLimit() int {
	if p.Limit < 1 {
		return 20
	}
	if p.Limit > 1000 {
		return 1000
	}
	return p.Limit
}

// GetOffset 计算数据库查询偏移量
func (p *PaginationQueryDTO) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
