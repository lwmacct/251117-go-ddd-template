/**
 * Admin 系统设置 Composable
 * 支持按 Category 懒加载，减少首屏数据量
 *
 * 懒加载功能由 useLazyCategorySchema 提供，包括：
 * - 按 category 加载数据
 * - 竞态条件防护（pending Set）
 * - 自动合并已有数据
 */
import { ref, computed } from "vue";
import { adminSettingsApi, adminSettingCategoriesApi, extractData } from "@/api";
import { useLazyCategorySchema, useSnackbar } from "@/composables";
import {
  type SettingSettingDTO,
  type SettingSettingsCategoryDTO,
  type SettingSettingsGroupDTO,
  type SettingSettingsItemDTO,
  type SettingCategoryDTO,
  type HandlerCreateSettingRequest,
  type HandlerUpdateSettingRequest,
} from "@models";

export function useSettings() {
  const settings = ref<SettingSettingDTO[]>([]);
  const categories = ref<SettingCategoryDTO[]>([]); // 分类列表（用于渲染 tabs）
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
    const response = await adminSettingsApi.apiAdminSettingsGet(categoryKey);
    return (response.data.data ?? []) as SettingSettingsCategoryDTO[];
  });

  // 其他加载状态
  const loading = ref(false);
  const saving = ref(false);

  // 按分类 ID 缓存的设置
  const settingsByCategory = computed(() => {
    const map = new Map<number, Map<string, SettingSettingDTO>>();
    settings.value.forEach((setting) => {
      const categoryId = setting.category_id ?? 0;
      const key = setting.key ?? "";
      if (!key) return;
      if (!map.has(categoryId)) {
        map.set(categoryId, new Map());
      }
      map.get(categoryId)!.set(key, setting);
    });
    return map;
  });

  // 按 key 索引的设置 Map（使用 Schema 数据）
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
      const response = await adminSettingCategoriesApi.apiAdminSettingsCategoriesGet();
      categories.value = (response.data.data ?? []) as SettingCategoryDTO[];
    } catch (error) {
      const msg = (error as Error).message || "获取分类列表失败";
      snackbar.error(msg);
      console.error("Failed to fetch setting categories:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取配置（全量）
   * 注意：全量加载会重置懒加载状态，然后逐个加载所有分类
   */
  const fetchSchema = async () => {
    loading.value = true;

    try {
      // 重置懒加载状态
      resetSchema();

      // 获取全量数据
      const response = await adminSettingsApi.apiAdminSettingsGet();
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
      console.error("Failed to fetch settings:", error);
    } finally {
      loading.value = false;
    }
  };

  // 获取单个设置的值
  const getSetting = <T = string>(key: string, defaultValue: T): T => {
    const setting = settings.value.find((s) => s.key === key);
    if (!setting || setting.default_value === undefined) return defaultValue;

    // 根据值类型解析（default_value 是 JSONB，可能是任意类型）
    const value = setting.default_value as unknown;
    switch (setting.value_type) {
      case "boolean":
        return (value === true || value === "true") as T;
      case "number":
        return Number(value) as T;
      case "json":
        try {
          return (typeof value === "string" ? JSON.parse(value) : value) as T;
        } catch {
          return defaultValue;
        }
      default:
        return value as T;
    }
  };

  // 创建单个设置
  const createSetting = async (data: HandlerCreateSettingRequest): Promise<boolean> => {
    saving.value = true;

    try {
      const response = await adminSettingsApi.apiAdminSettingsPost(data);
      const newSetting = extractData<SettingSettingDTO>(response.data);
      if (newSetting) {
        settings.value.push(newSetting);
      }
      snackbar.success("设置创建成功");
      return true;
    } catch (error) {
      const msg = (error as Error).message || "创建设置失败";
      snackbar.error(msg);
      console.error("Failed to create setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  // 更新单个设置
  const updateSetting = async (key: string, value: object): Promise<boolean> => {
    saving.value = true;

    try {
      const data: HandlerUpdateSettingRequest = { default_value: value };
      const response = await adminSettingsApi.apiAdminSettingsKeyPut(key, data);
      const updated = extractData<SettingSettingDTO>(response.data);
      const index = settings.value.findIndex((s) => s.key === key);
      if (index !== -1 && updated) {
        settings.value[index] = updated;
      }
      snackbar.success("设置更新成功");
      return true;
    } catch (error) {
      const msg = (error as Error).message || "更新设置失败";
      snackbar.error(msg);
      console.error("Failed to update setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  // 静默更新单个设置（用于开关等即时生效的控件，不显示成功提示）
  const updateSettingQuietly = async (key: string, value: object): Promise<boolean> => {
    try {
      const data: HandlerUpdateSettingRequest = { default_value: value };
      const response = await adminSettingsApi.apiAdminSettingsKeyPut(key, data);
      const updated = extractData<SettingSettingDTO>(response.data);
      const index = settings.value.findIndex((s) => s.key === key);
      if (index !== -1 && updated) {
        settings.value[index] = updated;
      }
      return true;
    } catch (error) {
      const msg = (error as Error).message || "更新设置失败";
      snackbar.error(msg);
      console.error("Failed to update setting:", error);
      return false;
    }
  };

  // 删除设置
  const deleteSetting = async (key: string): Promise<boolean> => {
    saving.value = true;

    try {
      await adminSettingsApi.apiAdminSettingsKeyDelete(key);
      settings.value = settings.value.filter((s) => s.key !== key);
      snackbar.success("设置删除成功");
      return true;
    } catch (error) {
      const msg = (error as Error).message || "删除设置失败";
      snackbar.error(msg);
      console.error("Failed to delete setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  return {
    settings,
    categories,
    schema,
    settingsByCategory,
    settingsMap,
    loadedCategories,
    loading,
    schemaLoading,
    saving,
    fetchCategories,
    fetchSchema,
    fetchSchemaByCategory,
    isCategoryLoaded,
    getSetting,
    createSetting,
    updateSetting,
    updateSettingQuietly,
    deleteSetting,
  };
}
