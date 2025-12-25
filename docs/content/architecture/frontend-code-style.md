---
title: 前端代码规范
description: Vue 3 + TypeScript 前端代码规范和最佳实践
outline: [2, 3]
---

# 前端代码规范

本文档定义了前端项目的代码规范和最佳实践，确保代码质量和团队协作效率。

<!--TOC-->

## Table of Contents

- [工具链概览](#工具链概览) `:45+11`
- [ESLint 配置](#eslint-配置) `:56+41`
  - [常用命令](#常用命令) `:87+10`
- [Prettier 配置](#prettier-配置) `:97+27`
  - [常用命令](#常用命令-1) `:114+10`
- [组件编写规范](#组件编写规范) `:124+85`
  - [文件命名](#文件命名) `:126+12`
  - [组件结构](#组件结构) `:138+50`
  - [Props 定义](#props-定义) `:188+21`
- [Composables 规范](#composables-规范) `:209+54`
  - [命名规范](#命名规范) `:211+5`
  - [结构模板](#结构模板) `:216+47`
- [TypeScript 规范](#typescript-规范) `:263+42`
  - [类型定义](#类型定义) `:265+24`
  - [未使用变量](#未使用变量) `:289+16`
- [测试规范](#测试规范) `:305+51`
  - [文件位置](#文件位置) `:307+13`
  - [测试结构](#测试结构) `:320+23`
  - [常用命令](#常用命令-2) `:343+13`
- [Git 提交规范](#git-提交规范) `:356+36`
  - [Commit Message 格式](#commit-message-格式) `:358+10`
  - [类型说明](#类型说明) `:368+12`
  - [示例](#示例) `:380+12`
- [IDE 配置](#ide-配置) `:392+23`
  - [VS Code 推荐扩展](#vs-code-推荐扩展) `:394+7`
  - [推荐设置](#推荐设置) `:401+14`
- [下一步](#下一步) `:415+5`

<!--TOC-->

## 工具链概览

项目使用以下工具保证代码质量：

| 工具     | 版本 | 用途                |
| -------- | ---- | ------------------- |
| ESLint   | 9.x  | 代码静态分析        |
| Prettier | 3.x  | 代码格式化          |
| Vitest   | 4.x  | 单元测试            |
| vue-tsc  | 3.x  | TypeScript 类型检查 |

## ESLint 配置

项目使用 ESLint 9.x 的 **Flat Config** 格式（`eslint.config.js`）：

```javascript
// eslint.config.js 核心配置
export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs["flat/recommended"],
  prettier,
  {
    rules: {
      // Vue 规则
      "vue/multi-word-component-names": "off",
      "vue/valid-v-slot": ["error", { allowModifiers: true }],

      // TypeScript 规则
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": [
        "warn",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
    },
  },
);
```

### 常用命令

```bash
# 检查代码
npm run lint

# 自动修复
npm run lint:fix
```

## Prettier 配置

格式化规则定义在 `.prettierrc`：

```json
{
  "semi": true,
  "singleQuote": false,
  "tabWidth": 2,
  "trailingComma": "es5",
  "printWidth": 120,
  "bracketSpacing": true,
  "arrowParens": "always",
  "endOfLine": "lf"
}
```

### 常用命令

```bash
# 格式化代码
npm run format

# 检查格式
npm run format:check
```

## 组件编写规范

### 文件命名

```
✅ 推荐
src/components/UserAvatar.vue
src/pages/admin/users/components/UserDialog.vue

❌ 避免
src/components/useravatar.vue
src/components/user_avatar.vue
```

### 组件结构

使用 `<script setup>` 语法，按以下顺序组织代码：

```vue
<script setup lang="ts">
// 1. 类型导入
import type { User } from "@/types";

// 2. 组件导入
import UserAvatar from "@/components/UserAvatar.vue";

// 3. Composables
const { users, loading, fetchUsers } = useUsers();

// 4. Props & Emits
const props = defineProps<{
  userId: number;
}>();

const emit = defineEmits<{
  (e: "update", user: User): void;
}>();

// 5. 响应式状态
const isEditing = ref(false);

// 6. 计算属性
const displayName = computed(() => `${props.user.firstName} ${props.user.lastName}`);

// 7. 方法
const handleSubmit = async () => {
  // ...
};

// 8. 生命周期
onMounted(() => {
  fetchUsers();
});
</script>

<template>
  <!-- 模板内容 -->
</template>

<style scoped>
/* 组件样式 */
</style>
```

### Props 定义

使用 TypeScript 类型定义：

```typescript
// ✅ 推荐：使用 TypeScript 接口
interface Props {
  user: User;
  readonly?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  readonly: false,
});

// ❌ 避免：运行时声明
const props = defineProps({
  user: { type: Object, required: true },
});
```

## Composables 规范

### 命名规范

- 文件名：`use{Feature}.ts`（如 `useClipboard.ts`）
- 函数名：`use{Feature}`（如 `useClipboard`）

### 结构模板

```typescript
// src/composables/useExample.ts
import { ref, computed, onUnmounted } from "vue";

export interface UseExampleOptions {
  initialValue?: string;
}

export interface UseExampleReturn {
  value: Ref<string>;
  isLoading: Ref<boolean>;
  update: (newValue: string) => Promise<void>;
}

export function useExample(options: UseExampleOptions = {}): UseExampleReturn {
  const { initialValue = "" } = options;

  // 状态
  const value = ref(initialValue);
  const isLoading = ref(false);

  // 方法
  const update = async (newValue: string) => {
    isLoading.value = true;
    try {
      // 业务逻辑
      value.value = newValue;
    } finally {
      isLoading.value = false;
    }
  };

  // 清理
  onUnmounted(() => {
    // 清理资源
  });

  return {
    value,
    isLoading,
    update,
  };
}
```

## TypeScript 规范

### 类型定义

```typescript
// ✅ 推荐：明确的类型定义
interface User {
  id: number;
  name: string;
  email: string;
  createdAt: Date;
}

// ✅ 推荐：使用 type 定义联合类型
type UserStatus = "active" | "inactive" | "pending";

// ❌ 避免：使用 any
const data: any = fetchData();

// ✅ 推荐：使用 unknown 并进行类型守卫
const data: unknown = fetchData();
if (isUser(data)) {
  console.log(data.name);
}
```

### 未使用变量

如果参数有意不使用，使用下划线前缀：

```typescript
// ✅ 正确
const handleClick = (_event: MouseEvent) => {
  // 不需要使用 event
};

// ❌ 错误（ESLint 会警告）
const handleClick = (event: MouseEvent) => {
  // event 未使用
};
```

## 测试规范

### 文件位置

测试文件放在被测试文件的 `__tests__` 目录下：

```
src/composables/
├── useClipboard.ts
├── useDebounce.ts
└── __tests__/
    ├── useClipboard.spec.ts
    └── useDebounce.spec.ts
```

### 测试结构

```typescript
import { describe, it, expect, vi, beforeEach } from "vitest";

describe("useClipboard", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should initialize with default values", () => {
    const { copied, error } = useClipboard();

    expect(copied.value).toBe(false);
    expect(error.value).toBeNull();
  });

  it("should copy text successfully", async () => {
    // 测试逻辑
  });
});
```

### 常用命令

```bash
# 运行测试（watch 模式）
npm run test

# 单次运行
npm run test:run

# 生成覆盖率报告
npm run test:coverage
```

## Git 提交规范

### Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 类型说明

| Type       | 说明                   |
| ---------- | ---------------------- |
| `feat`     | 新功能                 |
| `fix`      | Bug 修复               |
| `docs`     | 文档更新               |
| `style`    | 代码格式（不影响功能） |
| `refactor` | 重构                   |
| `test`     | 测试相关               |
| `chore`    | 构建/工具相关          |

### 示例

```bash
feat(web): add user avatar upload component

- Add AvatarUploader component with drag & drop support
- Implement image compression before upload
- Add preview functionality

Closes #123
```

## IDE 配置

### VS Code 推荐扩展

- Vue - Official（Volar）
- ESLint
- Prettier - Code formatter
- TypeScript Vue Plugin

### 推荐设置

```json
// .vscode/settings.json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": "explicit"
  },
  "typescript.preferences.importModuleSpecifier": "relative"
}
```

## 下一步

- [前端架构概述](./frontend-overview) - 了解整体架构
- [组件开发](./frontend-components) - 组件开发指南
- [状态管理](./frontend-state) - Pinia 使用指南
