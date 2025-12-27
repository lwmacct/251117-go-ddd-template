---
paths:
  - "src/pages/**/*.vue"
  - "src/pages/**/*.ts"
---

# Pages 目录规范

## 核心概念

`src/pages/` 采用**基于文件系统的路由**，目录结构即路由结构。

```
src/pages/
├── admin/                    # /admin（分级目录，无 index.vue）
│   ├── users/               # /admin/users
│   │   └── index.vue        # ← 页面组件
│   └── roles/               # /admin/roles
│       └── index.vue
├── auth/                     # /auth（分级目录）
│   ├── login/               # /auth/login
│   │   └── index.vue
│   └── register/            # /auth/register
│       └── index.vue
└── user/                     # /user（分级目录）
    └── profile/             # /user/profile
        └── index.vue
```

**判断规则**：

- 有 `index.vue` → **页面**（渲染内容）
- 无 `index.vue` → **分级目录**（仅组织路由层级）

## 页面目录结构

```
pages/{module}/{page}/
├── index.vue              # 页面主组件（必须）
├── composables/           # 页面状态（可选）
│   └── use{Feature}.ts
└── components/            # 页面子组件（可选）
    └── {Name}.vue
```

## 职责划分

| 层级           | 职责                         | API 调用 |
| -------------- | ---------------------------- | -------- |
| `index.vue`    | 编排布局、协调对话框状态     | ❌       |
| `composables/` | 业务逻辑、状态管理           | ✅       |
| `components/`  | UI 交互、表单验证、emit 事件 | ❌       |

## 禁止事项

- ❌ 在 `index.vue` 或子组件中直接调用 API
- ❌ 在 pages 中定义 DTO（使用 `@/generated/models`）
