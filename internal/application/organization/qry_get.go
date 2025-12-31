package organization

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
)

// GetHandler 获取组织查询处理器
type GetHandler struct {
	orgQueryRepo organization.QueryRepository
}

// NewGetHandler 创建获取组织查询处理器
func NewGetHandler(orgQueryRepo organization.QueryRepository) *GetHandler {
	return &GetHandler{orgQueryRepo: orgQueryRepo}
}

// Handle 处理获取组织查询
func (h *GetHandler) Handle(ctx context.Context, query GetOrgQuery) (*OrgDTO, error) {
	org, err := h.orgQueryRepo.GetByID(ctx, query.OrgID)
	if err != nil {
		return nil, err
	}
	return ToOrgDTO(org), nil
}
