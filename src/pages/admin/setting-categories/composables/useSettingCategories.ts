/**
 * Admin 配置分类管理 Composable
 */
import { ref } from "vue";
import { adminSettingCategoriesApi } from "@/api";
import type { SettingCategoryDTO, HandlerCreateCategoryRequest, HandlerUpdateCategoryRequest } from "@models";

export function useSettingCategories() {
  const categories = ref<SettingCategoryDTO[]>([]);
  const loading = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  /**
   * 获取分类列表
   */
  const fetchCategories = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminSettingCategoriesApi.apiAdminSettingsCategoriesGet();
      categories.value = (response.data.data ?? []) as SettingCategoryDTO[];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取分类列表失败";
      console.error("Failed to fetch categories:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 创建分类
   */
  const createCategory = async (data: HandlerCreateCategoryRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesPost(data);
      successMessage.value = "分类创建成功";
      await fetchCategories();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建分类失败";
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
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesIdPut(id, data);
      successMessage.value = "分类更新成功";
      await fetchCategories();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新分类失败";
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
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminSettingCategoriesApi.apiAdminSettingsCategoriesIdDelete(id);
      successMessage.value = "分类删除成功";
      await fetchCategories();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除分类失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  return {
    categories,
    loading,
    errorMessage,
    successMessage,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    clearMessages,
  };
}
