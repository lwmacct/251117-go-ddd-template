# 页面过渡动画

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:22+4`
- [已实现功能](#已实现功能) `:26+8`
  - [PageTransition 组件](#pagetransition-组件) `:28+6`
- [组件接口](#组件接口) `:34+18`
  - [Props](#props) `:44+8`
- [可用动画](#可用动画) `:52+12`
- [代码位置](#代码位置) `:64+10`
- [已集成位置](#已集成位置) `:74+3`

<!--TOC-->

## 需求背景

页面切换时增加过渡动画，提升用户体验，使页面切换更加流畅自然。

## 已实现功能

### PageTransition 组件

- 多种预设动画效果
- 支持禁用动画
- 可配置过渡模式

## 组件接口

```vue
<router-view v-slot="{ Component }">
  <PageTransition name="fade">
    <component :is="Component" />
  </PageTransition>
</router-view>
```

### Props

| 属性    | 类型    | 默认值   | 说明     |
| ------- | ------- | -------- | -------- |
| name    | string  | "fade"   | 动画名称 |
| mode    | string  | "out-in" | 过渡模式 |
| enabled | boolean | true     | 是否启用 |

## 可用动画

| 名称        | 效果           |
| ----------- | -------------- |
| fade        | 淡入淡出       |
| slide-left  | 左滑入，左滑出 |
| slide-right | 右滑入，右滑出 |
| slide-up    | 上滑入，上滑出 |
| slide-down  | 下滑入，下滑出 |
| scale       | 缩放效果       |
| none        | 无动画         |

## 代码位置

```
web/src/
├── components/
│   └── PageTransition.vue    # 页面过渡组件
└── layout/
    └── AdminLayout.vue       # 已集成页面过渡
```

## 已集成位置

- ✅ 管理后台布局 (`AdminLayout.vue`)
