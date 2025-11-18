# Login 页面架构说明

## 📁 目录结构

```mermaid
graph TD
    A[pages/Login/] --> B[index.vue]
    A --> C[types.ts]
    A --> D[composables/]
    A --> E[components/]
    D --> F[useLogin.ts]
    E --> G[LoginForm.vue]
    E --> H[TwoFactorForm.vue]
```

## 🎯 页面职责

提供用户登录功能，采用简洁的居中布局设计

## 🎨 布局设计

```mermaid
graph TD
    A[登录页面] --> B[居中卡片]
    B --> C[返回首页按钮]
    B --> D[登录标题]
    B --> E[登录表单]
    B --> F[测试账号提示]
    B --> G[注册入口]
    E --> H[账号输入]
    E --> I[密码输入]
    E --> J[图形验证码]
    E --> K[记住我选项]
    E --> L[登录按钮]
    E --> M[错误提示]
```

### 登录表单组件 (LoginForm.vue)

- 返回首页按钮（卡片内顶部）
- 登录标题
- 账号输入（支持手机号/用户名/邮箱）
- 密码输入
- 图形验证码
- 记住我选项
- 登录按钮
- 错误提示区域
- 测试账号提示（开发环境）
- 注册入口

## 🔐 安全特性

- 图形验证码验证
- 验证码自动刷新
- 验证码一次性使用
- 登录失败自动刷新验证码

## 📋 类型定义 (types.ts)

```mermaid
classDiagram
    class LoginForm {
        +string email
        +string password
        +boolean rememberMe
        +string captcha_id
        +string captcha
    }

    class LoginResponse {
        +boolean success
        +string token
        +string message
        +User user
    }

    class CaptchaData {
        +string id
        +string image
        +number expire_at
    }

    class LoginPageData {
        +string pageTitle
        +string pageIcon
        +string backgroundGradient
    }
```

## 📦 状态管理 (composables/useLogin.ts)

页面级状态管理，使用 Composition API 实现，不依赖 Pinia。

**注意**: 这是页面级 composable，每次调用创建新实例。如果需要全局状态，请使用 `/src/stores/useAuthStore`。

```mermaid
stateDiagram-v2
    [*] --> 加载验证码
    加载验证码 --> 未登录: 验证码加载完成
    未登录 --> 登录中: 提交表单
    登录中 --> 验证验证码: 验证中
    验证验证码 --> 已登录: 验证成功
    验证验证码 --> 登录失败: 验证失败
    登录失败 --> 加载验证码: 刷新验证码
    加载验证码 --> 未登录: 重新加载
    已登录 --> 未登录: 登出
```

### 使用方式

```typescript
import { useLogin } from "@/views/Auth/Login/composables";

const login = useLogin();
await login.fetchCaptcha();
const result = await login.login();
```

## 🔄 数据流

```mermaid
sequenceDiagram
    participant U as User
    participant V as index.vue
    participant S as Store
    participant API as Backend API
    participant Redis as Redis Cache

    V->>API: GET /api/platform/auth/captcha
    API->>Redis: 存储验证码
    API-->>V: 返回验证码图片
    U->>V: 输入凭据和验证码
    U->>V: 点击登录
    V->>S: login(form)
    S->>API: POST /api/platform/auth/login
    API->>Redis: 验证验证码
    Redis-->>API: 验证结果
    API->>Redis: 删除验证码
    API-->>S: 返回token
    S-->>V: 登录成功
    V-->>U: 跳转仪表板

    Note over V,Redis: 登录失败时自动刷新验证码
```

## 🎨 UI 组件

- 用户名输入框
- 密码输入框
- 验证码输入框
- 验证码图片（可点击刷新）
- 记住我复选框
- 登录按钮
