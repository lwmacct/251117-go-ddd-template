<script setup lang="ts">
import { ref } from "vue";

/**
 * 菜单管理页面
 * 用于查看和管理系统菜单
 */

interface MenuItem {
  id: number;
  title: string;
  path: string;
  icon: string;
  parent?: number;
  order: number;
  visible: boolean;
}

const menus = ref<MenuItem[]>([]);
const dialog = ref(false);
const search = ref("");
const treeView = ref(true);

// 表头配置
const headers = [
  { title: "菜单名称", key: "title" },
  { title: "路径", key: "path" },
  { title: "图标", key: "icon" },
  { title: "排序", key: "order" },
  { title: "可见", key: "visible" },
  { title: "操作", key: "actions", sortable: false },
];

// 添加/编辑菜单
const openDialog = () => {
  dialog.value = true;
};

// 切换视图模式
const toggleView = () => {
  treeView.value = !treeView.value;
};
</script>

<template>
  <div class="menus-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">菜单管理</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="4">
                <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索菜单" single-line hide-details variant="outlined" density="compact"></v-text-field>
              </v-col>
              <v-col cols="12" md="8" class="text-right">
                <v-btn-toggle v-model="treeView" mandatory class="mr-4">
                  <v-btn :value="true" icon="mdi-file-tree"></v-btn>
                  <v-btn :value="false" icon="mdi-table"></v-btn>
                </v-btn-toggle>
                <v-btn color="primary" @click="openDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建菜单
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text>
            <!-- 树形视图 -->
            <div v-if="treeView">
              <v-alert type="info" variant="tonal"> 树形视图功能开发中... </v-alert>
            </div>

            <!-- 表格视图 -->
            <v-data-table v-else :headers="headers" :items="menus" :search="search" no-data-text="暂无菜单数据">
              <template #item.icon="{ item }">
                <v-icon>{{ item.icon }}</v-icon>
              </template>
              <template #item.visible="{ item }">
                <v-chip :color="item.visible ? 'success' : 'error'" size="small">
                  {{ item.visible ? "显示" : "隐藏" }}
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
        <v-card-title>新建菜单</v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field label="菜单名称" variant="outlined" class="mb-4"></v-text-field>
            <v-text-field label="路径" variant="outlined" class="mb-4"></v-text-field>
            <v-text-field label="图标 (MDI) " placeholder="mdi-home" variant="outlined" class="mb-4"></v-text-field>
            <v-select label="父级菜单" :items="['无', '系统管理', '用户中心']" variant="outlined" class="mb-4"></v-select>
            <v-text-field label="排序" type="number" variant="outlined" class="mb-4"></v-text-field>
            <v-switch label="是否可见" color="primary"></v-switch>
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
.menus-page {
  width: 100%;
}
</style>
