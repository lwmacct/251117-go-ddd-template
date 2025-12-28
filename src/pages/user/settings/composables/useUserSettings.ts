/**
 * User 个人设置 Composable
 * 与 Admin Settings 的主要区别：
 * - 使用 value 字段（用户实际值）而非 default_value
 * - 支持重置到默认值（删除用户自定义）
 * - 支持按 Category 懒加载
 *
 * 懒加载功能由 useLazyCategorySchema 提供，包括：
 * - 按 category 加载数据
 * - 竞态条件防护（pending Set）
 * - 自动合并已有数据
 */
import { ref, computed } from "vue";
import { userSettingsApi } from "@/api";
import { useLazyCategorySchema } from "@/composables";
import type {
  SettingCategoryMetaDTO,
  SettingSchemaCategoryDTO,
  SettingSchemaSettingDTO,
  HandlerBatchSetUserSettingsRequest,
  HandlerSetUserSettingRequest,
} from "@models";

export function useUserSettings() {
  // 分类元信息（用于渲染 tabs，不含 settings 数据）
  const categories = ref<SettingCategoryMetaDTO[]>([]);

  // 使用封装的懒加载 composable
  const {
    schema,
    loadedCategories,
    loading: schemaLoading,
    fetchSchemaByCategory,
    isCategoryLoaded,
    reset: resetSchema,
  } = useLazyCategorySchema<SettingSchemaCategoryDTO>(async (categoryKey) => {
    const response = await userSettingsApi.apiUserSettingsGet(categoryKey);
    return (response.data.data ?? []) as SettingSchemaCategoryDTO[];
  });

  // 其他加载状态
  const loading = ref(false);
  const saving = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  // 按 key 索引的设置 Map
  const settingsMap = computed(() => {
    const map = new Map<string, SettingSchemaSettingDTO>();
    schema.value.forEach((cat) => {
      cat.groups?.forEach((group) => {
        group.settings?.forEach((s) => {
          if (s.key) {
            map.set(s.key, s);
          }
        });
      });
    });
    return map;
  });

  /**
   * 获取分类列表（用于渲染 tabs）
   * 应在页面初始化时首先调用
   */
  const fetchCategories = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await userSettingsApi.apiUserSettingsCategoriesGet();
      categories.value = (response.data.data ?? []) as SettingCategoryMetaDTO[];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取分类列表失败";
      console.error("Failed to fetch user setting categories:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取用户配置（全量）
   * 注意：全量加载会重置懒加载状态，然后逐个加载所有分类
   */
  const fetchSchema = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      // 重置懒加载状态
      resetSchema();

      // 获取全量数据
      const response = await userSettingsApi.apiUserSettingsGet();
      const allCategories = (response.data.data ?? []) as SettingSchemaCategoryDTO[];

      // 逐个分类触发加载（标记为已加载）
      for (const cat of allCategories) {
        if (cat.category) {
          await fetchSchemaByCategory(cat.category);
        }
      }
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取设置失败";
      console.error("Failed to fetch user settings:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 批量更新用户设置
   */
  const batchUpdateSettings = async (updates: { key: string; value: object }[]): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      const settingsData = updates.map((u) => ({
        key: u.key,
        value: u.value,
      }));

      const data: HandlerBatchSetUserSettingsRequest = { settings: settingsData };
      await userSettingsApi.apiUserSettingsBatchPost(data);

      // 重新获取 schema 以更新 is_customized 标识
      await fetchSchema();

      successMessage.value = "设置保存成功";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "保存设置失败";
      console.error("Failed to batch update user settings:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  /**
   * 静默更新单个设置（不显示成功提示，用于自动保存）
   */
  const updateSettingQuietly = async (key: string, value: object): Promise<boolean> => {
    try {
      const data: HandlerSetUserSettingRequest = { value };
      await userSettingsApi.apiUserSettingsKeyPut(key, data);

      // 更新本地 schema 中对应项的 value 和 is_customized
      // 由于 schema 使用 shallowRef，需要创建新引用触发响应式
      const updated = schema.value.map((cat) => ({
        ...cat,
        groups: cat.groups?.map((group) => ({
          ...group,
          settings: group.settings?.map((s) => (s.key === key ? { ...s, value, is_customized: true } : s)),
        })),
      }));
      schema.value = updated;

      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "保存失败";
      console.error("Failed to update user setting:", error);
      return false;
    }
  };

  /**
   * 重置单个设置到系统默认值
   */
  const resetSetting = async (key: string): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";

    try {
      // 从 schema 中找到包含该设置的分类（设置对象只有 category_id，没有 category 字符串）
      let categoryKey: string | undefined;
      for (const cat of schema.value) {
        const found = cat.groups?.some((g) => g.settings?.some((s) => s.key === key));
        if (found) {
          categoryKey = cat.category;
          break;
        }
      }

      await userSettingsApi.apiUserSettingsKeyDelete(key);

      // 只刷新该设置所属的分类，而非全量刷新
      if (categoryKey) {
        // 从 loaded 集合中移除该分类，强制重新加载
        // 注意：loadedCategories 是 ShallowRef<Set>，需创建新引用触发响应式
        const newLoaded = new Set(loadedCategories.value);
        newLoaded.delete(categoryKey);
        loadedCategories.value = newLoaded;

        // 只刷新这一个分类
        await fetchSchemaByCategory(categoryKey);
      }

      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "重置设置失败";
      console.error("Failed to reset user setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  /**
   * 检查设置是否已被用户自定义
   */
  const isCustomized = (key: string): boolean => {
    const setting = settingsMap.value.get(key);
    return setting?.is_customized ?? false;
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  return {
    categories,
    schema,
    settingsMap,
    loadedCategories,
    loading,
    schemaLoading,
    saving,
    errorMessage,
    successMessage,
    fetchCategories,
    fetchSchema,
    fetchSchemaByCategory,
    isCategoryLoaded,
    batchUpdateSettings,
    updateSettingQuietly,
    resetSetting,
    isCustomized,
    clearMessages,
  };
}
