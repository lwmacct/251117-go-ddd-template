/**
 * Admin 系统设置 Composable
 */
import { ref, computed } from "vue";
import { AdminSettingsAPI } from "@/api/admin";
import type { Setting, CreateSettingRequest } from "@/types/admin";

export function useSettings() {
  const settings = ref<Setting[]>([]);
  const loading = ref(false);
  const saving = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  // 按分类缓存的设置
  const settingsByCategory = computed(() => {
    const map = new Map<string, Map<string, Setting>>();
    settings.value.forEach((setting) => {
      if (!map.has(setting.category)) {
        map.set(setting.category, new Map());
      }
      map.get(setting.category)!.set(setting.key, setting);
    });
    return map;
  });

  // 获取所有设置
  const fetchSettings = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      settings.value = await AdminSettingsAPI.getSettings();
    } catch (error: any) {
      errorMessage.value = error.message || "获取设置失败";
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
      const categorySettings = await AdminSettingsAPI.getSettingsByCategory(category);
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
    } catch (error: any) {
      errorMessage.value = error.message || "获取设置失败";
      console.error("Failed to fetch settings by category:", error);
    } finally {
      loading.value = false;
    }
  };

  // 获取单个设置的值
  const getSetting = (key: string, defaultValue: any = ""): any => {
    const setting = settings.value.find((s) => s.key === key);
    if (!setting) return defaultValue;

    // 根据值类型解析
    switch (setting.value_type) {
      case "boolean":
        return setting.value === "true";
      case "number":
        return Number(setting.value);
      case "json":
        try {
          return JSON.parse(setting.value);
        } catch {
          return defaultValue;
        }
      default:
        return setting.value;
    }
  };

  // 创建单个设置
  const createSetting = async (data: CreateSettingRequest): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";

    try {
      const newSetting = await AdminSettingsAPI.createSetting(data);
      settings.value.push(newSetting);
      successMessage.value = "设置创建成功";
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || "创建设置失败";
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
      const updated = await AdminSettingsAPI.updateSetting(key, { value });
      const index = settings.value.findIndex((s) => s.key === key);
      if (index !== -1) {
        settings.value[index] = updated;
      }
      successMessage.value = "设置更新成功";
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || "更新设置失败";
      console.error("Failed to update setting:", error);
      return false;
    } finally {
      saving.value = false;
    }
  };

  // 批量更新设置
  const batchUpdateSettings = async (updates: { key: string; value: any }[]): Promise<boolean> => {
    saving.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      // 构建批量更新请求，将值转换为字符串
      const settingsData = updates.map((u) => ({
        key: u.key,
        value: typeof u.value === "object" ? JSON.stringify(u.value) : String(u.value),
      }));

      await AdminSettingsAPI.batchUpdateSettings({ settings: settingsData });

      // 更新本地缓存
      updates.forEach((update) => {
        const index = settings.value.findIndex((s) => s.key === update.key);
        if (index !== -1) {
          settings.value[index].value = typeof update.value === "object" ? JSON.stringify(update.value) : String(update.value);
        }
      });

      successMessage.value = "设置保存成功";
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || "批量更新设置失败";
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
      await AdminSettingsAPI.deleteSetting(key);
      settings.value = settings.value.filter((s) => s.key !== key);
      successMessage.value = "设置删除成功";
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || "删除设置失败";
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
