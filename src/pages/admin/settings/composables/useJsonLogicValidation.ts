/**
 * JSON Logic 验证 Composable
 * 与后端共用同一套验证规则，确保前后端验证一致
 */
import { ref, type Ref } from "vue";
import jsonLogic from "json-logic-js";
import type { SettingSchemaCategoryDTO, SettingSchemaSettingDTO } from "@models";

// 简单验证规则格式（向后兼容）
interface SimpleValidationRule {
  min?: number;
  max?: number;
  required?: boolean;
  pattern?: string;
  message?: string;
}

export function useJsonLogicValidation(schema: Ref<SettingSchemaCategoryDTO[]>, formValues: Ref<Record<string, unknown>>) {
  // 验证错误 Map: key -> error message
  const errors = ref<Map<string, string>>(new Map());

  /**
   * 检查规则是否为 JSON Logic 格式
   */
  const isJsonLogicRule = (rule: unknown): boolean => {
    if (typeof rule !== "object" || rule === null) return false;
    const keys = Object.keys(rule as object);
    // JSON Logic 规则的顶层 key 是操作符
    const jsonLogicOps = [
      "and",
      "or",
      "!",
      "!!",
      "if",
      "<",
      "<=",
      ">",
      ">=",
      "==",
      "===",
      "!=",
      "!==",
      "+",
      "-",
      "*",
      "/",
      "%",
      "var",
      "missing",
      "missing_some",
      "merge",
      "in",
      "cat",
      "substr",
      "log",
      "min",
      "max",
      "strlen",
    ];
    return keys.some((k) => jsonLogicOps.includes(k));
  };

  /**
   * 将简单规则转换为 JSON Logic 格式
   */
  const convertSimpleRule = (rule: SimpleValidationRule): object => {
    const conditions: object[] = [];

    if (rule.required) {
      conditions.push({ "!!": { var: "value" } });
    }

    if (rule.min !== undefined && rule.max !== undefined) {
      conditions.push({
        and: [{ ">=": [{ var: "value" }, rule.min] }, { "<=": [{ var: "value" }, rule.max] }],
      });
    } else if (rule.min !== undefined) {
      conditions.push({ ">=": [{ var: "value" }, rule.min] });
    } else if (rule.max !== undefined) {
      conditions.push({ "<=": [{ var: "value" }, rule.max] });
    }

    if (rule.pattern) {
      // JSON Logic 不直接支持 regex，需要自定义或跳过
      // 暂时跳过 pattern 验证
    }

    if (conditions.length === 0) return { "!!": [true] };
    if (conditions.length === 1) return conditions[0]!;
    return { and: conditions };
  };

  /**
   * 生成默认错误消息
   */
  const generateErrorMessage = (rule: unknown, label?: string): string => {
    const fieldName = label || "此字段";

    if (typeof rule !== "object" || rule === null) {
      return `${fieldName}验证失败`;
    }

    // 处理简单规则的消息
    const simple = rule as SimpleValidationRule;
    if (simple.message) return simple.message;

    if (simple.min !== undefined && simple.max !== undefined) {
      return `${fieldName}必须在 ${simple.min} 到 ${simple.max} 之间`;
    }
    if (simple.min !== undefined) {
      return `${fieldName}不能小于 ${simple.min}`;
    }
    if (simple.max !== undefined) {
      return `${fieldName}不能大于 ${simple.max}`;
    }
    if (simple.required) {
      return `${fieldName}不能为空`;
    }

    return `${fieldName}验证失败`;
  };

  /**
   * 构建验证数据上下文
   */
  const buildDataContext = (key: string, value: unknown): object => {
    // 构建 settings Map 用于跨字段验证
    const settings: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(formValues.value)) {
      settings[k] = v;
    }

    return {
      value,
      key,
      settings,
    };
  };

  /**
   * 从 Schema 中查找设置项
   */
  const findSettingInSchema = (key: string): SettingSchemaSettingDTO | undefined => {
    for (const cat of schema.value) {
      for (const group of cat.groups || []) {
        for (const setting of group.settings || []) {
          if (setting.key === key) return setting;
        }
      }
    }
    return undefined;
  };

  /**
   * 验证单个字段
   * @returns 错误消息，或 null 如果验证通过
   */
  const validate = (key: string): string | null => {
    const setting = findSettingInSchema(key);
    const validation = setting?.ui_config?.validation;
    if (!setting || !validation) {
      errors.value.delete(key);
      return null;
    }

    const value = formValues.value[key];
    const rule = validation;

    try {
      // 确定使用的规则
      let jsonLogicRule: object;
      if (isJsonLogicRule(rule)) {
        jsonLogicRule = rule as object;
      } else {
        jsonLogicRule = convertSimpleRule(rule as SimpleValidationRule);
      }

      // 构建数据上下文
      const data = buildDataContext(key, value);

      // 执行 JSON Logic
      const result = jsonLogic.apply(jsonLogicRule, data);

      if (!result) {
        const errorMsg = generateErrorMessage(rule, setting.label);
        errors.value.set(key, errorMsg);
        return errorMsg;
      }

      errors.value.delete(key);
      return null;
    } catch (e) {
      console.error(`Validation error for ${key}:`, e);
      const errorMsg = `${setting.label || key}验证规则执行失败`;
      errors.value.set(key, errorMsg);
      return errorMsg;
    }
  };

  /**
   * 验证所有字段
   * @returns 错误 Map，key -> error message
   */
  const validateAll = (): Map<string, string> => {
    errors.value.clear();

    for (const cat of schema.value) {
      for (const group of cat.groups || []) {
        for (const setting of group.settings || []) {
          if (setting.key && setting.ui_config?.validation) {
            validate(setting.key);
          }
        }
      }
    }

    return errors.value;
  };

  /**
   * 获取指定字段的错误
   */
  const getError = (key: string): string | undefined => {
    return errors.value.get(key);
  };

  /**
   * 检查是否有验证错误
   */
  const hasErrors = (): boolean => {
    return errors.value.size > 0;
  };

  /**
   * 清除所有错误
   */
  const clearErrors = (): void => {
    errors.value.clear();
  };

  /**
   * 清除指定字段的错误
   */
  const clearError = (key: string): void => {
    errors.value.delete(key);
  };

  return {
    errors,
    validate,
    validateAll,
    getError,
    hasErrors,
    clearErrors,
    clearError,
  };
}
