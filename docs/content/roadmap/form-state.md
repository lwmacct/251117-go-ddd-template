# 表单状态 Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:27+4`
- [已实现功能](#已实现功能) `:31+23`
  - [useForm](#useform) `:33+9`
  - [useFormDirtyGuard](#useformdirtyguard) `:42+6`
  - [useFieldArray](#usefieldarray) `:48+6`
- [使用方式](#使用方式) `:54+80`
  - [基础用法](#基础用法) `:56+28`
  - [与 Vuetify 结合](#与-vuetify-结合) `:84+13`
  - [脏状态保护](#脏状态保护) `:97+13`
  - [动态字段](#动态字段) `:110+24`
- [API](#api) `:134+18`
  - [useForm 返回值](#useform-返回值) `:136+16`
- [代码位置](#代码位置) `:152+7`

<!--TOC-->

## 需求背景

需要统一管理表单的脏状态、验证、提交、重置等，类似 Formik/React Hook Form 的功能。

## 已实现功能

### useForm

- 表单值管理
- 脏状态检测
- 字段验证
- 提交处理
- 重置表单
- 字段属性生成

### useFormDirtyGuard

- 离开页面警告
- 浏览器刷新/关闭保护
- 自定义确认逻辑

### useFieldArray

- 动态字段数组
- 添加/删除/移动项
- 用于可变长度表单

## 使用方式

### 基础用法

```typescript
import { useForm } from "@/composables/useForm";

const form = useForm({
  initialValues: {
    name: "",
    email: "",
  },
  validate: (values) => {
    const errors: Record<string, string> = {};
    if (!values.name) errors.name = "名称必填";
    if (!values.email) errors.email = "邮箱必填";
    return Object.keys(errors).length > 0 ? errors : null;
  },
  onSubmit: async (values) => {
    await api.createUser(values);
  },
  onSuccess: () => {
    toast.success("创建成功");
  },
});

// 使用 getFieldProps 绑定字段
// <v-text-field v-bind="form.getFieldProps('name')" />
```

### 与 Vuetify 结合

```vue
<template>
  <v-form @submit.prevent="form.submit">
    <v-text-field label="名称" v-bind="form.getFieldProps('name')" />
    <v-text-field label="邮箱" v-bind="form.getFieldProps('email')" />
    <v-btn type="submit" :loading="form.isSubmitting" :disabled="!form.isDirty || !form.isValid"> 提交 </v-btn>
    <v-btn @click="form.reset()">重置</v-btn>
  </v-form>
</template>
```

### 脏状态保护

```typescript
import { useForm, useFormDirtyGuard } from "@/composables/useForm";

const form = useForm({ ... });

useFormDirtyGuard({
  isDirty: form.isDirty,
  message: "您有未保存的更改，确定要离开吗？"
});
```

### 动态字段

```typescript
import { useFieldArray } from "@/composables/useForm";

const { fields, append, remove } = useFieldArray<string>([]);

// 添加项
append("新项目");

// 删除项
remove(0);
```

```vue
<template>
  <div v-for="(field, index) in fields" :key="index">
    <v-text-field v-model="fields[index]" />
    <v-btn @click="remove(index)">删除</v-btn>
  </div>
  <v-btn @click="append('')">添加</v-btn>
</template>
```

## API

### useForm 返回值

| 属性          | 类型           | 说明         |
| ------------- | -------------- | ------------ |
| values        | T              | 表单值       |
| errors        | Record         | 错误信息     |
| touched       | Record         | 触摸状态     |
| isDirty       | `Ref<boolean>` | 是否已修改   |
| isValid       | `Ref<boolean>` | 是否有效     |
| isSubmitting  | `Ref<boolean>` | 是否提交中   |
| setFieldValue | Function       | 设置字段值   |
| setFieldError | Function       | 设置字段错误 |
| reset         | Function       | 重置表单     |
| submit        | Function       | 提交表单     |
| getFieldProps | Function       | 获取字段属性 |

## 代码位置

```
web/src/
└── composables/
    └── useForm.ts    # 表单状态管理
```
