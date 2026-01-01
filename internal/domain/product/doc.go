// Package product 定义产品领域模型和仓储接口。
//
// # Overview
//
// 本包是多租户系统中可订阅产品的领域核心，定义了：
//   - [Product]: 产品实体
//   - [Type]: 产品类型（个人/团队）
//   - [Status]: 产品状态
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - 产品领域错误（见 errors.go）
//
// # 产品类型
//
// 产品分为两种类型：
//   - [TypePersonal]: 个人产品，由用户直接订阅
//   - [TypeTeam]: 团队产品，由组织订阅，支持席位管理
//
// # Usage
//
//	product := &Product{
//	    Code:      "teamTask",
//	    Name:      "团队任务",
//	    Type:      TypeTeam,
//	    LayoutRef: "TeamTaskLayout",
//	    MaxSeats:  0, // 0 表示无限制
//	    TrialDays: 14,
//	}
//
//	if product.IsTeamProduct() {
//	    // 团队产品需要席位管理
//	}
//
// # Thread Safety
//
// 所有导出函数都是并发安全的。
//
// # 依赖关系
//
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
//
// # 权限域
//
// 产品管理权限域：sys（系统管理员管理）
// 产品访问权限：订阅后可访问
package product
