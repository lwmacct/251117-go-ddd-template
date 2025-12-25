// Package menu 定义菜单领域模型。
//
// 本包管理系统导航菜单的领域逻辑，定义了：
//   - [Menu]: 菜单实体，支持树形结构
//   - [CommandRepository]: 写仓储接口（创建、更新、删除）
//   - [QueryRepository]: 读仓储接口（查询、树形构建）
//   - 菜单领域错误（见 errors.go）
//
// 树形结构：
// 通过 ParentID 实现无限级菜单嵌套：
//   - [Menu.IsRoot]: 检查是否为顶级菜单（ParentID 为 nil）
//   - [Menu.HasChildren]: 检查是否有子菜单
//   - [Menu.AddChild]: 添加子菜单
//   - [Menu.FindChild]: 递归查找子菜单
//
// 可见性控制：
//   - [Menu.Visible]: 控制菜单显示/隐藏
//   - [Menu.Order]: 定义同级菜单的显示顺序
//
// RBAC 集成：
// [Menu] 实体通常与权限系统配合使用，实现基于角色的菜单过滤。
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package menu
