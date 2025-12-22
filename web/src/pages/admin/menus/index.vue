<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useMenus } from "./composables/useMenus";
import MenuDialog from "./components/MenuDialog.vue";
import MenuTree from "./components/MenuTree.vue";
import type { Menu, CreateMenuRequest, UpdateMenuRequest } from "@/types/admin";

const {
  menus,
  loading,
  errorMessage,
  successMessage,
  fetchMenus,
  createMenu,
  updateMenu,
  deleteMenu,
  reorderMenus,
  clearMessages,
  exportMenus,
} = useMenus();

const menuDialog = ref(false);
const deleteDialog = ref(false);
const dialogMode = ref<"create" | "edit">("create");
const selectedMenu = ref<Menu | null>(null);
const menuToDelete = ref<Menu | null>(null);
const viewMode = ref<"tree" | "table">("tree");

onMounted(() => {
  fetchMenus();
});

const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedMenu.value = null;
  menuDialog.value = true;
};

const openEditDialog = (menu: Menu) => {
  dialogMode.value = "edit";
  selectedMenu.value = menu;
  menuDialog.value = true;
};

const openDeleteDialog = (menu: Menu) => {
  menuToDelete.value = menu;
  deleteDialog.value = true;
};

const handleSaveMenu = async (data: CreateMenuRequest | UpdateMenuRequest) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createMenu(data as CreateMenuRequest);
  } else if (selectedMenu.value) {
    success = await updateMenu(selectedMenu.value.id, data as UpdateMenuRequest);
  }

  if (success) {
    menuDialog.value = false;
  }
};

const confirmDelete = async () => {
  if (!menuToDelete.value) return;

  const success = await deleteMenu(menuToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    menuToDelete.value = null;
  }
};

// 处理拖拽排序
const handleMenusReorder = async (updatedMenus: Menu[]) => {
  // 构建排序数据
  const buildReorderData = (menus: Menu[], parentId?: number) => {
    const result: Array<{ id: number; order: number; parent_id?: number }> = [];
    menus.forEach((menu, index) => {
      result.push({
        id: menu.id,
        order: index,
        parent_id: parentId,
      });
      if (menu.children && menu.children.length > 0) {
        result.push(...buildReorderData(menu.children, menu.id));
      }
    });
    return result;
  };

  const reorderData = {
    menus: buildReorderData(updatedMenus),
  };

  await reorderMenus(reorderData);
  await fetchMenus(); // 刷新列表
};
</script>

<template>
  <div class="menus-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">菜单管理</h1>
      </v-col>
    </v-row>

    <v-row v-if="errorMessage || successMessage">
      <v-col cols="12">
        <v-alert v-if="errorMessage" type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
        </v-alert>
        <v-alert v-if="successMessage" type="success" closable @click:close="clearMessages">
          {{ successMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <v-btn-toggle v-model="viewMode" mandatory>
                  <v-btn value="tree" size="small">
                    <v-icon start>mdi-file-tree</v-icon>
                    树形视图
                  </v-btn>
                  <v-btn value="table" size="small">
                    <v-icon start>mdi-table</v-icon>
                    表格视图
                  </v-btn>
                </v-btn-toggle>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn variant="outlined" class="mr-2" :loading="loading" @click="exportMenus">
                  <v-icon start>mdi-download</v-icon>
                  导出
                </v-btn>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建菜单
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>

            <!-- 树形视图 -->
            <div v-if="viewMode === 'tree' && !loading">
              <MenuTree
                v-if="menus.length > 0"
                :menus="menus"
                @edit="openEditDialog"
                @delete="openDeleteDialog"
                @update:menus="handleMenusReorder"
              />
              <v-alert v-else type="info">暂无菜单，点击"新建菜单"创建第一个菜单</v-alert>
            </div>

            <!-- 表格视图（可选实现） -->
            <div v-if="viewMode === 'table' && !loading">
              <v-alert type="info" class="mb-4"> 表格视图展示所有菜单（扁平化） </v-alert>
              <v-table v-if="menus.length > 0">
                <template #default>
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>标题</th>
                      <th>路径</th>
                      <th>图标</th>
                      <th>排序</th>
                      <th>可见</th>
                      <th>操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <!-- 扁平化菜单展示（简化版） -->
                    <template v-for="menu in menus" :key="menu.id">
                      <tr>
                        <td>{{ menu.id }}</td>
                        <td>{{ menu.title }}</td>
                        <td>
                          <code>{{ menu.path }}</code>
                        </td>
                        <td>
                          <v-icon v-if="menu.icon" size="small">{{ menu.icon }}</v-icon>
                        </td>
                        <td>{{ menu.order }}</td>
                        <td>
                          <v-chip :color="menu.visible ? 'success' : 'error'" size="small">{{
                            menu.visible ? "是" : "否"
                          }}</v-chip>
                        </td>
                        <td>
                          <v-btn icon="mdi-pencil" size="small" variant="text" @click="openEditDialog(menu)"></v-btn>
                          <v-btn
                            icon="mdi-delete"
                            size="small"
                            variant="text"
                            color="error"
                            @click="openDeleteDialog(menu)"
                          ></v-btn>
                        </td>
                      </tr>
                    </template>
                  </tbody>
                </template>
              </v-table>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <MenuDialog
      v-model="menuDialog"
      :menu="selectedMenu"
      :mode="dialogMode"
      :parent-menus="menus"
      @save="handleSaveMenu"
    />

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除菜单 <strong>{{ menuToDelete?.title }}</strong> 吗？
          <v-alert
            v-if="menuToDelete?.children && menuToDelete.children.length > 0"
            type="warning"
            class="mt-2"
            density="compact"
          >
            该菜单有子菜单，请先删除子菜单
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="elevated" :loading="loading" @click="confirmDelete">删除</v-btn>
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
