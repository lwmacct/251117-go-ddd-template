<script setup lang="ts">
import { ref } from "vue";

/**
 * 用户管理页面
 * 用于查看和管理系统用户
 */

interface UserItem {
  username: string;
  email: string;
  role: string;
  status: string;
  lastLogin: string;
}

const users = ref<UserItem[]>([]);
const dialog = ref(false);
const search = ref("");

// 表头配置
const headers = [
  {
    title: "用户名",
    key: "username",
  },
  {
    title: "邮箱",
    key: "email",
  },
  {
    title: "角色",
    key: "role",
  },
  {
    title: "状态",
    key: "status",
  },
  {
    title: "最后登录",
    key: "lastLogin",
  },
  {
    title: "操作",
    key: "actions",
    sortable: false,
  },
];

// 添加/编辑用户
const openDialog = () => {
  dialog.value = true;
};
</script>

<template>
  <div class="users-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">用户管理</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索用户" single-line hide-details variant="outlined" density="compact"></v-text-field>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="openDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建用户
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text>
            <v-data-table :headers="headers" :items="users" :search="search" no-data-text="暂无用户数据">
              <template #item.status="{ item }">
                <v-chip :color="item.status === 'active' ? 'success' : 'error'" size="small">
                  {{ item.status === "active" ? "启用" : "禁用" }}
                </v-chip>
              </template>
              <template #item.actions>
                <v-btn icon="mdi-pencil" size="small" variant="text"></v-btn>
                <v-btn icon="mdi-delete" size="small" variant="text" color="error"></v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 新建/编辑对话框 -->
    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title>新建用户</v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field label="用户名" variant="outlined" class="mb-4"></v-text-field>
            <v-text-field label="邮箱" type="email" variant="outlined" class="mb-4"></v-text-field>
            <v-text-field label="密码" type="password" variant="outlined" class="mb-4"></v-text-field>
            <v-select label="角色" :items="['管理员', '普通用户', '访客']" variant="outlined"></v-select>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="dialog = false">取消</v-btn>
          <v-btn color="primary" @click="dialog = false">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.users-page {
  width: 100%;
}
</style>
