---
paths:
  - "src/**/*.vue"
  - "src/**/*.ts"
  - "src/**/*.tsx"
---

# 检查命令规范

前端编码完成后，必须执行以下检查确保代码质量。

## 类型检查

**命令**：`pnpm vue-tsc --noEmit`

**说明**：使用 Vue TypeScript 编译器进行类型检查，不生成输出文件。

**通过标准**：无 TypeScript 错误输出。

```bash
pnpm vue-tsc --noEmit
```

## ESLint 检查

**命令**：`pnpm eslint`

**说明**：检查代码风格和规范问题。

**通过标准**：

- 无 `error` 级别问题
- 警告（`warning`）数量应尽量少

```bash
# 检查所有文件
pnpm eslint

# 检查特定目录
pnpm eslint src/pages/

# 自动修复格式问题
pnpm eslint --fix
```

## ESLint 规则优先级

| 规则类别    | 优先级 | 说明                                 |
| ----------- | ------ | ------------------------------------ |
| 类型错误    | P0     | 必须修复，如未定义变量、类型不匹配   |
| Vue 规范    | P0     | Vue 3 特定规则，如模板语法、组件结构 |
| Import 规范 | P1     | 导入顺序和路径                       |
| 格式规范    | P2     | Prettier 格式问题，可自动修复        |

## 提交前检查清单

在提交代码或创建 PR 前，确认以下项：

- [ ] `pnpm vue-tsc --noEmit` 通过（无错误）
- [ ] `pnpm eslint` 无 error 级别问题
- [ ] 如有格式问题，运行 `pnpm eslint --fix` 自动修复
- [ ] 新增页面已添加路由配置和菜单项
- [ ] 新增 API 已在 `src/api/index.ts` 中导出
- [ ] DTO 类型来自 `@models` 而非前端定义
- [ ] Vuetify 组件使用主题系统而非硬编码颜色

## 常见问题修复

### 类型错误

```bash
# 查看具体错误
pnpm vue-tsc --noEmit 2>&1 | grep error
```

### 格式问题

```bash
# 自动修复大部分格式问题
pnpm eslint --fix

# 修复后再次检查
pnpm eslint
```

### 未使用的导入

```typescript
// ❌ 错误：导入但未使用
import { extractData, extractList } from "@/api";

// ✅ 正确：只导入需要的
import { extractList } from "@/api";
```
