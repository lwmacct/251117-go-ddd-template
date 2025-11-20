# Go DDD Template

[![Go Version](https://img.shields.io/badge/Go-1.25.4+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-vitepress-3eaf7c.svg)](https://lwmacct.github.io/251117-go-ddd-template/)

基于 Go 的领域驱动设计 (DDD) 模板应用，采用 DDD 四层架构 + CQRS 模式，提供认证授权、RBAC 权限、审计日志等企业级特性。

## 📚 完整文档

访问在线文档查看详细的架构说明、API 参考、开发指南等内容：

**https://lwmacct.github.io/251117-go-ddd-template/**

## 🏗️ 项目组成

本仓库是一个 Monorepo，包含三个主要部分：

- **后端服务** (`/internal`, `/cmd`, `main.go`) - Go 应用核心代码
- **前端项目** (`/web`) - Vue 3 单页应用
- **文档站点** (`/docs`) - VitePress 2.0 文档系统

## 🛠️ 开发约定

### Pre-commit 文档构建检查

本项目使用 [pre-commit](https://pre-commit.com/) 框架在提交前自动验证文档构建。

#### 快速开始

```bash
# 1. 安装 pre-commit 工具
pip install pre-commit

# 2. 在项目中安装 Git hooks
pre-commit install
```

#### 工作原理

安装后，每次 `git commit` 时：

- ✅ 自动检测 `docs/` 目录下的文件变更
- ✅ 如果有文档修改，自动运行 VitePress 构建验证
- ✅ 只有构建成功才允许提交，确保文档质量

#### 常用命令

```bash
# 手动运行文档构建检查
pre-commit run docs-build --all-files

# 跳过文档构建检查（不推荐）
git commit --no-verify
# 或者
SKIP=docs-build git commit -m "message"

# 更新 hooks 到最新版本
pre-commit autoupdate
```

#### 扩展配置

如需添加更多检查（如 Go 代码格式化、linting 等），可以编辑 `.pre-commit-config.yaml` 文件。参考官方文档：https://pre-commit.com/hooks.html
