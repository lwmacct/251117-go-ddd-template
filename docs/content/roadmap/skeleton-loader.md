# 骨架屏加载组件

<!--TOC-->

- [需求背景](#需求背景) `:22:25`
- [已实现功能](#已实现功能) `:26:27`
  - [SkeletonLoader 组件](#skeletonloader-组件) `:28:33`
- [组件接口](#组件接口) `:34:35`
  - [基础用法](#基础用法) `:36:43`
  - [预设类型](#预设类型) `:44:62`
  - [自定义骨架](#自定义骨架) `:63:75`
- [Props](#props) `:76:85`
- [代码位置](#代码位置) `:86:93`
- [使用场景](#使用场景) `:94:99`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

数据加载时显示骨架屏，提供更好的视觉反馈，避免页面空白或闪烁。

## 已实现功能

### SkeletonLoader 组件

- 多种预设类型：text、avatar、button、card、table、list
- 支持自定义骨架
- 条件渲染：加载时显示骨架，完成后显示内容

## 组件接口

### 基础用法

```vue
<SkeletonLoader :loading="isLoading" type="table">
  <ActualContent />
</SkeletonLoader>
```

### 预设类型

```vue
<!-- 文本骨架 -->
<SkeletonLoader type="text" :lines="3" />

<!-- 头像骨架 -->
<SkeletonLoader type="avatar" />

<!-- 表格骨架 -->
<SkeletonLoader type="table" :table-rows="5" :table-cols="4" />

<!-- 列表骨架 -->
<SkeletonLoader type="list" :lines="5" />

<!-- 卡片骨架 -->
<SkeletonLoader type="card" />
```

### 自定义骨架

```vue
<SkeletonLoader :loading="loading" type="custom">
  <template #skeleton>
    <v-skeleton-loader type="image" />
    <v-skeleton-loader type="text" />
  </template>

  <ActualContent />
</SkeletonLoader>
```

## Props

| 属性      | 类型    | 默认值 | 说明              |
| --------- | ------- | ------ | ----------------- |
| type      | string  | "text" | 骨架类型          |
| loading   | boolean | true   | 是否加载中        |
| lines     | number  | 3      | 行数（text/list） |
| tableRows | number  | 5      | 表格行数          |
| tableCols | number  | 4      | 表格列数          |

## 代码位置

```
web/src/
└── components/
    └── SkeletonLoader.vue    # 骨架屏组件
```

## 使用场景

- 列表数据加载
- 表格数据加载
- 卡片内容加载
- 初始页面渲染
