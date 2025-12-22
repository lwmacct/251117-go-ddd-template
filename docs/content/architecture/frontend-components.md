# 组件开发

前端采用基于页面的组件组织方式，每个页面包含独立的组件、Composables 和类型定义。

<!--TOC-->

## Table of Contents

- [页面结构](#页面结构) `:29+14`
- [页面模块](#页面模块) `:43+71`
  - [Auth 模块](#auth-模块) `:45+20`
  - [Admin 模块](#admin-模块) `:65+27`
  - [User 模块](#user-模块) `:92+22`
- [Composable 模式](#composable-模式) `:114+95`
- [页面组件示例](#页面组件示例) `:209+54`
- [对话框组件示例](#对话框组件示例) `:263+98`
- [布局组件](#布局组件) `:361+31`
  - [AdminLayout](#adminlayout) `:363+29`
- [共享视图组件](#共享视图组件) `:392+35`
  - [导航菜单](#导航菜单) `:394+33`
- [最佳实践](#最佳实践) `:427+57`
  - [1. 组件职责单一](#1-组件职责单一) `:429+6`
  - [2. Props/Events 明确](#2-propsevents-明确) `:435+18`
  - [3. v-model 双向绑定](#3-v-model-双向绑定) `:453+13`
  - [4. 异步状态处理](#4-异步状态处理) `:466+18`

<!--TOC-->

## 页面结构

```
src/pages/{module}/{feature}/
├── index.vue                  # 页面入口组件
├── components/                # 页面私有组件
│   ├── SomeDialog.vue
│   └── SomeForm.vue
├── composables/               # 页面逻辑
│   ├── index.ts
│   └── useFeature.ts
└── types.ts                   # 页面类型（可选）
```

## 页面模块

### Auth 模块

```
pages/auth/
├── login/
│   ├── index.vue              # 登录页（管理登录和 2FA 状态）
│   ├── components/
│   │   ├── LoginForm.vue      # 登录表单
│   │   └── TwoFactorForm.vue  # 2FA 验证表单
│   └── composables/
│       └── useLogin.ts        # 登录逻辑
└── register/
    ├── index.vue
    ├── components/
    │   ├── RegisterForm.vue
    │   └── VerifyEmailForm.vue
    └── composables/
        └── useRegister.ts
```

### Admin 模块

```
pages/admin/
├── overview/                  # 数据概览
├── users/                     # 用户管理
│   ├── components/
│   │   ├── UserDialog.vue     # 用户编辑对话框
│   │   └── RoleSelector.vue   # 角色选择器
│   └── composables/
│       └── useAdminUsers.ts
├── roles/                     # 角色管理
│   ├── components/
│   │   ├── RoleDialog.vue
│   │   └── PermissionSelector.vue
│   └── composables/
│       └── useRoles.ts
├── menus/                     # 菜单管理
│   ├── components/
│   │   ├── MenuTree.vue       # 树形菜单
│   │   └── MenuDialog.vue
│   └── composables/
│       └── useMenus.ts
├── settings/                  # 系统设置
└── auditlogs/                 # 审计日志
```

### User 模块

```
pages/user/
├── profile/                   # 个人资料
│   ├── components/
│   │   └── BasicInfoForm.vue
│   └── composables/
├── security/                  # 安全设置
│   ├── components/
│   │   ├── PasswordSettings.vue
│   │   └── TwoFactorSettings.vue
│   └── composables/
│       └── useTwoFactor.ts
└── tokens/                    # 访问令牌
    ├── components/
    │   ├── TokenDialog.vue
    │   └── TokenDisplay.vue
    └── composables/
        └── useTokens.ts
```

## Composable 模式

每个页面的业务逻辑封装在 Composable 中：

```typescript
// src/pages/admin/users/composables/useAdminUsers.ts
import { ref, reactive, computed, onMounted } from "vue";
import { listUsers, createUser, updateUser, deleteUser } from "@/api/admin/users";
import type { AdminUser, CreateUserRequest, UpdateUserRequest } from "@/types/admin";

export function useAdminUsers() {
  // ========== 状态 ==========
  const users = ref<AdminUser[]>([]);
  const loading = ref(false);
  const dialogVisible = ref(false);
  const editingUser = ref<AdminUser | null>(null);

  const pagination = reactive({
    page: 1,
    limit: 10,
    total: 0,
  });

  // ========== 计算属性 ==========
  const isEditing = computed(() => !!editingUser.value);
  const dialogTitle = computed(() => (isEditing.value ? "编辑用户" : "创建用户"));

  // ========== 方法 ==========
  async function fetchUsers() {
    loading.value = true;
    try {
      const response = await listUsers({
        page: pagination.page,
        limit: pagination.limit,
      });
      users.value = response.data;
      pagination.total = response.pagination.total;
    } catch (error) {
      console.error("Failed to fetch users:", error);
    } finally {
      loading.value = false;
    }
  }

  function openCreateDialog() {
    editingUser.value = null;
    dialogVisible.value = true;
  }

  function openEditDialog(user: AdminUser) {
    editingUser.value = { ...user };
    dialogVisible.value = true;
  }

  async function handleSubmit(data: CreateUserRequest | UpdateUserRequest) {
    if (isEditing.value) {
      await updateUser(editingUser.value!.id, data as UpdateUserRequest);
    } else {
      await createUser(data as CreateUserRequest);
    }
    dialogVisible.value = false;
    await fetchUsers();
  }

  async function handleDelete(id: number) {
    await deleteUser(id);
    await fetchUsers();
  }

  // ========== 生命周期 ==========
  onMounted(() => {
    fetchUsers();
  });

  // ========== 导出 ==========
  return {
    // 状态
    users,
    loading,
    dialogVisible,
    editingUser,
    pagination,
    // 计算属性
    isEditing,
    dialogTitle,
    // 方法
    fetchUsers,
    openCreateDialog,
    openEditDialog,
    handleSubmit,
    handleDelete,
  };
}
```

## 页面组件示例

```vue
<!-- src/pages/admin/users/index.vue -->
<script setup lang="ts">
import { useAdminUsers } from "./composables/useAdminUsers";
import UserDialog from "./components/UserDialog.vue";

const { users, loading, dialogVisible, editingUser, pagination, dialogTitle, openCreateDialog, openEditDialog, handleSubmit, handleDelete } = useAdminUsers();

const headers = [
  { title: "ID", key: "id" },
  { title: "用户名", key: "username" },
  { title: "邮箱", key: "email" },
  { title: "状态", key: "status" },
  { title: "操作", key: "actions", sortable: false },
];
</script>

<template>
  <v-container>
    <v-card>
      <v-card-title class="d-flex align-center">
        <span>用户管理</span>
        <v-spacer />
        <v-btn color="primary" @click="openCreateDialog">
          <v-icon left>mdi-plus</v-icon>
          创建用户
        </v-btn>
      </v-card-title>

      <v-data-table :headers="headers" :items="users" :loading="loading" :items-per-page="pagination.limit">
        <template #item.status="{ item }">
          <v-chip :color="item.status === 'active' ? 'success' : 'error'">
            {{ item.status }}
          </v-chip>
        </template>

        <template #item.actions="{ item }">
          <v-btn icon size="small" @click="openEditDialog(item)">
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
          <v-btn icon size="small" color="error" @click="handleDelete(item.id)">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </template>
      </v-data-table>
    </v-card>

    <UserDialog v-model="dialogVisible" :title="dialogTitle" :user="editingUser" @submit="handleSubmit" />
  </v-container>
</template>
```

## 对话框组件示例

```vue
<!-- src/pages/admin/users/components/UserDialog.vue -->
<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { AdminUser, CreateUserRequest } from "@/types/admin";

const props = defineProps<{
  modelValue: boolean;
  title: string;
  user: AdminUser | null;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  submit: [data: CreateUserRequest];
}>();

const form = ref({
  username: "",
  email: "",
  password: "",
  full_name: "",
  status: "active" as const,
});

const isEditing = computed(() => !!props.user);

// 编辑时填充表单
watch(
  () => props.user,
  (user) => {
    if (user) {
      form.value = {
        username: user.username,
        email: user.email,
        password: "",
        full_name: user.full_name ?? "",
        status: user.status,
      };
    } else {
      // 重置表单
      form.value = {
        username: "",
        email: "",
        password: "",
        full_name: "",
        status: "active",
      };
    }
  },
  { immediate: true },
);

function handleSubmit() {
  emit("submit", form.value);
}

function handleClose() {
  emit("update:modelValue", false);
}
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="500" @update:model-value="handleClose">
    <v-card>
      <v-card-title>{{ title }}</v-card-title>

      <v-card-text>
        <v-form @submit.prevent="handleSubmit">
          <v-text-field v-model="form.username" label="用户名" :disabled="isEditing" required />
          <v-text-field v-model="form.email" label="邮箱" type="email" required />
          <v-text-field v-if="!isEditing" v-model="form.password" label="密码" type="password" required />
          <v-text-field v-model="form.full_name" label="姓名" />
          <v-select
            v-model="form.status"
            label="状态"
            :items="[
              { title: '激活', value: 'active' },
              { title: '禁用', value: 'inactive' },
            ]"
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn @click="handleClose">取消</v-btn>
        <v-btn color="primary" @click="handleSubmit">
          {{ isEditing ? "保存" : "创建" }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
```

## 布局组件

### AdminLayout

```vue
<!-- src/layout/AdminLayout.vue -->
<script setup lang="ts">
import AppBars from "@/views/AppBars/index.vue";
import Navigation from "@/views/Navigation/index.vue";

const menuItems = [
  { title: "数据概览", path: "/admin/overview", icon: "mdi-speedometer" },
  { title: "用户管理", path: "/admin/users", icon: "mdi-account" },
  { title: "角色管理", path: "/admin/roles", icon: "mdi-account-group" },
  { title: "菜单管理", path: "/admin/menus", icon: "mdi-menu" },
  { title: "系统设置", path: "/admin/settings", icon: "mdi-cog" },
  { title: "审计日志", path: "/admin/auditlogs", icon: "mdi-file-document-outline" },
];
</script>

<template>
  <v-app>
    <AppBars />
    <Navigation :items="menuItems" />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>
```

## 共享视图组件

### 导航菜单

```vue
<!-- src/views/Navigation/index.vue -->
<script setup lang="ts">
import { useRoute } from "vue-router";
import type { MenuItem } from "./types";

defineProps<{
  items: MenuItem[];
}>();

const route = useRoute();

function isActive(path: string) {
  return route.path === path;
}
</script>

<template>
  <v-navigation-drawer permanent>
    <v-list nav>
      <v-list-item v-for="item in items" :key="item.path" :to="item.path" :active="isActive(item.path)">
        <template #prepend>
          <v-icon>{{ item.icon }}</v-icon>
        </template>
        <v-list-item-title>{{ item.title }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>
```

## 最佳实践

### 1. 组件职责单一

- 页面组件负责布局和组装
- 子组件负责特定功能
- Composable 负责业务逻辑

### 2. Props/Events 明确

```vue
<script setup lang="ts">
// 明确定义 Props 类型
defineProps<{
  user: AdminUser | null;
  loading: boolean;
}>();

// 明确定义 Events 类型
defineEmits<{
  submit: [data: CreateUserRequest];
  cancel: [];
}>();
</script>
```

### 3. v-model 双向绑定

```vue
<!-- 父组件 -->
<UserDialog v-model="dialogVisible" />

<!-- 子组件 -->
<script setup>
defineProps<{ modelValue: boolean }>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>
```

### 4. 异步状态处理

```typescript
const loading = ref(false);
const error = ref<string | null>(null);

async function fetchData() {
  loading.value = true;
  error.value = null;
  try {
    // ...
  } catch (e) {
    error.value = formatError(e);
  } finally {
    loading.value = false;
  }
}
```
