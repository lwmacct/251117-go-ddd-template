# Go DDD Template

[![Go Version](https://img.shields.io/badge/Go-1.25.4+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-vitepress-3eaf7c.svg)](https://lwmacct.github.io/251117-go-ddd-template/)

基于 Go 的领域驱动设计（DDD）模板应用，采用整洁架构原则。

## ✨ 特性

🏗️ DDD 整洁架构 · 🔐 JWT 认证 · 🗄️ PostgreSQL · ⚡ Redis · ⚙️ 配置管理 · 🚀 生产就绪

## 🚀 快速开始

```bash
# 启动依赖服务
docker-compose up -d

# 运行应用
task go:run -- api

# 健康检查
curl http://localhost:8080/health
```

## 📚 文档

**在线文档**: https://lwmacct.github.io/251117-go-ddd-template/

本地开发文档：

```bash
cd docs
npm install  # 首次需要安装依赖
npm run dev  # 访问 http://localhost:5173/docs/
```

## 🛠️ 技术栈

**后端**: Go 1.25.4 · Gin · GORM · JWT · Koanf
**数据**: PostgreSQL · Redis
**工具**: Docker · Taskfile
**文档**: VitePress 2.0 · Vue 3

## 📁 项目结构

```
├── cmd/                # 命令行入口
├── internal/           # 核心代码（DDD 分层架构）
├── configs/            # 配置文件
├── web/                # 前端项目 (Vue3)
├── docs/               # VitePress 文档
└── main.go             # 主入口
```

## 📄 许可证

[MIT License](LICENSE)
