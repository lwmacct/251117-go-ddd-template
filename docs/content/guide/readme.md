# Go DDD Template

基于领域驱动设计（DDD）和 CQRS 模式的企业级 Go 应用模板。

## 🚀 快速开始

```bash
# 克隆项目
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 启动依赖服务
docker-compose up -d

# 运行应用
task go:run -- api
```

详细安装步骤请参阅 [快速开始指南](./getting-started.md)。

## 🏗️ 项目架构

本项目采用 **DDD 四层架构 + CQRS 模式**：

```
Adapters → Application → Domain ← Infrastructure
```

- **Adapters**: HTTP Handler，仅处理请求转换
- **Application**: Use Case Handler，业务编排
- **Domain**: 领域模型和业务规则
- **Infrastructure**: 技术实现细节

详见 [架构文档](/architecture/ddd-cqrs)。

## ✨ 核心特性

- 🏗️ **DDD + CQRS**: 完整的领域驱动设计实现
- 🔐 **双重认证**: JWT + PAT 支持不同场景
- 👥 **RBAC 权限**: 三段式细粒度权限控制
- 📊 **审计日志**: 完整的操作追踪
- 🗄️ **PostgreSQL**: GORM ORM 支持
- ⚡ **Redis 缓存**: 高性能缓存层
- 🐳 **Docker 支持**: 容器化部署
- 📝 **Swagger API**: 自动生成 API 文档

## 📁 项目结构

```
internal/
├── adapters/      # HTTP 适配器层
├── application/   # 应用层 (Use Cases)
├── domain/        # 领域层 (业务模型)
├── infrastructure/# 基础设施层
└── bootstrap/     # 依赖注入容器

web/               # Vue 3 前端应用
docs/              # VitePress 文档
testing/           # 测试脚本
```

## 🔨 开发

### 使用 Taskfile

```bash
# 查看所有可用任务
task -a

# 运行开发服务器
task dev

# 构建项目
task build

# 运行测试
task test

# 代码检查
task lint
```

### 使用 Pre-commit

```bash
# 初始化 pre-commit hooks
pre-commit install

# 手动运行所有检查
pre-commit run --all-files
```

### 热重载开发

```bash
# 安装 air
go install github.com/air-verse/air@latest

# 运行热重载
air
```

详见 [开发文档](/development/setup)。

## 📚 文档

- [用户指南](/guide/) - 快速上手和基础使用
- [架构文档](/architecture/) - 深入理解系统设计
- [API 文档](/api/overview) - 接口规范和示例
- [开发指南](/development/) - 开发环境和最佳实践

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

贡献前请确保：
1. 遵循代码规范
2. 编写充分的测试
3. 更新相关文档

## 📄 许可证

MIT License

## 🔗 链接

- [GitHub 仓库](https://github.com/lwmacct/251117-go-ddd-template)
- [作者主页](https://github.com/lwmacct)
- [语雀文档](https://yuque.com/lwmacct)

---

*本项目是一个持续演进的模板，欢迎社区贡献和改进建议。*
