---
paths:
  - "src/**/*.ts"
  - "src/**/*.vue"
---

# API 层开发规范

本目录是前端与后端 API 的桥梁，所有业务类型定义来自 OpenAPI 生成代码。

<!--TOC-->

## Table of Contents

- [核心原则](#核心原则) `:24+6`
- [类型来源](#类型来源) `:30+9`
- [禁止事项](#禁止事项) `:39+17`
- [允许的扩展](#允许的扩展) `:56+16`
- [设计反思原则](#设计反思原则) `:72+31`
- [目录结构](#目录结构) `:103+12`

<!--TOC-->

## 核心原则

**前端不定义与后端重复的 DTO，只消费 OpenAPI 生成的类型。**

发现后端类型缺陷时，**修复后端**，禁止前端补救。

## 类型来源

| 类型分类     | 来源                    | 示例                                   |
| ------------ | ----------------------- | -------------------------------------- |
| **业务 DTO** | `src/generated/models/` | `UserUserWithRolesDTO`, `AuthLoginDTO` |
| **响应包装** | `src/api/types.ts`      | `ApiResponse<T>`, `ListApiResponse<T>` |
| **前端派生** | `src/api/types.ts`      | `Menu.children`（树形结构）            |
| **前端状态** | `src/api/types.ts`      | `PaginationState`, `LoginResult`       |

## 禁止事项

```typescript
// ❌ 禁止：在 pages/components 中重复定义 DTO
interface UserDTO {
  id: number;
  username: string;
}

// ❌ 禁止：前端扩展补救后端缺陷
export interface LoginResponse extends AuthLoginResultDTO {
  user?: User; // 后端返回但 Swagger 未定义 → 修复后端！
}

// ❌ 禁止：在前端临时定义类型"修复"问题
```

## 允许的扩展

**唯一允许 `extends` 的场景**：前端根据已有数据派生的字段（后端不需要返回）。

```typescript
// ✅ 允许：前端构建的树形结构
export interface Menu extends MenuMenuDTO {
  children?: Menu[]; // 前端根据 parent_id 构建，非后端返回
}
```

**判断标准**：字段数据来源是什么？

- 后端返回但未定义 → **修复后端**
- 前端根据已有数据计算/派生 → **允许 extends**

## 设计反思原则

**前端类型困难是后端设计问题的信号。**

| 症状                  | 可能的设计问题                     |
| --------------------- | ---------------------------------- |
| 需要在前端定义 DTO    | 后端 DTO 未暴露或 Swagger 注解缺失 |
| 需要 `extends` 补字段 | 后端返回了未定义的字段             |
| 类型与实际响应不匹配  | Handler 注解与实际返回不一致       |
| 需要 `as any` 绕过    | API 设计存在根本问题               |

**正确的修复方向**：

1. Handler `@Success` 注解是否使用正确的 DTO 类型？
2. DTO 字段是否都有 `json` tag？
3. DTO 是否在 `internal/application/*/dto.go` 定义？

**修复流程**：

```
后端 DTO 补全
    ↓
pre-commit run swagger-generate
    ↓
pnpm api:generate
    ↓
前端直接使用生成类型
```

**验证**：`go test ./internal/...` + `pnpm vue-tsc --noEmit`

## 目录结构

```
src/api/
├── auth/              # 认证相关
│   ├── client.ts      # axios 实例 + API 实例化
│   └── platformAuth.ts # 自定义认证 API
├── types.ts           # 响应包装 + 前端派生/状态类型
├── helpers.ts         # 响应提取辅助函数
├── errors.ts          # 统一错误处理
└── index.ts           # 统一导出
```
