/**
 * JSON Logic 验证工具
 *
 * 前后端共享的验证规则引擎，确保前端预校验和后端最终校验行为一致。
 *
 * @see https://jsonlogic.com
 */
import jsonLogic from "json-logic-js";

// 注册自定义操作符（与后端保持一致）

/**
 * strlen - 计算字符串的 UTF-8 字符长度
 *
 * @example
 * {"strlen": {"var": "value"}} // 获取 value 的字符长度
 * {">=": [{"strlen": {"var": "value"}}, 6]} // 长度 >= 6
 */
jsonLogic.add_operation("strlen", (str: unknown): number => {
  if (typeof str !== "string") return 0;
  return [...str].length; // 正确处理 Unicode（中文、emoji 等）
});

/**
 * 验证结果
 */
export interface ValidationResult {
  valid: boolean;
  message?: string;
}

/**
 * 验证设置值
 *
 * @param rule - JSON Logic 规则（字符串或对象）
 * @param value - 要验证的值
 * @param allSettings - 所有设置的当前值（用于跨字段验证）
 * @returns 验证结果
 *
 * @example
 * // 简单验证
 * validateSetting({">=": [{"var": "value"}, 6]}, 10) // { valid: true }
 *
 * // 跨字段验证
 * validateSetting(
 *   {"<=": [{"var": "value"}, {"var": "settings.max_value"}]},
 *   5,
 *   { max_value: 10 }
 * ) // { valid: true }
 */
export function validateSetting(
  rule: string | object,
  value: unknown,
  allSettings?: Record<string, unknown>,
): ValidationResult {
  try {
    const parsedRule = typeof rule === "string" ? JSON.parse(rule) : rule;
    const data = {
      value,
      key: "",
      settings: allSettings || {},
    };

    const result = jsonLogic.apply(parsedRule, data);
    return { valid: Boolean(result) };
  } catch (error) {
    console.error("JSON Logic validation error:", error);
    return { valid: false, message: "验证规则执行失败" };
  }
}

/**
 * 批量验证多个设置
 *
 * @param items - 验证项列表
 * @returns 验证失败的 key -> message 映射
 */
export function validateSettings(
  items: Array<{
    key: string;
    rule: string | object;
    value: unknown;
    message?: string;
  }>,
  allSettings?: Record<string, unknown>,
): Record<string, string> {
  const errors: Record<string, string> = {};

  for (const item of items) {
    const result = validateSetting(item.rule, item.value, allSettings);
    if (!result.valid) {
      errors[item.key] = item.message || result.message || `${item.key} 验证失败`;
    }
  }

  return errors;
}

// 导出 jsonLogic 实例（已注册自定义操作符）
export { jsonLogic };
