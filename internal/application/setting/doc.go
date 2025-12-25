// Package setting 实现系统设置的应用层用例。
//
// 本包提供 CQRS 模式的 Command 和 Query Handler：
//
// # Command（写操作）
//
//   - [command.CreateSettingHandler]: 创建设置项
//   - [command.UpdateSettingHandler]: 更新设置值
//   - [command.DeleteSettingHandler]: 删除设置项
//   - [command.BatchUpdateSettingsHandler]: 批量更新设置
//
// # Query（读操作）
//
//   - [query.GetSettingHandler]: 获取设置详情
//   - [query.ListSettingsHandler]: 设置列表查询（支持分类筛选）
//
// # DTO 与映射
//
// 请求 DTO：
//   - [CreateSettingDTO]: 创建设置请求
//   - [UpdateSettingDTO]: 更新设置请求
//   - [BatchUpdateSettingsDTO]: 批量更新请求
//
// 响应 DTO：
//   - [SettingResponse]: 设置信息响应
//
// 映射函数：
//   - [ToSettingResponse]: Setting -> SettingResponse
//
// 设置分类：
// 设置项通过 Category 字段分组（如 system、security、notification 等）。
//
// 类型安全：
// 设置值支持多种类型（string、int、bool、json），通过领域层方法解析。
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package setting
