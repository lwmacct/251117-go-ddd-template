<script setup lang="ts">
import { ref } from "vue";

/**
 * 角色管理页面
 * 用于查看和管理角色权限
 */

interface RoleItem {
  name: string;
  code: string;
  description: string;
  status: string;
  createdAt: string;
}

const roles = ref<RoleItem[]>([]);
const dialog = ref(false);
const search = ref("");

// 表头配置
const headers = [
  {
    title: "角色名称",
    key: "name",
  },
  {
    title: "角色代码",
    key: "code",
  },
  {
    title: "描述",
    key: "description",
  },
  {
    title: "状态",
    key: "status",
  },
  {
    title: "创建时间",
    key: "createdAt",
  },
  {
    title: "操作",
    key: "actions",
    sortable: false,
  },
];

// 添加/编辑角色
const openDialog = () => {
  dialog.value = true;
};
</script>

<template>
  <div class="roles-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">角色管理</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索角色" single-line hide-details variant="outlined" density="compact"></v-text-field>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="openDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建角色
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text>
            <v-data-table :headers="headers" :items="roles" :search="search" no-data-text="暂无角色数据">
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
        <v-card-title>新建角色</v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field label="角色名称" variant="outlined" class="mb-4"></v-text-field>
            <v-text-field label="角色代码" variant="outlined" class="mb-4"></v-text-field>
            <v-textarea label="描述" variant="outlined" rows="3"></v-textarea>
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
.roles-page {
  width: 100%;
}
</style>
