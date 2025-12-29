package auditlog

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
)

// ToAuditLogDTO 将领域实体转换为 DTO
func ToAuditLogDTO(log *auditlog.AuditLog) *AuditLogDTO {
	if log == nil {
		return nil
	}

	return &AuditLogDTO{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		Resource:  log.Resource,
		Details:   log.Details,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		Status:    log.Status,
		CreatedAt: log.CreatedAt,
	}
}

// ToAuditActionsResponseDTO 将 operation 包的审计操作定义转换为 DTO
func ToAuditActionsResponseDTO() AuditActionsResponseDTO {
	return AuditActionsResponseDTO{
		Actions:    ToAuditActionDTOs(operation.AllAuditActions()),
		Categories: ToCategoryOptionDTOs(operation.AllAuditCategories()),
		Operations: ToOperationTypeDTOs(operation.AllAuditOperations()),
	}
}

// ToAuditActionDTOs 将审计操作定义列表转换为 DTO 列表
func ToAuditActionDTOs(actions []operation.AuditActionDefinition) []AuditActionDTO {
	result := make([]AuditActionDTO, len(actions))
	for i, a := range actions {
		result[i] = AuditActionDTO{
			Action:      a.Action,
			Operation:   string(a.Operation),
			Category:    a.Category,
			Label:       a.Label,
			Description: a.Description,
		}
	}
	return result
}

// ToCategoryOptionDTOs 将分类选项列表转换为 DTO 列表
func ToCategoryOptionDTOs(options []operation.CategoryOption) []CategoryOptionDTO {
	result := make([]CategoryOptionDTO, len(options))
	for i, o := range options {
		result[i] = CategoryOptionDTO{
			Value: o.Value,
			Label: o.Label,
		}
	}
	return result
}

// ToOperationTypeDTOs 将操作类型选项列表转换为 DTO 列表
func ToOperationTypeDTOs(options []operation.OperationTypeOption) []OperationTypeDTO {
	result := make([]OperationTypeDTO, len(options))
	for i, o := range options {
		result[i] = OperationTypeDTO{
			Value: o.Value,
			Label: o.Label,
		}
	}
	return result
}
