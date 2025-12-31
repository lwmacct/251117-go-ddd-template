<script setup lang="ts">
import { ref } from "vue";
import { useAdminProducts } from "./composables/useAdminProducts";
import { ITEMS_PER_PAGE_OPTIONS } from "@/composables";
import ProductDialog from "./components/ProductDialog.vue";
import CopyButton from "@/components/CopyButton.vue";
import type { ProductProductDTO, ProductCreateProductDTO, ProductUpdateProductDTO } from "@models";

/**
 * 产品管理页面
 * 用于查看和管理系统产品
 */

// 使用 composable
const {
  products,
  loading,
  searchQuery,
  statusFilter,
  pagination,
  fetchProducts: _fetchProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  onTableOptionsUpdate,
  exportProducts,
} = useAdminProducts();

// 对话框状态
const productDialog = ref(false);
const deleteDialog = ref(false);

// 编辑状态
const dialogMode = ref<"create" | "edit">("create");
const selectedProduct = ref<ProductProductDTO | null>(null);
const productToDelete = ref<ProductProductDTO | null>(null);

// 状态选项
const statusOptions = [
  { title: "全部", value: "" },
  { title: "启用", value: "active" },
  { title: "禁用", value: "inactive" },
];

// 表头配置
const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "产品名称", key: "name", sortable: true },
  { title: "描述", key: "description" },
  { title: "价格", key: "price", sortable: true },
  { title: "状态", key: "status", sortable: true },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

// 打开创建对话框
const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedProduct.value = null;
  productDialog.value = true;
};

// 打开编辑对话框
const openEditDialog = (product: ProductProductDTO) => {
  dialogMode.value = "edit";
  selectedProduct.value = product;
  productDialog.value = true;
};

// 打开删除确认对话框
const openDeleteDialog = (product: ProductProductDTO) => {
  productToDelete.value = product;
  deleteDialog.value = true;
};

// 保存产品（创建或编辑）
const handleSaveProduct = async (data: ProductCreateProductDTO | ProductUpdateProductDTO) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createProduct(data as ProductCreateProductDTO);
  } else if (selectedProduct.value?.id) {
    success = await updateProduct(selectedProduct.value.id, data as ProductUpdateProductDTO);
  }

  if (success) {
    productDialog.value = false;
  }
};

// 确认删除
const confirmDelete = async () => {
  if (!productToDelete.value?.id) return;

  const success = await deleteProduct(productToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    productToDelete.value = null;
  }
};

// 格式化日期
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

// 格式化价格
const formatPrice = (price: number) => {
  return `¥${price.toFixed(2)}`;
};

// 状态颜色
const getStatusColor = (status?: string) => {
  const colors: Record<string, string> = {
    active: "success",
    inactive: "error",
  };
  return colors[status ?? ""] || "default";
};

// 状态文本
const getStatusText = (status?: string) => {
  const texts: Record<string, string> = {
    active: "启用",
    inactive: "禁用",
  };
  return texts[status ?? ""] || (status ?? "-");
};
</script>

<template>
  <div class="products-page">
    <!-- 标题 -->
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">产品管理</h1>
      </v-col>
    </v-row>

    <!-- 产品列表卡片 -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="searchQuery"
                  prepend-inner-icon="mdi-magnify"
                  label="搜索产品"
                  single-line
                  hide-details
                  clearable
                  variant="outlined"
                  density="compact"
                  placeholder="输入后自动搜索..."
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="2">
                <v-select
                  v-model="statusFilter"
                  label="状态筛选"
                  :items="statusOptions"
                  item-title="title"
                  item-value="value"
                  hide-details
                  clearable
                  variant="outlined"
                  density="compact"
                ></v-select>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn variant="outlined" class="mr-2" :loading="loading" @click="exportProducts">
                  <v-icon start>mdi-download</v-icon>
                  导出
                </v-btn>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建产品
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :items-per-page-options="ITEMS_PER_PAGE_OPTIONS"
              :page="pagination.page"
              :headers="headers"
              :items="products"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无产品数据"
              @update:options="onTableOptionsUpdate"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <div class="d-flex align-center">
                  <span>{{ item.id }}</span>
                  <CopyButton :text="String(item.id)" size="x-small" />
                </div>
              </template>

              <!-- 描述列 -->
              <template #item.description="{ item }">
                <span class="text-body-2 text-truncate d-inline-block" style="max-width: 200px">
                  {{ item.description || "-" }}
                </span>
              </template>

              <!-- 价格列 -->
              <template #item.price="{ item }">
                <span class="text-body-2">{{ item.price !== undefined ? formatPrice(item.price) : "-" }}</span>
              </template>

              <!-- 状态列 -->
              <template #item.status="{ item }">
                <v-chip :color="getStatusColor(item.status)" size="small">
                  {{ getStatusText(item.status) }}
                </v-chip>
              </template>

              <!-- 创建时间列 -->
              <template #item.created_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
              </template>

              <!-- 操作列 -->
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
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 创建/编辑对话框 -->
    <ProductDialog v-model="productDialog" :product="selectedProduct" :mode="dialogMode" @save="handleSaveProduct" />

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除产品 <strong>{{ productToDelete?.name }}</strong> 吗？此操作不可恢复。
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
.products-page {
  width: 100%;
}
</style>
