/**
 * Admin 产品管理 Composable
 */
import { ref, watch } from "vue";
import { adminProductApi, extractList, extractData } from "@/api";
import type { ProductProductDTO, ProductCreateProductDTO, ProductUpdateProductDTO } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";
import { useServerPagination, useSnackbar } from "@/composables";

export function useAdminProducts() {
  // 状态
  const products = ref<ProductProductDTO[]>([]);
  const searchQuery = ref("");
  const statusFilter = ref<"active" | "inactive" | "">("");
  // 防抖搜索值，300ms 延迟
  const debouncedSearchQuery = refDebounced(searchQuery, 300);

  // 使用通用分页 composable
  const {
    pagination,
    loading,
    onTableOptionsUpdate: baseOnTableOptionsUpdate,
    updateTotal,
    resetAndFetch,
    getParams,
  } = useServerPagination();

  // 消息提示
  const { success, error } = useSnackbar();

  /**
   * 获取产品列表
   */
  const fetchProducts = async () => {
    loading.value = true;

    try {
      const { limit, page } = getParams();
      const response = await adminProductApi.apiAdminProductsGet(limit, page);
      const result = extractList<ProductProductDTO>(response.data);
      products.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (err) {
      error((err as Error).message || "获取产品列表失败");
      console.error("Failed to fetch products:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取单个产品详情
   */
  const fetchProduct = async (id: number): Promise<ProductProductDTO | null> => {
    try {
      const response = await adminProductApi.apiAdminProductsIdGet(id);
      return extractData<ProductProductDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取产品详情失败");
      console.error("Failed to fetch product:", err);
      return null;
    }
  };

  /**
   * 创建产品
   */
  const createProduct = async (data: ProductCreateProductDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminProductApi.apiAdminProductsPost(data);
      success("产品创建成功");
      await fetchProducts(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "创建产品失败");
      console.error("Failed to create product:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新产品
   */
  const updateProduct = async (id: number, data: ProductUpdateProductDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminProductApi.apiAdminProductsIdPut(id, data);
      success("产品更新成功");
      await fetchProducts(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "更新产品失败");
      console.error("Failed to update product:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 删除产品
   */
  const deleteProduct = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await adminProductApi.apiAdminProductsIdDelete(id);
      success("产品删除成功");
      await fetchProducts(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "删除产品失败");
      console.error("Failed to delete product:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 监听搜索值变化，自动触发搜索
  watch(debouncedSearchQuery, () => {
    resetAndFetch(fetchProducts);
  });

  // 监听状态筛选变化
  watch(statusFilter, () => {
    resetAndFetch(fetchProducts);
  });

  /**
   * 表格选项变化处理（分页、每页条数、排序）
   * 由 v-data-table-server 的 @update:options 触发
   */
  const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    baseOnTableOptionsUpdate(options, fetchProducts);
  };

  /**
   * 导出产品列表为 CSV
   */
  const exportProducts = async () => {
    loading.value = true;

    try {
      // 获取所有产品（最多 1000 条）
      const response = await adminProductApi.apiAdminProductsGet(1000, 1);
      const result = extractList<ProductProductDTO>(response.data);

      if (result.data.length === 0) {
        error("没有可导出的数据");
        return;
      }

      // 过滤数据
      let exportData = result.data;
      if (statusFilter.value) {
        exportData = exportData.filter((p) => p.status === statusFilter.value);
      }

      // 定义 CSV 列
      const columns: CSVColumn<ProductProductDTO>[] = [
        { header: "ID", key: "id" },
        { header: "产品名称", key: "name" },
        { header: "描述", key: "description" },
        { header: "价格", key: "price" },
        {
          header: "状态",
          key: (item) => (item.status === "active" ? "启用" : "禁用"),
        },
        {
          header: "创建时间",
          key: (item) =>
            item.created_at
              ? new Date(item.created_at).toLocaleString("zh-CN", {
                  year: "numeric",
                  month: "2-digit",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                })
              : "-",
        },
        {
          header: "更新时间",
          key: (item) =>
            item.updated_at
              ? new Date(item.updated_at).toLocaleString("zh-CN", {
                  year: "numeric",
                  month: "2-digit",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                })
              : "-",
        },
      ];

      // 生成文件名
      const timestamp = new Date().toISOString().slice(0, 10);
      const filename = `产品列表_${timestamp}.csv`;

      // 导出
      exportToCSV(exportData, columns, { filename, withBOM: true });
      success(`成功导出 ${exportData.length} 条产品记录`);
    } catch (err) {
      error((err as Error).message || "导出失败");
      console.error("Failed to export products:", err);
    } finally {
      loading.value = false;
    }
  };

  return {
    // 状态
    products,
    loading,
    searchQuery,
    debouncedSearchQuery,
    statusFilter,
    pagination,

    // 方法
    fetchProducts,
    fetchProduct,
    createProduct,
    updateProduct,
    deleteProduct,
    onTableOptionsUpdate,
    exportProducts,
  };
}
