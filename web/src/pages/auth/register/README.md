# 注册页面 (Register Page)

用户注册页面模块，提供完整的用户注册功能。

## 📁 目录结构

```
register/
├── components/
│   └── RegisterForm.vue   # 注册表单组件
├── composables/
│   ├── index.ts           # 统一导出
│   └── useRegister.ts     # 注册逻辑 composable
├── index.vue              # 注册页面入口
├── types.ts               # 页面专用类型（邮箱验证等）
└── README.md
```

## 🚀 使用方式

### 基本使用

注册页面通过路由访问：`/auth/register`

### 组件结构

- **index.vue**：页面布局容器（左侧装饰 + 右侧表单）
- **RegisterForm.vue**：核心表单组件
- **useRegister**：封装注册逻辑的 composable

## 🧩 Composables

### useRegister

提供注册相关的状态和方法。

```typescript
import { useRegister } from './composables'

const {
  formData,          // 表单数据
  confirmPassword,   // 确认密码
  loading,           // 加载状态
  errorMessage,      // 错误消息
  handleRegister,    // 执行注册
  clearError,        // 清除错误
  resetForm,         // 重置表单
} = useRegister()
```

## 🎨 功能特性

- 用户名、邮箱、密码注册
- 姓名（可选字段）
- 密码确认验证
- 密码可见性切换
- 实时表单验证
- 友好的错误提示
- 响应式设计（移动端适配）
- 注册成功自动跳转

## ✅ 验证规则

- **用户名**：3-50 字符，只能包含字母、数字和下划线
- **邮箱**：标准邮箱格式
- **密码**：至少 6 个字符
- **确认密码**：必须与密码一致
- **姓名**：可选，最多 100 个字符

## 🔮 未来扩展

- [ ] 邮箱验证（VerifyEmailForm.vue）
- [ ] 手机号注册
- [ ] 图形验证码
- [ ] 用户协议和隐私政策同意

## 📝 相关文档

- [共享 API](../../shared/api/)
- [共享类型](../../shared/types/)
- [认证模块总览](../../README.md)
