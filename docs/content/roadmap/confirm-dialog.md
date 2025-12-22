# 可复用确认对话框

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:24+4`
- [已实现功能](#已实现功能) `:28+15`
  - [ConfirmDialog 组件](#confirmdialog-组件) `:30+7`
  - [useConfirm Composable](#useconfirm-composable) `:37+6`
- [组件接口](#组件接口) `:43+21`
  - [ConfirmDialog](#confirmdialog) `:45+6`
  - [Props](#props) `:51+13`
- [Composable 使用](#composable-使用) `:64+23`
- [代码位置](#代码位置) `:87+10`
- [类型配置](#类型配置) `:97+7`

<!--TOC-->

## 需求背景

应用中存在大量确认对话框（删除确认、危险操作确认等），代码重复且风格不统一。需要统一的确认对话框组件和调用方式。

## 已实现功能

### ConfirmDialog 组件

- 三种类型：删除（红色）、警告（黄色）、信息（蓝色）
- 自动图标和颜色配置
- 加载状态支持
- 支持插槽自定义内容

### useConfirm Composable

- Promise 式调用，支持 async/await
- 快捷方法：`confirmDelete`、`confirmDanger`
- 加载状态控制

## 组件接口

### ConfirmDialog

```vue
<ConfirmDialog v-model="visible" type="delete" title="确认删除" message="确定要删除此项吗？" :loading="loading" @confirm="handleConfirm" @cancel="handleCancel" />
```

### Props

| 属性        | 类型                            | 默认值     | 说明           |
| ----------- | ------------------------------- | ---------- | -------------- |
| modelValue  | boolean                         | 必填       | 显示状态       |
| title       | string                          | "确认操作" | 标题           |
| message     | string                          | ""         | 消息内容       |
| type        | "delete" \| "warning" \| "info" | "info"     | 类型           |
| confirmText | string                          | "确认"     | 确认按钮文本   |
| cancelText  | string                          | "取消"     | 取消按钮文本   |
| loading     | boolean                         | false      | 加载状态       |
| persistent  | boolean                         | false      | 点击遮罩不关闭 |

## Composable 使用

```typescript
import { useConfirm } from "@/composables/useConfirm";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

const { visible, options, loading, confirmDelete, handleConfirm, handleCancel } = useConfirm();

// 删除确认
const handleDelete = async (item: Item) => {
  const confirmed = await confirmDelete(item.name);
  if (confirmed) {
    // 执行删除
  }
};
```

```vue
<template>
  <ConfirmDialog v-model="visible" :type="options.type" :title="options.title" :message="options.message" :loading="loading" @confirm="handleConfirm" @cancel="handleCancel" />
</template>
```

## 代码位置

```
web/src/
├── components/
│   └── ConfirmDialog.vue       # 确认对话框组件
└── composables/
    └── useConfirm.ts           # 确认对话框 Composable
```

## 类型配置

| 类型    | 图标             | 按钮颜色       |
| ------- | ---------------- | -------------- |
| delete  | mdi-delete-alert | error (红色)   |
| warning | mdi-alert        | warning (黄色) |
| info    | mdi-help-circle  | primary (蓝色) |
