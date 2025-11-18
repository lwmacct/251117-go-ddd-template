# 登录页面 (Login Page)

用户登录页面模块，提供完整的登录功能和用户体验。

## 📁 目录结构

```
login/
├── components/
│   └── LoginForm.vue      # 登录表单组件
├── composables/
│   ├── index.ts           # 统一导出
│   └── useLogin.ts        # 登录逻辑 composable
├── index.vue              # 登录页面入口
├── types.ts               # 页面专用类型（2FA 等）
└── README.md
```

## 🚀 使用方式

### 基本使用

登录页面通过路由访问：`/auth/login`

### 组件结构

- **index.vue**：页面布局容器（左侧装饰 + 右侧表单）
- **LoginForm.vue**：核心表单组件
- **useLogin**：封装登录逻辑的 composable

## 🧩 Composables

### useLogin

提供登录相关的状态和方法。

```typescript
import { useLogin } from './composables'

const {
  formData,         // 表单数据
  loading,          // 加载状态
  errorMessage,     // 错误消息
  handleLogin,      // 执行登录
  clearError,       // 清除错误
  resetForm,        // 重置表单
} = useLogin()
```

## 🎨 功能特性

- 用户名或邮箱登录
- 密码可见性切换
- 实时表单验证
- 友好的错误提示
- 响应式设计（移动端适配）
- 登录成功自动跳转

## 🔮 未来扩展

- [ ] 双因素认证（TwoFactorForm.vue）
- [ ] 记住我功能
- [ ] 社交账号登录
- [ ] 验证码登录

## 📝 相关文档

- [共享 API](../../shared/api/)
- [共享类型](../../shared/types/)
- [认证模块总览](../../README.md)
