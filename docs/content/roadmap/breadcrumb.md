# 面包屑导航

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:22+4`
- [已实现功能](#已实现功能) `:26+11`
  - [AppBreadcrumb 组件](#appbreadcrumb-组件) `:28+9`
- [组件接口](#组件接口) `:37+14`
  - [Props](#props) `:43+8`
- [路由配置要求](#路由配置要求) `:51+16`
- [代码位置](#代码位置) `:67+10`
- [效果预览](#效果预览) `:77+5`

<!--TOC-->

## 需求背景

管理后台需要显示当前页面的层级位置，帮助用户了解所在位置并快速导航到上级页面。

## 已实现功能

### AppBreadcrumb 组件

- 自动根据路由生成面包屑
- 使用路由 `meta.title` 作为显示文本
- 支持首页图标自定义
- 可点击导航到父级页面
- 语义化 HTML 结构（`<nav>` + `<ol>`）
- 良好的可访问性（`aria-label`）

## 组件接口

```vue
<AppBreadcrumb home-icon="mdi-home" home-text="首页" :show-home="true" />
```

### Props

| 属性     | 类型    | 默认值     | 说明         |
| -------- | ------- | ---------- | ------------ |
| homeIcon | string  | "mdi-home" | 首页图标     |
| homeText | string  | "首页"     | 首页文本     |
| showHome | boolean | true       | 是否显示首页 |

## 路由配置要求

路由需要在 `meta` 中定义 `title` 字段：

```typescript
{
  path: "users",
  name: "AdminUsers",
  component: () => import("@/pages/admin/users/index.vue"),
  meta: {
    title: "用户管理",
    icon: "mdi-account",
  },
}
```

## 代码位置

```
web/src/
├── components/
│   └── AppBreadcrumb.vue   # 面包屑组件
└── layout/
    └── AdminLayout.vue     # 已集成面包屑
```

## 效果预览

```
首页 > 管理后台 > 用户管理
```
