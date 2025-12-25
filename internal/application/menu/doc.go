// Package menu 实现菜单管理的应用层用例。
//
// 本包提供 CQRS 模式的 Command 和 Query Handler：
//
// # Command（写操作）
//
//   - [command.CreateMenuHandler]: 创建菜单项
//   - [command.UpdateMenuHandler]: 更新菜单信息
//   - [command.DeleteMenuHandler]: 删除菜单项
//   - [command.ReorderMenusHandler]: 重新排序菜单
//
// # Query（读操作）
//
//   - [query.GetMenuHandler]: 获取菜单详情
//   - [query.ListMenusHandler]: 菜单树形列表查询
//
// # DTO 与映射
//
// 请求 DTO：
//   - [CreateMenuDTO]: 创建菜单请求
//   - [UpdateMenuDTO]: 更新菜单请求
//   - [ReorderMenusDTO]: 菜单排序请求
//
// 响应 DTO：
//   - [MenuResponse]: 菜单信息响应
//   - [MenuTreeResponse]: 菜单树形结构响应
//
// 映射函数：
//   - [ToMenuResponse]: Menu -> MenuResponse
//   - [ToMenuTreeResponse]: []Menu -> MenuTreeResponse
//
// 菜单层级：
// 菜单支持父子层级结构，通过 ParentID 建立关联。
// 根菜单 ParentID 为 0。
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package menu
