# 空状态组件

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:24+4`
- [已实现功能](#已实现功能) `:28+9`
  - [EmptyState 组件](#emptystate-组件) `:30+7`
- [组件接口](#组件接口) `:37+24`
  - [基础用法](#基础用法) `:39+6`
  - [自定义内容](#自定义内容) `:45+6`
  - [使用插槽](#使用插槽) `:51+10`
- [Props](#props) `:61+13`
- [预设类型](#预设类型) `:74+9`
- [代码位置](#代码位置) `:83+7`

<!--TOC-->

## 需求背景

应用中多处需要显示空状态（列表为空、搜索无结果、加载失败等），需要统一的空状态组件和预设样式。

## 已实现功能

### EmptyState 组件

- 多种预设类型：empty、search、error、no-permission
- 自定义图标、标题、描述
- 可选操作按钮
- 支持插槽自定义内容

## 组件接口

### 基础用法

```vue
<EmptyState type="empty" />
```

### 自定义内容

```vue
<EmptyState icon="mdi-file-outline" title="暂无文件" description="上传文件开始使用" action-text="上传文件" :show-action="true" @action="handleUpload" />
```

### 使用插槽

```vue
<EmptyState type="search">
  <template #action>
    <v-btn variant="text" @click="clearSearch">清空搜索</v-btn>
  </template>
</EmptyState>
```

## Props

| 属性        | 类型                                                          | 默认值           | 说明             |
| ----------- | ------------------------------------------------------------- | ---------------- | ---------------- |
| type        | "empty" \| "search" \| "error" \| "no-permission" \| "custom" | "custom"         | 预设类型         |
| icon        | string                                                        | -                | 自定义图标       |
| title       | string                                                        | -                | 自定义标题       |
| description | string                                                        | -                | 自定义描述       |
| iconSize    | string \| number                                              | 80               | 图标大小         |
| iconColor   | string                                                        | "grey-lighten-1" | 图标颜色         |
| actionText  | string                                                        | -                | 操作按钮文本     |
| showAction  | boolean                                                       | false            | 是否显示操作按钮 |

## 预设类型

| 类型          | 图标                     | 标题       | 描述                       |
| ------------- | ------------------------ | ---------- | -------------------------- |
| empty         | mdi-inbox-outline        | 暂无数据   | 当前列表为空               |
| search        | mdi-file-search-outline  | 未找到结果 | 尝试修改搜索条件或清空筛选 |
| error         | mdi-alert-circle-outline | 加载失败   | 数据加载出错，请稍后重试   |
| no-permission | mdi-lock-outline         | 无权限访问 | 您没有权限查看此内容       |

## 代码位置

```
web/src/
└── components/
    └── EmptyState.vue    # 空状态组件
```
