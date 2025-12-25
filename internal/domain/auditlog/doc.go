// Package auditlog 定义审计日志领域模型。
//
// 审计日志用于记录系统中的关键操作，本包定义了：
//   - [AuditLog]: 审计日志实体
//   - [FilterOptions]: 日志查询过滤条件
//   - [CommandRepository]: 写仓储接口（记录日志）
//   - [QueryRepository]: 读仓储接口（查询、统计）
//   - 审计日志领域错误（见 errors.go）
//
// 记录维度：
//   - 操作者：UserID、Username
//   - 操作：Action（create、update、delete、login 等）
//   - 资源：Resource、ResourceID
//   - 上下文：IPAddress、UserAgent
//   - 状态：Status（success、failed、pending）
//
// 查询能力：
// [FilterOptions] 提供灵活的过滤条件：
//   - 按用户、操作类型、资源类型过滤
//   - 按时间范围过滤（StartDate、EndDate）
//   - 分页支持（Page、Limit）
//
// 合规支持：
// 满足 SOC2、GDPR 等合规性审计要求。
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package auditlog
