/**
 * Admin 系统设置 Composable
 */
import { ref, computed } from "vue";
import { adminSettingsApi, extractData } from "@/api";
import {
  type SettingSettingDTO,
  type HandlerCreateSettingRequest,
  type HandlerUpdateSettingRequest,
  type HandlerBatchUpdateSettingsRequest,
} from "@models";

export function useSettings() {
  const settings = ref<SettingSettingDTO[]>([]);
  const loading = ref(false);
  const saving = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  // 按分类缓存的设置
  const settingsByCategory = computed(() => {
    const map = new Map<string, Map<string, SettingSettingDTO>>();
    settings.value.forEach((setting) => {
      const category = setting.category ?? "default";
      const key = setting.key ?? "";
      if (!key) return;
      if (!map.has(category)) {
        map.set(category, new Map());
      }
      map.get(category)!.set(key, setting);
    });
    return map;
  });

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
  const fetchSettingsByCategory = async (category: string) => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminSettingsApi.apiAdminSettingsGet(category);
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
    if (!setting || setting.value === undefined) return defaultValue;

    // 根据值类型解析
    switch (setting.value_type) {
      case "boolean":
        return (setting.value === "true") as T;
      case "number":
        return Number(setting.value) as T;
      case "json":
        try {
          return JSON.parse(setting.value) as T;
        } catch {
          return defaultValue;
        }
      default:
        return setting.value as T;
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
  const updateSetting = async (key: string, value: string): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";

    try {
      const data: HandlerUpdateSettingRequest = { value };
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

  // 批量更新设置
  const batchUpdateSettings = async (
    updates: { key: string; value: string | number | boolean | object }[],
  ): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      // 构建批量更新请求，将值转换为字符串
      const settingsData = updates.map((u) => ({
        key: u.key,
        value: typeof u.value === "object" ? JSON.stringify(u.value) : String(u.value),
      }));

      const data: HandlerBatchUpdateSettingsRequest = { settings: settingsData };
      await adminSettingsApi.apiAdminSettingsBatchPost(data);

      // 更新本地缓存
      updates.forEach((update) => {
        const index = settings.value.findIndex((s) => s.key === update.key);
        const current = index !== -1 ? settings.value[index] : undefined;
        if (current) {
          current.value = typeof update.value === "object" ? JSON.stringify(update.value) : String(update.value);
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
    settingsByCategory,
    loading,
    saving,
    errorMessage,
    successMessage,
    fetchSettings,
    fetchSettingsByCategory,
    getSetting,
    createSetting,
    updateSetting,
    batchUpdateSettings,
    deleteSetting,
    clearMessages,
  };
}
