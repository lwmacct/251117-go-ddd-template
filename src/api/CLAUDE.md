# API 层开发规范

本目录是前端与后端 API 的桥梁，所有类型定义应来自 OpenAPI 生成代码。

<!--TOC-->

## Table of Contents

- [类型使用原则](#类型使用原则) `:22+59`
  - [正确做法](#正确做法) `:26+15`
  - [禁止做法](#禁止做法) `:41+10`
  - [类型来源](#类型来源) `:51+12`
  - [扩展类型的正确方式](#扩展类型的正确方式) `:63+18`
- [当需要新类型时](#当需要新类型时) `:81+32`
  - [检查清单](#检查清单) `:85+14`
  - [修复流程](#修复流程) `:99+14`
- [目录结构](#目录结构) `:113+13`
- [反思：前端扩展类型是后端缺陷的信号](#反思前端扩展类型是后端缺陷的信号) `:126+12`

<!--TOC-->

## 类型使用原则

**核心原则：前端不定义与后端重复的 DTO，只消费 OpenAPI 生成的类型。**

### 正确做法

```typescript
// 从 @/api 统一导入类型
import type { User, LoginRequest, CaptchaData } from "@/api";

// 使用 API 类型定义组件 Props
interface Props {
  user: User;
}

// 使用 API 类型定义表单数据
const formData = ref<LoginRequest>({ ... });
```

### 禁止做法

```typescript
// 禁止在 pages/components 中重复定义 DTO
interface UserDTO { id: number; username: string; }  // ❌

// 禁止在 src/types/ 中定义与 API 重复的类型
export interface LoginRequest { ... }  // ❌
```

### 类型来源

| 类型分类     | 来源                         | 示例                                       |
| ------------ | ---------------------------- | ------------------------------------------ |
| **业务 DTO** | `src/api/generated/models/`  | `UserUserWithRolesDTO`, `AuthLoginDTO`     |
| **简化别名** | `src/api/types.ts`           | `User = UserUserWithRolesDTO`              |
| **扩展类型** | `src/api/types.ts` (extends) | `LoginResponse extends AuthLoginResultDTO` |
| **响应包装** | `src/types/response/`        | `ApiResponse<T>` (通用结构)                |
| **前端状态** | `src/api/types.ts` 或组件内  | `PaginationState`, `LoginResult`           |

> **注意**：前端状态类型（如 `LoginResult`）虽然字段与 API 相关，但表示的是前端操作结果，不是后端响应结构。

### 扩展类型的正确方式

当后端 DTO 缺少前端需要的字段时，**临时**使用 `extends` 扩展：

```typescript
// ⚠️ 临时方案：扩展生成的类型（应尽快修复后端）
export interface LoginResponse extends AuthLoginResultDTO {
  user?: User; // 后端实际返回但 Swagger 未定义
}

// ✅ 合理扩展：添加前端特有的派生字段
export interface Menu extends MenuMenuDTO {
  children?: Menu[]; // 树形结构，前端根据 parent_id 构建
}
```

> **重要**：`extends` 扩展是临时补救措施，不是最终方案。发现需要扩展时，应立即记录并安排后端修复。

## 当需要新类型时

**如果发现需要在前端定义新的 DTO，首先反思后端是否有问题：**

### 检查清单

1. **Swagger 注解是否完整？**
   - Handler 的 `@Success` 注解是否使用了正确的 DTO 类型？
   - DTO 字段是否都有 `json` tag？

2. **API 响应格式是否规范？**
   - 是否使用 `response.OK()` / `response.List()` 统一封装？
   - 响应结构是否与 Swagger 定义一致？

3. **DTO 是否在 Application 层定义？**
   - 所有 DTO 应在 `internal/application/*/dto.go` 中定义
   - Handler 层不应定义业务 DTO

### 修复流程

```
发现前端需要新类型
    ↓
检查后端 Swagger 注解 → 补充缺失的注解
    ↓
运行 swag init + pnpm api:generate
    ↓
前端使用生成的类型
```

**禁止绕过：在前端临时定义类型来"修复"问题。**

## 目录结构

```
src/api/
├── auth/              # 认证相关
│   ├── client.ts      # axios 实例 + API 实例化
│   └── platformAuth.ts # 自定义认证 API
├── types.ts           # 类型别名 + 扩展（唯一允许手写类型的地方）
├── helpers.ts         # 响应提取辅助函数
├── errors.ts          # 统一错误处理
└── index.ts           # 统一导出
```

## 反思：前端扩展类型是后端缺陷的信号

当发现需要 `interface X extends GeneratedDTO { missingField }` 时，应质疑后端：

| 前端症状     | 后端根因                   |
| ------------ | -------------------------- |
| 扩展添加字段 | DTO 字段缺失               |
| 扩展覆盖类型 | Swagger 注解使用错误的 DTO |

**修复流程**：后端 DTO 补全 → `pre-commit run swagger-generate` → 前端扩展改为别名

**验证**：`go test ./internal/...` + `pnpm vue-tsc --noEmit`
