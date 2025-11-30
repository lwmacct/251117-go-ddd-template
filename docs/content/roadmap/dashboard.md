# 数据仪表盘

> **状态**: ✅ 已实现
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

管理员需要一个系统概览页面，快速了解系统运行状况和关键指标。

## 已实现功能

### 统计卡片

- 总用户数
- 活跃用户数
- 禁用用户数
- 封禁用户数
- 总角色数
- 总权限数
- 菜单数量

### 用户状态分布

- 活跃/未激活/封禁用户比例
- 进度条可视化展示

### 最近审计日志

- 显示最近 5 条审计记录
- 包含用户、操作、资源、状态、时间
- 链接到完整审计日志页面

### 快速操作

- 管理用户
- 管理角色
- 管理菜单
- 审计日志

### 系统信息

- 系统名称和版本
- 架构模式说明
- 数据库/Redis 连接状态
- 服务运行状态

## 技术实现

### 后端 API

```
GET /api/admin/overview/stats
Response: {
  "total_users": number,
  "active_users": number,
  "inactive_users": number,
  "banned_users": number,
  "total_roles": number,
  "total_permissions": number,
  "total_menus": number,
  "recent_audit_logs": [...]
}
```

### 前端实现

- 页面: `web/src/pages/admin/overview/index.vue`
- Composable: `web/src/pages/admin/overview/composables/useOverview.ts`
- API: `web/src/api/admin/overview.ts`

## 相关文件

```
web/src/
├── api/admin/
│   └── overview.ts              # 概览 API
├── pages/admin/overview/
│   ├── index.vue                # 概览页面
│   └── composables/
│       └── useOverview.ts       # 页面逻辑

internal/adapters/http/handler/
└── overview.go                  # 后端 Handler
```

## 未来增强

- [ ] 添加图表库（ECharts/Chart.js）实现趋势图
- [ ] 用户注册趋势图
- [ ] 活动热力图
- [ ] 实时数据刷新
