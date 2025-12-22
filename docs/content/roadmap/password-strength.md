# 密码强度指示器

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:24+4`
- [已实现功能](#已实现功能) `:28+24`
  - [可视化强度指示](#可视化强度指示) `:30+6`
  - [密码要求检查清单](#密码要求检查清单) `:36+8`
  - [集成位置](#集成位置) `:44+8`
- [技术实现](#技术实现) `:52+32`
  - [组件接口](#组件接口) `:54+9`
  - [Props](#props) `:63+7`
  - [代码位置](#代码位置) `:70+14`
- [强度算法](#强度算法) `:84+17`

<!--TOC-->

## 需求背景

用户在设置密码时需要直观的强度反馈，帮助用户创建安全的密码。

## 已实现功能

### 可视化强度指示

- 三级强度显示：弱（红色）、中（黄色）、强（绿色）
- 进度条动态反映强度等级
- 图标提示增强辨识度

### 密码要求检查清单

- 至少 8 个字符
- 包含小写字母
- 包含大写字母
- 包含数字
- 包含特殊字符

### 集成位置

1. **修改密码页面** (`/user/security`)
   - 完整显示强度条和要求清单

2. **用户创建对话框** (`/admin/users`)
   - 仅显示强度条（创建模式）

## 技术实现

### 组件接口

```vue
<PasswordStrengthIndicator
  :password="password"
  :show-hints="true"  // 是否显示要求清单
/>
```

### Props

| 属性      | 类型    | 默认值 | 说明             |
| --------- | ------- | ------ | ---------------- |
| password  | string  | 必填   | 要检测的密码     |
| showHints | boolean | true   | 是否显示要求清单 |

### 代码位置

```
web/src/
├── components/
│   └── PasswordStrengthIndicator.vue    # 可复用组件
├── utils/auth/
│   └── validation.ts                    # 强度检测函数
├── pages/user/security/components/
│   └── PasswordSettings.vue             # 修改密码集成
└── pages/admin/users/components/
    └── UserDialog.vue                   # 用户创建集成
```

## 强度算法

```typescript
function checkPasswordStrength(password: string): "weak" | "medium" | "strong" {
  let strength = 0;

  if (/[a-z]/.test(password)) strength++; // 小写字母
  if (/[A-Z]/.test(password)) strength++; // 大写字母
  if (/\d/.test(password)) strength++; // 数字
  if (/[^a-zA-Z0-9]/.test(password)) strength++; // 特殊字符
  if (password.length >= 8) strength++; // 长度

  if (strength <= 2) return "weak";
  if (strength <= 3) return "medium";
  return "strong";
}
```
