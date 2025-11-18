# Register 页面架构说明

## 📁 目录结构

```mermaid
graph TD
    A[pages/Register/] --> B[index.vue]
    A --> C[types.ts]
    A --> D[composables/]
    A --> E[components/]
    D --> F[useRegister.ts]
    E --> G[RegisterForm.vue]
    E --> H[VerifyEmailForm.vue]
```

## 🎯 页面职责

提供用户注册功能，采用简洁的居中布局设计 (风格与登录页面一致)

**注意**: 邮箱验证功能已合并到此页面，不再需要独立路由。如果直接访问验证链接 (带 `email` 和 `code` query 参数) ，页面会自动显示验证表单。

## 🎨 布局设计

```mermaid
graph TD
    A[注册页面] --> B[居中卡片]
    B --> C[返回首页按钮]
    B --> D[注册标题]
    B --> E[注册表单]
    B --> F[验证邮箱表单]
    E --> G[邮箱输入]
    E --> H[密码输入]
    E --> I[确认密码输入]
    E --> J[图形验证码]
    E --> K[注册按钮]
    E --> L[错误/成功提示]
    E --> M[登录入口]
    F --> N[邮箱显示]
    F --> O[验证码输入]
    F --> P[验证按钮]
    F --> Q[返回注册]
```

### 注册表单组件 (RegisterForm.vue)

- 返回首页按钮 (卡片内顶部)
- 注册标题
- 邮箱输入框 (格式验证)
- 密码输入框 (至少 6 个字符)
- 确认密码输入框 (密码匹配验证)
- 图形验证码 (可点击刷新)
- 注册按钮
- 错误/成功提示区域
- 用户协议提示
- 跳转登录按钮

### 邮箱验证表单组件 (VerifyEmailForm.vue)

- 返回注册按钮 (卡片内顶部)
- 验证邮箱标题
- 邮箱地址显示 (从注册表单或 URL 参数获取)
- 验证码输入框 (6 位数字，自动验证)
- 验证按钮
- 成功状态显示 (验证成功后)
- 支持从 URL 参数初始化 (独立访问场景)

## 🔐 安全特性

- 图形验证码验证
- 验证码自动刷新
- 验证码一次性使用
- 密码强度要求 (至少 6 个字符)
- 密码确认验证
- 注册失败自动刷新验证码
- 邮箱验证码验证 (6 位数字)
- **密码管理器兼容性** (1Password、LastPass、Bitwarden 等)
  - 使用标准 `autocomplete` 属性
  - 正确的 `name` 属性设置
  - 支持密码自动生成和填充

## 📋 类型定义 (types.ts)

```mermaid
classDiagram
    class RegisterForm {
        +string username
        +string email
        +string password
        +string confirmPassword
        +string nickname
        +string captcha_id
        +string captcha
    }

    class RegisterResponse {
        +boolean success
        +string message
        +User user
    }

    class CaptchaData {
        +string id
        +string image
        +number expire_at
    }

    class RegisterPageData {
        +string pageTitle
        +string pageIcon
        +string backgroundGradient
    }
```

## 📦 状态管理 (composables/useRegister.ts)

页面级状态管理，使用 Composition API 实现，不依赖 Pinia。

**注意**: 这是页面级 composable，每次调用创建新实例。

```mermaid
stateDiagram-v2
    [*] --> 加载验证码
    加载验证码 --> 未注册: 验证码加载完成
    未注册 --> 注册中: 提交表单
    注册中 --> 验证表单: 验证中
    验证表单 --> 验证密码: 表单有效
    验证密码 --> 验证验证码: 密码匹配
    验证验证码 --> 注册成功: 验证通过
    验证验证码 --> 注册成功: 保存session_token
    注册成功 --> 注册失败: 验证失败
    注册失败 --> 加载验证码: 刷新验证码
    注册成功 --> 邮箱验证: 需要验证
    邮箱验证 --> 验证中: 输入验证码
    验证中 --> 验证成功: 验证通过
    验证中 --> 验证失败: 验证失败
    验证失败 --> 邮箱验证: 重新输入
    验证成功 --> 跳转登录: 3秒后
```

### 使用方式

```typescript
import { useRegister } from "@/views/Auth/Register/composables";

const register = useRegister();
await register.fetchCaptcha();
const result = await register.register();
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
    U->>V: 填写注册信息
    U->>V: 点击注册
    V->>S: register(form)
    S->>S: 验证表单
    S->>S: 验证密码匹配
    S->>API: POST /api/platform/auth/register
    API->>Redis: 验证验证码
    Redis-->>API: 验证结果
    API->>Redis: 删除验证码
    API->>API: 创建用户
    API->>API: 发送验证邮件
    API-->>S: 返回session_token
    S-->>V: 注册成功，需要验证
    V-->>U: 显示验证表单
    U->>V: 输入邮箱验证码
    V->>API: POST /api/platform/auth/verify-email
    API->>Redis: 验证验证码
    Redis-->>API: 验证结果
    API-->>V: 验证成功
    V-->>U: 跳转登录页

    Note over V,Redis: 注册失败时自动刷新验证码
    Note over V,API: 邮箱验证失败时可重新输入
```

## 🎨 UI 组件架构

### 注册表单组件 (RegisterForm.vue)

- 返回首页按钮 (卡片内顶部)
- 注册标题
- **邮箱输入框** (必填，格式验证)
  - `type="email"`, `name="email"`, `autocomplete="email"`
- **密码输入框** (必填，至少 6 个字符)
  - `type="password"`, `name="password"`, `autocomplete="new-password"`
  - 支持密码管理器生成强密码
- **确认密码输入框** (必填，需要匹配)
  - `type="password"`, `name="confirm-password"`, `autocomplete="new-password"`
  - 支持密码管理器自动填充
- **验证码输入框** (必填)
  - `autocomplete="off"` 防止密码管理器干扰
- 验证码图片 (可点击刷新)
- 注册按钮
- 错误/成功提示区域
- 用户协议提示
- 跳转登录按钮

### 邮箱验证表单组件 (VerifyEmailForm.vue)

- 返回注册按钮 (卡片内顶部)
- 验证邮箱标题
- 邮箱地址显示 (来自注册表单或 URL 参数)
- 验证码输入框 (6 位数字，自动验证)
  - `autocomplete="one-time-code"` 支持密码管理器识别
- 验证按钮
- 成功状态显示 (验证成功后自动跳转)

## 🔄 表单验证

### 实时验证

- ✅ 邮箱包含 @
- ✅ 密码长度 >= 6
- ✅ 确认密码与密码一致
- ✅ 验证码不为空
- ✅ 邮箱验证码长度为 6 位

### 错误提示

- 实时显示错误信息
- 5 秒后自动消失
- 可手动关闭

## 🚀 使用流程

### 注册流程

1. 页面加载，自动获取验证码
2. 用户填写注册信息
3. 实时验证表单字段
4. 提交注册请求
5. 成功后显示提示，切换到验证表单
6. 用户输入邮箱验证码
7. 验证成功后跳转登录页

### 独立访问验证 (从邮件链接)

1. 用户点击邮件中的验证链接 (带 `email` 和 `code` query 参数)
2. 页面自动识别 URL 参数，直接显示验证表单
3. 如果 URL 中有验证码，自动填充并验证
4. 验证成功后跳转登录页

## 📝 注意事项

- 邮箱验证已合并到 Register 页面，不再需要独立路由 `/auth/verify-email`
- 支持从邮件链接直接访问验证 (通过 query 参数)
- Session token 用于注册流程中的验证
- 独立访问场景使用邮箱进行验证
