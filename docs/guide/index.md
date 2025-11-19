# 用户指南

欢迎使用 Go DDD Template！本指南将帮助您快速上手并掌握应用的使用方法。

## 📚 指南概览

本指南包含以下内容，帮助您从零开始使用和部署应用：

### 快速入门

- **[快速开始](./getting-started)** - 5 分钟快速上手，从安装到运行第一个请求
- **[配置系统](./configuration)** - 理解配置优先级，学会自定义配置
- **[CLI 命令](./cli-commands)** - 掌握所有命令行工具

### 部署运维

- **[应用部署](./application-deployment)** - 生产环境部署指南
- **[文档部署](./docs-deployment)** - 部署 VitePress 文档到 GitHub Pages

### 开发指南

- **[测试指南](./testing)** - 编写和运行测试
- **[贡献指南](./contributing)** - 如何为项目贡献代码

### 示例

- **[Mermaid 图表](./mermaid-examples)** - 图表示例和最佳实践
- **[任务列表示例](./task-list-examples)** - 任务列表使用示例

## 🎯 我应该从哪里开始？

### 我是新用户

1. 从 [快速开始](./getting-started) 开始
2. 了解 [配置系统](./configuration)
3. 查看 [CLI 命令](./cli-commands)

### 我要部署到生产环境

1. 阅读 [应用部署](./application-deployment)
2. 了解配置最佳实践
3. 查看安全注意事项

### 我要贡献代码

1. 阅读 [贡献指南](./contributing)
2. 了解 [测试指南](./testing)
3. 查看 [架构文档](/architecture/)

## 🔗 相关文档

- **[架构文档](/architecture/)** - 深入了解系统设计和实现原理
- **[API 文档](/api/)** - API 接口详细说明
- **[前端文档](/frontend/)** - Vue 3 前端开发指南
- **[开发文档](/development/)** - VitePress 文档系统使用

## 💡 快速链接

### 常用命令

```bash
# 启动开发服务器
task dev

# 运行数据库迁移
go run . migrate

# 查看所有命令
go run . --help
```

### 环境变量

```bash
# 服务器地址
APP_SERVER_ADDR=:8080

# 数据库连接
APP_DATA_PGSQL_URL=postgresql://user:pass@localhost:5432/db

# Redis 连接
APP_DATA_REDIS_URL=redis://localhost:6379/0

# JWT 密钥
APP_JWT_SECRET=your-secret-key
```

### 配置文件

```bash
# 默认配置文件位置
configs/config.yaml

# 示例配置
configs/config.example.yaml
```

## ❓ 需要帮助？

- 查看 [常见问题](./getting-started#常见问题)
- 阅读 [架构文档](/architecture/) 了解设计原理
- 查看 [API 文档](/api/) 了解接口详情
- 提交 [GitHub Issue](https://github.com/lwmacct/251117-go-ddd-template/issues)

---

准备好了吗？从 [快速开始](./getting-started) 开始您的旅程！
