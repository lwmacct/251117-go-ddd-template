<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useSettingCategories } from "./composables/useSettingCategories";
import CategoryDialog from "./components/CategoryDialog.vue";
import type { SettingCategoryDTO, HandlerCreateCategoryRequest, HandlerUpdateCategoryRequest } from "@models";

const { categories, loading, fetchCategories, createCategory, updateCategory, deleteCategory } = useSettingCategories();

const categoryDialog = ref(false);
const deleteDialog = ref(false);
const dialogMode = ref<"create" | "edit">("create");
const selectedCategory = ref<SettingCategoryDTO | null>(null);
const categoryToDelete = ref<SettingCategoryDTO | null>(null);

const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "分类标识", key: "key", sortable: true },
  { title: "显示名称", key: "label", sortable: true },
  { title: "图标", key: "icon" },
  { title: "排序权重", key: "order", sortable: true },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

onMounted(() => {
  fetchCategories();
});

const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedCategory.value = null;
  categoryDialog.value = true;
};

const openEditDialog = (category: SettingCategoryDTO) => {
  dialogMode.value = "edit";
  selectedCategory.value = category;
  categoryDialog.value = true;
};

const openDeleteDialog = (category: SettingCategoryDTO) => {
  categoryToDelete.value = category;
  deleteDialog.value = true;
};

const handleSaveCategory = async (data: HandlerCreateCategoryRequest | HandlerUpdateCategoryRequest) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createCategory(data as HandlerCreateCategoryRequest);
  } else if (selectedCategory.value?.id) {
    success = await updateCategory(selectedCategory.value.id, data as HandlerUpdateCategoryRequest);
  }

  if (success) {
    categoryDialog.value = false;
  }
};

const confirmDelete = async () => {
  if (!categoryToDelete.value?.id) return;

  const success = await deleteCategory(categoryToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    categoryToDelete.value = null;
  }
};

const formatDate = (dateString?: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
};
</script>

<template>
  <div class="setting-categories-page">
    <v-row>
      <v-col cols="12">
        <div class="d-flex align-center mb-6">
          <h1 class="text-h4">配置分类管理</h1>
          <v-spacer />
          <v-btn variant="outlined" class="mr-2" to="/admin/settings">
            <v-icon start>mdi-arrow-left</v-icon>
            返回系统设置
          </v-btn>
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <span class="text-h6">分类列表</span>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn variant="outlined" class="mr-2" :loading="loading" @click="fetchCategories">
                  <v-icon start>mdi-refresh</v-icon>
                  刷新
                </v-btn>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建分类
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table
              :headers="headers"
              :items="categories"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无分类数据"
              items-per-page="-1"
              hide-default-footer
            >
              <template #item.icon="{ item }">
                <v-chip size="small" variant="outlined">
                  <v-icon start size="small">{{ item.icon }}</v-icon>
                  {{ item.icon }}
                </v-chip>
              </template>

              <template #item.order="{ item }">
                <v-chip size="small" color="primary" variant="tonal">
                  {{ item.order }}
                </v-chip>
              </template>

              <template #item.created_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
              </template>

              <template #item.actions="{ item }">
                <v-tooltip text="编辑">
                  <template #activator="{ props }">
                    <v-btn icon="mdi-pencil" size="small" variant="text" v-bind="props" @click="openEditDialog(item)"></v-btn>
                  </template>
                </v-tooltip>

                <v-tooltip text="删除">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      v-bind="props"
                      @click="openDeleteDialog(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <CategoryDialog v-model="categoryDialog" :category="selectedCategory" :mode="dialogMode" @save="handleSaveCategory" />

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除分类
          <strong>{{ categoryToDelete?.label }}</strong>
          （{{ categoryToDelete?.key }}）吗？此操作不可恢复。
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="elevated" :loading="loading" @click="confirmDelete"> 删除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.setting-categories-page {
  width: 100%;
}
</style>
