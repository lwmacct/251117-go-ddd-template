# VitePress 高级功能展示

本页面展示 VitePress 的高级功能和自定义组件。

## 🖼️ 图片缩放 (Medium Zoom)

点击下方图片可以放大查看：

![Go Logo](https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png)

**特性**：

- ✅ 点击图片放大
- ✅ 背景自适应主题
- ✅ 响应式设计
- ✅ 路由切换自动重新初始化

## 📡 API 端点展示

使用自定义 `ApiEndpoint` 组件展示 API：

<ApiEndpoint
method="POST"
path="/api/auth/login"
description="用户登录接口"
version="v2.0">

**请求体**：

```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应**：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

</ApiEndpoint>

<ApiEndpoint
method="GET"
path="/api/users/:id"
description="获取用户详情">

**路径参数**：

- `id` (required): 用户 ID

**响应**：

```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "created_at": "2025-11-18T10:00:00Z"
}
```

</ApiEndpoint>

<ApiEndpoint
method="DELETE"
path="/api/users/:id"
description="删除用户 (此接口已废弃)"
deprecated>

请使用 `PUT /api/users/:id` 并设置 `status: inactive`。

</ApiEndpoint>

## 🎯 功能卡片

<FeatureCard
title="JWT 认证"
description="基于 JWT 的用户认证系统"
icon="🔐">

- 支持 Token 刷新
- 自动过期处理
- 安全的密钥管理

</FeatureCard>

<FeatureCard
title="PostgreSQL 集成"
description="使用 GORM 进行数据库操作"
icon="🗄️"
highlighted>

- 自动迁移
- 软删除支持
- 事务管理
- 连接池优化

</FeatureCard>

<FeatureCard
title="Redis 缓存"
description="高性能缓存和分布式锁"
icon="⚡">

- 缓存策略
- 分布式锁
- 过期时间管理

</FeatureCard>

## 📝 步骤指南

<script setup>
const setupSteps = [
  {
    title: '安装依赖',
    description: '使用 Docker Compose 启动 PostgreSQL 和 Redis 服务'
  },
  {
    title: '配置环境变量',
    description: '复制 .env.example 为 .env 并填写配置'
  },
  {
    title: '运行数据库迁移',
    description: '执行 task db:migrate 创建数据表'
  },
  {
    title: '启动应用',
    description: '运行 task go:run -- api 启动 HTTP 服务器'
  }
]
</script>

<StepsGuide :steps="setupSteps" />

## 🎨 主题自定义

本文档系统已自定义以下主题元素：

### 品牌颜色

- **主色调**: `#3eaf7c` <span style="display: inline-block; width: 20px; height: 20px; background: #3eaf7c; border-radius: 4px; vertical-align: middle;"></span>
- **辅助色**: `#42b983` <span style="display: inline-block; width: 20px; height: 20px; background: #42b983; border-radius: 4px; vertical-align: middle;"></span>
- **深色**: `#35495e` <span style="display: inline-block; width: 20px; height: 20px; background: #35495e; border-radius: 4px; vertical-align: middle;"></span>

### UI 增强

- ✅ 外部链接自动添加 ↗ 图标
- ✅ 圆角代码块 (8px)
- ✅ 美化的滚动条
- ✅ 圆角表格
- ✅ 平滑过渡动画

## 🔧 代码实现

### ApiEndpoint 组件

```vue
<ApiEndpoint method="POST" path="/api/users" description="创建新用户" version="v2.0">
  <!-- 你的内容 -->
</ApiEndpoint>
```

**Props**:

- `method`: HTTP 方法 (`GET` | `POST` | `PUT` | `PATCH` | `DELETE`)
- `path`: API 路径
- `description`: 描述 (可选)
- `version`: 版本标记 (可选)
- `deprecated`: 是否废弃 (可选)

### FeatureCard 组件

```vue
<FeatureCard title="功能标题" description="功能描述" icon="🎯" highlighted>
  <!-- 详细内容 -->
</FeatureCard>
```

**Props**:

- `title`: 功能标题
- `description`: 功能描述 (可选)
- `icon`: Emoji 图标 (可选)
- `highlighted`: 是否高亮 (可选)

### StepsGuide 组件

```vue
<script setup>
const steps = [
  {
    title: "步骤 1",
    description: "描述 1",
  },
  {
    title: "步骤 2",
    description: "描述 2",
  },
];
</script>

<StepsGuide :steps="steps" />
```

**Props**:

- `steps`: 步骤数组，每个步骤包含 `title` 和 `description`

## 📚 扩展阅读

- [创建自定义组件](/development/features#自定义组件)
- [主题配置](https://vitepress.dev/reference/default-theme-config)
- [Vue 组件集成](https://vitepress.dev/guide/using-vue)

## 💡 使用建议

1. **API 文档**：使用 `ApiEndpoint` 组件展示 RESTful API
2. **功能展示**：使用 `FeatureCard` 突出核心功能
3. **教程指南**：使用 `StepsGuide` 展示操作步骤
4. **图片展示**：利用 Medium Zoom 提供更好的查看体验
5. **主题定制**：根据品牌调整 CSS 变量

## 🚀 更多可能

你还可以创建更多自定义组件：

- **代码对比组件**：并排展示不同版本的代码
- **时间线组件**：展示项目发展历程
- **状态指示器**：显示服务状态
- **进度追踪**：展示项目完成度
- **交互式演示**：嵌入 CodeSandbox/StackBlitz

所有这些都可以通过 Vue 组件轻松实现！
