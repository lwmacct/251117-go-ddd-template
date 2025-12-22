# 项目结构

本文档详细介绍 Go DDD Template 的目录结构和文件组织方式。

<!--TOC-->

## Table of Contents

- [目录概览](#目录概览) `:25+24`
- [核心目录详解](#核心目录详解) `:49+128`
  - [internal/ - 业务核心](#internal-业务核心) `:51+56`
  - [web/ - 前端应用](#web-前端应用) `:107+36`
  - [docs/ - 文档系统](#docs-文档系统) `:143+18`
  - [configs/ - 配置文件](#configs-配置文件) `:161+8`
  - [testing/ - 测试脚本](#testing-测试脚本) `:169+8`
- [文件命名规范](#文件命名规范) `:177+27`
  - [Go 文件](#go-文件) `:179+11`
  - [前端文件](#前端文件) `:190+14`
- [构建产物](#构建产物) `:204+12`
- [配置优先级](#配置优先级) `:216+7`
- [开发工具配置](#开发工具配置) `:223+7`

<!--TOC-->

## 目录概览

```
251117-go-ddd-template/
├── internal/               # 核心业务代码（Go 私有包）
│   ├── adapters/          # 适配器层
│   ├── application/       # 应用层
│   ├── domain/            # 领域层
│   ├── infrastructure/    # 基础设施层
│   ├── bootstrap/         # 依赖注入
│   └── commands/          # CLI 命令
├── web/                   # 前端应用（Vue 3）
├── docs/                  # 项目文档（VitePress）
├── configs/               # 配置文件
├── testing/               # 测试脚本
├── .taskfile/             # Task 配置
├── .github/               # GitHub Actions
├── .local/                # 本地构建输出
├── docker-compose.yml     # Docker 编排
├── Taskfile.yml           # Task 任务定义
├── main.go                # 应用入口
└── README.md              # 项目说明
```

## 核心目录详解

### internal/ - 业务核心

遵循 Go 语言约定，使用 internal 包确保代码不被外部引用。

#### internal/adapters/ - 适配器层

```
adapters/
└── http/
    ├── handler/       # HTTP 处理器
    ├── middleware/    # 中间件
    ├── response/      # 响应格式化
    └── router.go      # 路由配置
```

#### internal/application/ - 应用层

```
application/
├── user/
│   ├── command/       # 写操作 Use Cases
│   ├── query/         # 读操作 Use Cases
│   └── dto.go         # 数据传输对象
├── role/
├── menu/
└── ...
```

#### internal/domain/ - 领域层

```
domain/
├── user/
│   ├── entity_user.go          # 用户实体
│   ├── command_repository.go   # 写仓储接口
│   ├── query_repository.go     # 读仓储接口
│   └── errors.go               # 领域错误
├── role/
├── auth/
└── ...
```

#### internal/infrastructure/ - 基础设施层

```
infrastructure/
├── persistence/          # 数据持久化
│   ├── user_model.go
│   ├── user_command_repository.go
│   └── user_query_repository.go
├── auth/                # 认证服务
├── config/              # 配置管理
├── database/            # 数据库连接
└── redis/               # Redis 客户端
```

### web/ - 前端应用

Vue 3 + TypeScript + Vuetify 的单页管理应用。

```
web/
├── src/
│   ├── api/               # API 层（按模块划分）
│   │   ├── admin/         # 管理端 API
│   │   ├── auth/          # 认证 API（含 Axios 客户端）
│   │   ├── user/          # 用户端 API
│   │   └── helpers/       # 工具函数
│   ├── types/             # TypeScript 类型定义
│   │   ├── admin/         # 管理端类型
│   │   ├── auth/          # 认证类型
│   │   ├── user/          # 用户类型
│   │   ├── response/      # 响应类型
│   │   └── common/        # 通用类型
│   ├── composables/       # 组合式函数（自动导入）
│   ├── components/        # 共享组件（自动注册）
│   ├── stores/            # Pinia 状态管理
│   ├── router/            # Vue Router 路由
│   ├── layout/            # 布局组件
│   ├── views/             # 共享视图组件
│   ├── pages/             # 页面组件（按模块划分）
│   │   ├── admin/         # 管理后台页面
│   │   ├── auth/          # 认证页面
│   │   └── user/          # 用户中心页面
│   └── utils/             # 工具函数
├── public/                # 静态资源
├── vite.config.ts         # Vite 配置
└── package.json           # 依赖管理
```

详细前端架构请参阅 [前端架构概述](/architecture/frontend-overview)。

### docs/ - 文档系统

基于 VitePress 的文档站点。

```
docs/
├── .vitepress/
│   ├── config.ts       # VitePress 配置
│   ├── theme/          # 主题定制
│   └── config/         # 配置片段
├── content/
│   ├── guide/          # 用户指南
│   ├── architecture/   # 架构文档
│   ├── development/    # 开发文档
│   └── api/            # API 文档
└── package.json
```

### configs/ - 配置文件

```
configs/
├── config.example.yaml  # 配置示例
└── config.yaml          # 实际配置（gitignore）
```

### testing/ - 测试脚本

```
testing/
├── api_*.py            # API 测试脚本
└── requirements.txt    # Python 依赖
```

## 文件命名规范

### Go 文件

| 类型     | 命名规范                      | 示例                         |
| -------- | ----------------------------- | ---------------------------- |
| 实体     | `entity_{模块}.go`            | `entity_user.go`             |
| 仓储接口 | `{操作}_repository.go`        | `command_repository.go`      |
| 仓储实现 | `{模块}_{操作}_repository.go` | `user_command_repository.go` |
| Use Case | `{动作}_{模块}.go`            | `create_user.go`             |
| Handler  | `{模块}.go`（单数）           | `user.go`                    |
| 模型     | `{模块}_model.go`             | `user_model.go`              |

### 前端文件

| 类型                 | 命名规范            | 示例                                  |
| -------------------- | ------------------- | ------------------------------------- |
| 页面入口             | `index.vue`         | `pages/admin/users/index.vue`         |
| 页面组件             | PascalCase          | `UserDialog.vue`, `RoleSelector.vue`  |
| 共享组件             | PascalCase          | `ConfirmDialog.vue`, `EmptyState.vue` |
| Composable（全局）   | camelCase，use 前缀 | `useClipboard.ts`, `useConfirm.ts`    |
| Composable（页面级） | camelCase，use 前缀 | `useAdminUsers.ts`, `useLogin.ts`     |
| API 模块             | camelCase           | `users.ts`, `roles.ts`                |
| 类型定义             | camelCase           | `user.ts`, `auth.ts`                  |
| Store                | camelCase，use 前缀 | `auth.ts`（导出 `useAuthStore`）      |
| 工具函数             | camelCase           | `storage.ts`, `validation.ts`         |

## 构建产物

```
.local/                 # Task 构建输出
├── bin/               # 二进制文件
└── tmp/               # 临时文件

dist/                  # 前端构建输出
web/dist/              # Vue 应用
docs/.vitepress/dist/  # 文档站点
```

## 配置优先级

1. CLI 参数（最高优先级）
2. 环境变量（APP\_ 前缀）
3. 配置文件（config.yaml）
4. 默认值（最低优先级）

## 开发工具配置

```
.air.toml              # Air 热重载配置
.pre-commit-config.yaml # Pre-commit hooks
.gitignore             # Git 忽略规则
```
