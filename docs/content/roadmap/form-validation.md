# 表单验证工具

<!--TOC-->

- [需求背景](#需求背景) `:20+4`
- [已实现功能](#已实现功能) `:24+18`
  - [验证规则](#验证规则) `:26+16`
- [使用方式](#使用方式) `:42+31`
  - [与 Vuetify 表单配合](#与-vuetify-表单配合) `:44+6`
  - [手动验证](#手动验证) `:50+12`
  - [验证整个对象](#验证整个对象) `:62+11`
- [代码位置](#代码位置) `:73+7`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

统一表单验证规则，减少重复代码，提供与 Vuetify 兼容的验证函数。

## 已实现功能

### 验证规则

- `required` - 必填
- `minLength` / `maxLength` / `lengthBetween` - 长度验证
- `email` - 邮箱格式
- `phone` - 手机号格式
- `url` - URL 格式
- `number` / `integer` - 数字验证
- `min` / `max` - 数值范围
- `pattern` - 正则表达式
- `username` - 用户名格式
- `password` - 密码强度
- `sameAs` - 确认匹配
- `chinese` - 中文字符
- `idCard` - 身份证号

## 使用方式

### 与 Vuetify 表单配合

```vue
<v-text-field v-model="email" :rules="[rules.required(), rules.email()]" />
```

### 手动验证

```typescript
import { validate, rules } from "@/utils/validation";

const error = validate(email, [rules.required(), rules.email()]);

if (error) {
  console.log(error); // 错误消息
}
```

### 验证整个对象

```typescript
import { validateObject, rules } from "@/utils/validation";

const errors = validateObject(formData, {
  email: [rules.required(), rules.email()],
  password: [rules.required(), rules.password()],
});
```

## 代码位置

```
web/src/
└── utils/
    └── validation.ts    # 表单验证工具
```
