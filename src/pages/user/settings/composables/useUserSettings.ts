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
import { useLazyCategorySchema, useSnackbar } from "@/composables";
import type {
  SettingCategoryMetaDTO,
  SettingSettingsCategoryDTO,
  SettingSettingsGroupDTO,
  SettingSettingsItemDTO,
  HandlerSetUserSettingRequest,
} from "@models";

export function useUserSettings() {
  // 分类元信息（用于渲染 tabs，不含 settings 数据）
  const categories = ref<SettingCategoryMetaDTO[]>([]);
  const snackbar = useSnackbar();

  // 使用封装的懒加载 composable
  const {
    schema,
    loadedCategories,
    loading: schemaLoading,
    fetchSchemaByCategory,
    isCategoryLoaded,
    reset: resetSchema,
  } = useLazyCategorySchema<SettingSettingsCategoryDTO>(async (categoryKey) => {
    const response = await userSettingsApi.apiUserSettingsGet(categoryKey);
    return (response.data.data ?? []) as SettingSettingsCategoryDTO[];
  });

  // 其他加载状态
  const loading = ref(false);
  const saving = ref(false);

  // 按 key 索引的设置 Map
  const settingsMap = computed(() => {
    const map = new Map<string, SettingSettingsItemDTO>();
    schema.value.forEach((cat) => {
      cat.groups?.forEach((group: SettingSettingsGroupDTO) => {
        group.settings?.forEach((s: SettingSettingsItemDTO) => {
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

    try {
      const response = await userSettingsApi.apiUserSettingsCategoriesGet();
      categories.value = (response.data.data ?? []) as SettingCategoryMetaDTO[];
    } catch (error) {
      const msg = (error as Error).message || "获取分类列表失败";
      snackbar.error(msg);
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

    try {
      // 重置懒加载状态
      resetSchema();

      // 获取全量数据
      const response = await userSettingsApi.apiUserSettingsGet();
      const allCategories = (response.data.data ?? []) as SettingSettingsCategoryDTO[];

      // 逐个分类触发加载（标记为已加载）
      for (const cat of allCategories) {
        if (cat.category) {
          await fetchSchemaByCategory(cat.category);
        }
      }
    } catch (error) {
      const msg = (error as Error).message || "获取设置失败";
      snackbar.error(msg);
      console.error("Failed to fetch user settings:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 静默更新单个设置（不显示成功提示，用于自动保存）
   * 注意：不更新 schema，formValues 已是最新值，避免触发不必要的重绘
   */
  const updateSettingQuietly = async (key: string, value: object): Promise<boolean> => {
    try {
      const data: HandlerSetUserSettingRequest = { value };
      await userSettingsApi.apiUserSettingsKeyPut(key, data);
      return true;
    } catch (error) {
      const msg = (error as Error).message || "保存失败";
      snackbar.error(msg);
      console.error("Failed to update user setting:", error);
      return false;
    }
  };

  /**
   * 重置单个设置到系统默认值
   * 注意：调用方需要手动更新表单值为默认值
   */
  const resetSetting = async (key: string): Promise<boolean> => {
    saving.value = true;

    try {
      await userSettingsApi.apiUserSettingsKeyDelete(key);
      snackbar.success("已重置为默认值");
      return true;
    } catch (error) {
      const msg = (error as Error).message || "重置设置失败";
      snackbar.error(msg);
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

  return {
    categories,
    schema,
    settingsMap,
    loadedCategories,
    loading,
    schemaLoading,
    saving,
    fetchCategories,
    fetchSchema,
    fetchSchemaByCategory,
    isCategoryLoaded,
    updateSettingQuietly,
    resetSetting,
    isCustomized,
  };
}
