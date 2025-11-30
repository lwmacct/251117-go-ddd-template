/**
 * 表单验证工具
 * 提供常用验证规则和验证函数
 */

export type ValidationRule<T = string> = (value: T) => boolean | string;

/**
 * 验证规则工厂函数
 */
export const rules = {
  /**
   * 必填
   */
  required: (message = "此字段为必填项"): ValidationRule => {
    return (value) => {
      if (value === null || value === undefined || value === "") {
        return message;
      }
      if (Array.isArray(value) && value.length === 0) {
        return message;
      }
      return true;
    };
  },

  /**
   * 最小长度
   */
  minLength: (min: number, message?: string): ValidationRule => {
    return (value) => {
      if (!value) return true; // 空值由 required 处理
      if (value.length < min) {
        return message ?? `长度不能少于 ${min} 个字符`;
      }
      return true;
    };
  },

  /**
   * 最大长度
   */
  maxLength: (max: number, message?: string): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (value.length > max) {
        return message ?? `长度不能超过 ${max} 个字符`;
      }
      return true;
    };
  },

  /**
   * 长度范围
   */
  lengthBetween: (min: number, max: number, message?: string): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (value.length < min || value.length > max) {
        return message ?? `长度必须在 ${min} 到 ${max} 个字符之间`;
      }
      return true;
    };
  },

  /**
   * 邮箱格式
   */
  email: (message = "请输入有效的邮箱地址"): ValidationRule => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return (value) => {
      if (!value) return true;
      if (!emailRegex.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * 手机号格式（中国大陆）
   */
  phone: (message = "请输入有效的手机号"): ValidationRule => {
    const phoneRegex = /^1[3-9]\d{9}$/;
    return (value) => {
      if (!value) return true;
      if (!phoneRegex.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * URL 格式
   */
  url: (message = "请输入有效的网址"): ValidationRule => {
    return (value) => {
      if (!value) return true;
      try {
        new URL(value);
        return true;
      } catch {
        return message;
      }
    };
  },

  /**
   * 数字
   */
  number: (message = "请输入有效的数字"): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (isNaN(Number(value))) {
        return message;
      }
      return true;
    };
  },

  /**
   * 整数
   */
  integer: (message = "请输入整数"): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (!Number.isInteger(Number(value))) {
        return message;
      }
      return true;
    };
  },

  /**
   * 最小值
   */
  min: (minValue: number, message?: string): ValidationRule<number | string> => {
    return (value) => {
      if (value === null || value === undefined || value === "") return true;
      const num = Number(value);
      if (isNaN(num) || num < minValue) {
        return message ?? `不能小于 ${minValue}`;
      }
      return true;
    };
  },

  /**
   * 最大值
   */
  max: (maxValue: number, message?: string): ValidationRule<number | string> => {
    return (value) => {
      if (value === null || value === undefined || value === "") return true;
      const num = Number(value);
      if (isNaN(num) || num > maxValue) {
        return message ?? `不能大于 ${maxValue}`;
      }
      return true;
    };
  },

  /**
   * 正则表达式
   */
  pattern: (regex: RegExp, message = "格式不正确"): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (!regex.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * 用户名格式（字母数字下划线，3-20 位）
   */
  username: (message = "用户名只能包含字母、数字和下划线，长度 3-20 位"): ValidationRule => {
    const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
    return (value) => {
      if (!value) return true;
      if (!usernameRegex.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * 密码强度（至少包含大小写字母和数字，8-32 位）
   */
  password: (message = "密码需包含大小写字母和数字，长度 8-32 位"): ValidationRule => {
    return (value) => {
      if (!value) return true;
      if (value.length < 8 || value.length > 32) {
        return message;
      }
      if (!/[a-z]/.test(value) || !/[A-Z]/.test(value) || !/\d/.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * 确认密码匹配
   */
  sameAs: <T>(otherValue: () => T, message = "两次输入不一致"): ValidationRule<T> => {
    return (value) => {
      if (!value) return true;
      if (value !== otherValue()) {
        return message;
      }
      return true;
    };
  },

  /**
   * 中文字符
   */
  chinese: (message = "请输入中文"): ValidationRule => {
    const chineseRegex = /^[\u4e00-\u9fa5]+$/;
    return (value) => {
      if (!value) return true;
      if (!chineseRegex.test(value)) {
        return message;
      }
      return true;
    };
  },

  /**
   * 身份证号
   */
  idCard: (message = "请输入有效的身份证号"): ValidationRule => {
    const idCardRegex = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/;
    return (value) => {
      if (!value) return true;
      if (!idCardRegex.test(value)) {
        return message;
      }
      return true;
    };
  },
};

/**
 * 组合多个验证规则
 */
export function composeRules<T = string>(...ruleList: ValidationRule<T>[]): ValidationRule<T>[] {
  return ruleList;
}

/**
 * 验证单个值
 */
export function validate<T = string>(value: T, ruleList: ValidationRule<T>[]): string | null {
  for (const rule of ruleList) {
    const result = rule(value);
    if (result !== true && typeof result === "string") {
      return result;
    }
  }
  return null;
}

/**
 * 验证对象
 */
export function validateObject<T extends Record<string, unknown>>(
  data: T,
  schema: Partial<Record<keyof T, ValidationRule<T[keyof T]>[]>>
): Record<keyof T, string | null> {
  const errors = {} as Record<keyof T, string | null>;

  for (const key of Object.keys(schema) as (keyof T)[]) {
    const ruleList = schema[key];
    if (ruleList) {
      errors[key] = validate(data[key] as T[keyof T], ruleList);
    }
  }

  return errors;
}
