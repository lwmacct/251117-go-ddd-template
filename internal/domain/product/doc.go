// Package product 定义产品领域模型和仓储接口。
//
// 本包是多租户系统中可订阅产品的领域核心，定义了：
//   - [Product]: 产品实体
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - 产品领域错误（见 errors.go）
//
// 权限域：sys（系统管理员管理）
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence 包。
package product
