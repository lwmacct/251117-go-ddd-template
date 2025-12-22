# 加载遮罩组件

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:21+4`
- [已实现功能](#已实现功能) `:25+9`
  - [LoadingOverlay 组件](#loadingoverlay-组件) `:27+7`
- [组件接口](#组件接口) `:34+20`
  - [Props](#props) `:42+12`
- [代码位置](#代码位置) `:54+8`
- [使用场景](#使用场景) `:62+6`

<!--TOC-->

## 需求背景

在数据加载或异步操作期间，需要显示加载状态覆盖内容区域，防止用户重复操作。

## 已实现功能

### LoadingOverlay 组件

- 支持绝对定位覆盖父容器
- 可自定义加载文本
- 可配置遮罩透明度和颜色
- 淡入淡出动画效果

## 组件接口

```vue
<LoadingOverlay :loading="isLoading" text="数据加载中...">
  <YourContent />
</LoadingOverlay>
```

### Props

| 属性         | 类型             | 默认值    | 说明           |
| ------------ | ---------------- | --------- | -------------- |
| loading      | boolean          | false     | 是否显示加载   |
| text         | string           | ""        | 加载提示文本   |
| absolute     | boolean          | true      | 是否绝对定位   |
| opacity      | number           | 0.8       | 遮罩透明度     |
| size         | string \| number | 48        | 加载指示器大小 |
| color        | string           | "primary" | 加载指示器颜色 |
| overlayColor | string           | "white"   | 遮罩背景颜色   |

## 代码位置

```
web/src/
└── components/
    └── LoadingOverlay.vue    # 加载遮罩组件
```

## 使用场景

- 卡片内容加载
- 表格数据刷新
- 表单提交中
- 页面初始化
