/**
 * 设置项依赖关系处理 Composable
 */
import type { Ref } from "vue";
import type { SettingSettingsItemDTO } from "@models";

// 依赖关系配置
interface DependsOnConfig {
  key: string;
  value?: unknown;
  operator?: "eq" | "ne" | "gt" | "lt";
}

export function useSettingsDependency(
  settingsMap: Ref<Map<string, SettingSettingsItemDTO>>,
  formValues: Ref<Record<string, unknown>>,
) {
  /**
   * 解析依赖关系配置
   */
  const parseDependsOn = (setting: SettingSettingsItemDTO): DependsOnConfig | null => {
    const dependsOn = setting.ui_config?.depends_on;
    if (!dependsOn) return null;
    try {
      return dependsOn as DependsOnConfig;
    } catch {
      return null;
    }
  };

  /**
   * 检查设置项是否因依赖关系被禁用
   */
  const isDisabled = (setting: SettingSettingsItemDTO): boolean => {
    const dep = parseDependsOn(setting);
    if (!dep) return false;

    const depValue = formValues.value[dep.key];
    const expectedValue = dep.value;
    const operator = dep.operator || "eq";

    switch (operator) {
      case "eq":
        return depValue !== expectedValue;
      case "ne":
        return depValue === expectedValue;
      case "gt":
        return !(Number(depValue) > Number(expectedValue));
      case "lt":
        return !(Number(depValue) < Number(expectedValue));
      default:
        return false;
    }
  };

  /**
   * 获取禁用原因提示
   */
  const getDisabledHint = (setting: SettingSettingsItemDTO): string | undefined => {
    const dep = parseDependsOn(setting);
    if (!dep || !isDisabled(setting)) return undefined;

    const depSetting = settingsMap.value.get(dep.key);
    return `需要先启用「${depSetting?.label || dep.key}」`;
  };

  /**
   * 获取设置项的最终 hint（优先显示禁用原因）
   */
  const getFinalHint = (setting: SettingSettingsItemDTO): string | undefined => {
    return getDisabledHint(setting) || setting.ui_config?.hint;
  };

  return {
    isDisabled,
    getDisabledHint,
    getFinalHint,
  };
}
