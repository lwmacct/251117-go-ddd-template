/**
 * Admin 系统设置 Composable
 */
import { ref, computed } from "vue";
import { adminSettingsApi, extractData } from "@/api";
import {
  type SettingSettingDTO,
  type SettingSchemaCategoryDTO,
  type SettingSchemaSettingDTO,
  type HandlerCreateSettingRequest,
  type HandlerUpdateSettingRequest,
  type HandlerBatchUpdateSettingsRequest,
} from "@models";

export function useSettings() {
  const settings = ref<SettingSettingDTO[]>([]);
  const schema = ref<SettingSchemaCategoryDTO[]>([]);
  const loading = ref(false);
  const saving = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

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

  // 获取配置 Schema（层级结构）
  const fetchSchema = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminSettingsApi.apiAdminSettingsSchemaGet();
      schema.value = (response.data.data ?? []) as SettingSchemaCategoryDTO[];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取设置 Schema 失败";
      console.error("Failed to fetch settings schema:", error);
    } finally {
      loading.value = false;
    }
  };

  // 获取所有设置
  const fetchSettings = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminSettingsApi.apiAdminSettingsGet();
      settings.value = (response.data.data ?? []) as SettingSettingDTO[];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取设置失败";
      console.error("Failed to fetch settings:", error);
    } finally {
      loading.value = false;
    }
  };

  // 获取指定分类的设置
  const fetchSettingsByCategory = async (categoryId: number) => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminSettingsApi.apiAdminSettingsGet(categoryId);
      const categorySettings = (response.data.data ?? []) as SettingSettingDTO[];
      // 更新现有设置列表（合并）
      const existingKeys = new Set(settings.value.map((s) => s.key));
      categorySettings.forEach((s) => {
        if (existingKeys.has(s.key)) {
          // 更新现有
          const index = settings.value.findIndex((setting) => setting.key === s.key);
          if (index !== -1) {
            settings.value[index] = s;
          }
        } else {
          // 添加新设置
          settings.value.push(s);
        }
      });
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取设置失败";
      console.error("Failed to fetch settings by category:", error);
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
    errorMessage.value = "";

    try {
      const response = await adminSettingsApi.apiAdminSettingsPost(data);
      const newSetting = extractData<SettingSettingDTO>(response.data);
      if (newSetting) {
        settings.value.push(newSetting);
      }
      successMessage.value = "设置创建成功";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建设置失败";
      console.error("Failed to create setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  // 更新单个设置
  const updateSetting = async (key: string, value: object): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";

    try {
      const data: HandlerUpdateSettingRequest = { default_value: value };
      const response = await adminSettingsApi.apiAdminSettingsKeyPut(key, data);
      const updated = extractData<SettingSettingDTO>(response.data);
      const index = settings.value.findIndex((s) => s.key === key);
      if (index !== -1 && updated) {
        settings.value[index] = updated;
      }
      successMessage.value = "设置更新成功";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新设置失败";
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
      errorMessage.value = (error as Error).message || "更新设置失败";
      console.error("Failed to update setting:", error);
      return false;
    }
  };

  // 批量更新设置
  const batchUpdateSettings = async (updates: { key: string; value: object }[]): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      // 构建批量更新请求
      const settingsData = updates.map((u) => ({
        key: u.key,
        value: u.value,
      }));

      const data: HandlerBatchUpdateSettingsRequest = { settings: settingsData };
      await adminSettingsApi.apiAdminSettingsBatchPost(data);

      // 更新本地缓存
      updates.forEach((update) => {
        const index = settings.value.findIndex((s) => s.key === update.key);
        const current = index !== -1 ? settings.value[index] : undefined;
        if (current) {
          current.default_value = update.value;
        }
      });

      successMessage.value = "设置保存成功";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "批量更新设置失败";
      console.error("Failed to batch update settings:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  // 删除设置
  const deleteSetting = async (key: string): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";

    try {
      await adminSettingsApi.apiAdminSettingsKeyDelete(key);
      settings.value = settings.value.filter((s) => s.key !== key);
      successMessage.value = "设置删除成功";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除设置失败";
      console.error("Failed to delete setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  return {
    settings,
    schema,
    settingsByCategory,
    settingsMap,
    loading,
    saving,
    errorMessage,
    successMessage,
    fetchSchema,
    fetchSettings,
    fetchSettingsByCategory,
    getSetting,
    createSetting,
    updateSetting,
    updateSettingQuietly,
    batchUpdateSettings,
    deleteSetting,
    clearMessages,
  };
}
