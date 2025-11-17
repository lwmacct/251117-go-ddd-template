# Go DDD Template

[![Go Version](https://img.shields.io/badge/Go-1.25.4+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-vitepress-3eaf7c.svg)](https://lwmacct.github.io/251117-go-ddd-template/)

基于 Go 的领域驱动设计（DDD）模板应用，采用整洁架构原则。提供完整的用户认证、数据库集成、缓存管理等功能。

## ✨ 核心特性

- 🏗️ DDD 整洁架构
- 🔐 JWT 认证授权
- 🗄️ PostgreSQL + GORM
- ⚡ Redis 缓存
- ⚙️ Koanf 配置管理
- 🚀 生产就绪

## 🚀 快速开始

```bash
# 克隆项目
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 启动依赖服务
docker-compose up -d

# 运行应用
task go:run -- api
# 或：go run main.go api

# 测试
curl http://localhost:8080/health
```

## 📚 文档

完整文档：**https://lwmacct.github.io/251117-go-ddd-template/**

- [快速开始](https://lwmacct.github.io/251117-go-ddd-template/guide/getting-started)
- [项目架构](https://lwmacct.github.io/251117-go-ddd-template/guide/architecture)
- [API 文档](https://lwmacct.github.io/251117-go-ddd-template/api/)

## 🛠️ 技术栈

Go 1.25.4 · Gin · PostgreSQL · Redis · GORM · JWT · Koanf · Docker

## 🤝 贡献

查看[贡献指南](https://lwmacct.github.io/251117-go-ddd-template/guide/contributing)

## 📄 许可证

[MIT License](LICENSE)
