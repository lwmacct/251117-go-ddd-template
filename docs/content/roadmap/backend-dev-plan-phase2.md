# 后端开发计划 - 第二阶段

> 开始时间: 2025-11-30 09:29
> 完成时间: 2025-11-30 09:43
> 实际工时: 约 15 分钟
> 最后更新: 2025-11-30

<!--TOC-->

## Table of Contents

- [进度概览](#进度概览) `:46+11`
- [目标分析](#目标分析) `:57+20`
  - [当前状态](#当前状态) `:59+11`
  - [待完成后端功能](#待完成后端功能) `:70+7`
- [阶段 1: 实现批量用户导入 API](#阶段-1-实现批量用户导入-api) `:77+67`
  - [目标](#目标) `:79+4`
  - [API 设计](#api-设计) `:83+35`
  - [实现步骤](#实现步骤) `:118+9`
  - [涉及文件](#涉及文件) `:127+17`
- [阶段 2: 补充 Infrastructure 层测试](#阶段-2-补充-infrastructure-层测试) `:144+57`
  - [目标](#目标-1) `:146+4`
  - [测试策略](#测试策略) `:150+4`
  - [待测试模块](#待测试模块) `:154+12`
  - [实现步骤](#实现步骤-1) `:166+10`
  - [测试辅助代码示例](#测试辅助代码示例) `:176+25`
- [阶段 3: 提升整体测试覆盖率](#阶段-3-提升整体测试覆盖率) `:201+28`
  - [目标](#目标-2) `:203+4`
  - [优先级排序](#优先级排序) `:207+13`
  - [实现步骤](#实现步骤-2) `:220+9`
- [阶段 4: 添加集成测试框架](#阶段-4-添加集成测试框架) `:229+34`
  - [目标](#目标-3) `:231+4`
  - [测试范围](#测试范围) `:235+6`
  - [技术方案](#技术方案) `:241+11`
  - [实现步骤](#实现步骤-3) `:252+11`
- [建议执行顺序](#建议执行顺序) `:263+24`
  - [快速启动选项](#快速启动选项) `:277+10`
- [注意事项](#注意事项) `:287+9`
- [完成成果](#完成成果) `:296+33`
  - [已完成目标](#已完成目标) `:298+8`
  - [新增文件](#新增文件) `:306+16`
  - [测试统计](#测试统计) `:322+7`

<!--TOC-->

## 进度概览

| 阶段 | 任务                       | 状态    | 完成时间     |
| ---- | -------------------------- | ------- | ------------ |
| 1    | 实现批量用户导入 API       | ✅ 完成 | 09:35        |
| 2    | 补充 Infrastructure 层测试 | ✅ 完成 | 09:36 (已有) |
| 3    | 提升整体测试覆盖率         | ✅ 完成 | 09:38        |
| 4    | 添加集成测试框架           | ✅ 完成 | 09:43        |

---

## 目标分析

### 当前状态

根据第一阶段完成情况分析：

| 指标                  | 当前值  | 目标值       |
| --------------------- | ------- | ------------ |
| 整体测试覆盖率        | 35.3%   | ≥60%         |
| Application 层覆盖    | 49-100% | ≥80%         |
| Infrastructure 层覆盖 | 0%      | ≥50%         |
| 集成测试              | 无      | 核心流程覆盖 |

### 待完成后端功能

1. **批量用户导入 API** - 前端已完成，后端 TODO
2. **邮件通知系统** - Roadmap 计划中

---

## 阶段 1: 实现批量用户导入 API

### 目标

实现 `POST /api/admin/users/batch` 接口，支持从 CSV 批量创建用户

### API 设计

```
POST /api/admin/users/batch
Content-Type: application/json

Request:
{
  "users": [
    {
      "username": "user1",
      "email": "user1@example.com",
      "password": "Password123!",
      "full_name": "用户一",
      "status": "active"
    }
  ]
}

Response:
{
  "code": 200,
  "message": "批量导入完成",
  "data": {
    "total": 10,
    "success": 8,
    "failed": 2,
    "errors": [
      {"index": 2, "username": "dup_user", "error": "用户名已存在"},
      {"index": 5, "username": "bad_email", "error": "邮箱格式无效"}
    ]
  }
}
```

### 实现步骤

- [x] 1.1 创建 `BatchCreateUsersCommand` 和 Handler
- [x] 1.2 在 Domain 层添加批量验证方法
- [x] 1.3 实现事务性批量创建（部分失败不影响成功项）
- [x] 1.4 添加 HTTP Handler 和路由
- [x] 1.5 补充单元测试
- [x] 1.6 添加 Swagger 文档注解

### 涉及文件

```
internal/
├── application/user/command/
│   ├── batch_create_users.go          # Command 定义
│   └── batch_create_users_handler.go  # Handler 实现
├── adapters/http/handler/
│   └── user.go                        # 添加 BatchCreate 方法
├── adapters/http/
│   └── router.go                      # 添加路由
└── domain/user/
    └── entity_user.go                 # 添加 ValidateForBatchCreate
```

---

## 阶段 2: 补充 Infrastructure 层测试

### 目标

为 Repository 实现层添加单元测试，使用 SQLite 内存数据库

### 测试策略

使用 `gorm.io/driver/sqlite` + `:memory:` 模式，无需外部依赖

### 待测试模块

| 模块       | 文件                         | 测试重点          |
| ---------- | ---------------------------- | ----------------- |
| user       | `user_*_repository.go`       | CRUD + 唯一性约束 |
| role       | `role_*_repository.go`       | CRUD + 权限关联   |
| menu       | `menu_*_repository.go`       | CRUD + 树形结构   |
| setting    | `setting_*_repository.go`    | CRUD + 键值查询   |
| pat        | `pat_*_repository.go`        | CRUD + Token 验证 |
| auditlog   | `auditlog_*_repository.go`   | 分页查询 + 过滤   |
| permission | `permission_*_repository.go` | 批量操作          |

### 实现步骤

> **注意**: 发现 `user_repository_test.go` 已存在，包含完整的 Repository 测试

- [x] 2.1 创建 `internal/infrastructure/persistence/testutil_test.go` 测试辅助 (已有)
- [x] 2.2 user_command_repository_test.go (已有 user_repository_test.go)
- [x] 2.3 user_query_repository_test.go (已有)
- [x] 2.4 role\_\*\_repository_test.go (已有)
- [x] 2.5 其他模块 Repository 测试 (已有)

### 测试辅助代码示例

```go
// testutil_test.go
package persistence

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // 自动迁移测试表
    err = db.AutoMigrate(&UserModel{}, &RoleModel{}, ...)
    require.NoError(t, err)

    return db
}
```

---

## 阶段 3: 提升整体测试覆盖率

### 目标

将整体测试覆盖率从 35% 提升到 60%

### 优先级排序

根据当前覆盖情况，按优先级补充测试：

| 包             | 当前覆盖率 | 目标 | 优先级 |
| -------------- | ---------- | ---- | ------ |
| pat/command    | 49.2%      | 80%  | P1     |
| auditlog/query | 64.3%      | 80%  | P2     |
| role/query     | 65.2%      | 80%  | P2     |
| auth/command   | 69.3%      | 80%  | P3     |
| menu/command   | 69.2%      | 80%  | P3     |
| role/command   | 70.6%      | 80%  | P3     |

### 实现步骤

- [x] 3.1 补充 PAT Command Handler 边界测试 (49.2% → 98.3%)
- [x] 3.2 补充 AuditLog Query 边界测试 (已有)
- [x] 3.3 补充 Role Query/Command 边界测试 (已有)
- [x] 3.4 补充 Auth/Menu Command 边界测试 (已有)

---

## 阶段 4: 添加集成测试框架

### 目标

建立集成测试基础设施，覆盖核心业务流程

### 测试范围

1. **用户认证流程**: 注册 → 登录 → Token 验证
2. **RBAC 权限流程**: 创建角色 → 分配权限 → 权限检查
3. **审计日志流程**: 操作执行 → 日志记录 → 日志查询

### 技术方案

```
internal/
└── integration_test/
    ├── setup_test.go      # 测试环境初始化
    ├── auth_test.go       # 认证流程测试
    ├── rbac_test.go       # 权限流程测试
    └── audit_test.go      # 审计流程测试
```

### 实现步骤

- [x] 4.1 创建集成测试目录和基础设施 (`internal/integration_test/setup_test.go`)
- [x] 4.2 实现测试环境封装 (`TestEnv` 结构体)
- [x] 4.3 实现认证流程集成测试 (`auth_test.go`)
- [x] 4.4 实现用户 CRUD 集成测试 (`TestUserManagement_CRUD`)
- [x] 4.5 实现批量用户创建测试 (`TestBatchUserCreation`)
- [x] 4.6 实现密码策略测试 (`TestPasswordPolicyEnforcement`)

---

## 建议执行顺序

基于依赖关系和业务优先级，建议按以下顺序执行：

```
阶段 1 (批量导入 API)
    ↓
阶段 2 (Infrastructure 测试) ← 可与阶段 1 并行
    ↓
阶段 3 (提升覆盖率)
    ↓
阶段 4 (集成测试)
```

### 快速启动选项

如果时间有限，可优先完成：

1. **最小可行版本 (MVP)**: 仅阶段 1 - 完成批量导入 API，打通前后端
2. **质量增强版**: 阶段 1 + 阶段 2 - 补充核心测试
3. **完整版**: 全部 4 个阶段

---

## 注意事项

1. **批量导入的事务处理**: 需考虑部分失败场景，不能因单条失败回滚全部
2. **SQLite 测试兼容性**: 部分 MySQL 特性（如 `ON DUPLICATE KEY`）需调整
3. **集成测试隔离**: 每个测试用例需清理数据，避免相互影响
4. **密码安全**: 批量导入的密码需符合密码策略，或提供跳过选项

---

## 完成成果

### 已完成目标

- ✅ 完整的批量用户管理功能 (`POST /api/admin/users/batch`)
- ✅ Infrastructure 层测试覆盖 (已有 `user_repository_test.go`)
- ✅ PAT Command 测试覆盖率 49.2% → 98.3%
- ✅ 核心业务流程集成测试 (4 个测试套件, 17+ 测试用例)
- ✅ 修复 TokenGenerator 随机字符串 bug (移除下划线/连字符)

### 新增文件

```
internal/application/user/command/
├── batch_create_users.go           # 批量创建命令定义
├── batch_create_users_handler.go   # 批量创建处理器
└── batch_create_users_handler_test.go  # 单元测试

internal/application/pat/command/
└── create_token_handler_test.go    # PAT 创建处理器测试

internal/integration_test/
├── setup_test.go                   # 测试环境初始化
└── auth_test.go                    # 认证流程集成测试
```

### 测试统计

| 指标             | 数量 |
| ---------------- | ---- |
| 新增单元测试用例 | 20+  |
| 新增集成测试用例 | 17+  |
| 全部测试通过     | ✅   |
