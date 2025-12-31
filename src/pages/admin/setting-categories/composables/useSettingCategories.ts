/**
 * Admin 配置分类管理 Composable
 */
import { ref } from "vue";
import { adminSettingCategoriesApi } from "@/api";
import type { SettingCategoryDTO, HandlerCreateCategoryRequest, HandlerUpdateCategoryRequest } from "@models";
import { useSnackbar } from "@/composables";

export function useSettingCategories() {
  const categories = ref<SettingCategoryDTO[]>([]);
  const loading = ref(false);

  // 消息提示
  const { success, error } = useSnackbar();

  /**
   * 获取分类列表
   */
  const fetchCategories = async () => {
    loading.value = true;

    try {
      const response = await adminSettingCategoriesApi.apiAdminSettingsCategoriesGet();
      categories.value = (response.data.data ?? []) as SettingCategoryDTO[];
    } catch (err) {
      error((err as Error).message || "获取分类列表失败");
      console.error("Failed to fetch categories:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 创建分类
   */
  const createCategory = async (data: HandlerCreateCategoryRequest): Promise<boolean> => {
    loading.value = true;

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesPost(data);
      success("分类创建成功");
      await fetchCategories();
      return true;
    } catch (err) {
      error((err as Error).message || "创建分类失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新分类
   */
  const updateCategory = async (id: number, data: HandlerUpdateCategoryRequest): Promise<boolean> => {
    loading.value = true;

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesIdPut(id, data);
      success("分类更新成功");
      await fetchCategories();
      return true;
    } catch (err) {
      error((err as Error).message || "更新分类失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 删除分类
   */
  const deleteCategory = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesIdDelete(id);
      success("分类删除成功");
      await fetchCategories();
      return true;
    } catch (err) {
      error((err as Error).message || "删除分类失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  return {
    categories,
    loading,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  };
}
