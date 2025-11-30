/**
 * 表单状态 Composable
 * 管理表单的脏状态、提交、重置等
 */

import { ref, reactive, computed, watch, type Ref, type UnwrapRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseFormOptions<T extends Record<string, unknown>> {
  /** 初始值 */
  initialValues: T;
  /** 提交处理函数 */
  onSubmit?: (values: T) => void | Promise<void>;
  /** 验证函数 */
  validate?: (values: T) => Record<keyof T, string> | null;
  /** 字段变化回调 */
  onChange?: (name: keyof T, value: T[keyof T]) => void;
  /** 提交成功回调 */
  onSuccess?: () => void;
  /** 提交失败回调 */
  onError?: (error: Error) => void;
}

export interface UseFormReturn<T extends Record<string, unknown>> {
  /** 表单值 */
  values: T;
  /** 初始值 */
  initialValues: T;
  /** 错误信息 */
  errors: Partial<Record<keyof T, string>>;
  /** 触摸状态 */
  touched: Partial<Record<keyof T, boolean>>;
  /** 是否已修改 */
  isDirty: Ref<boolean>;
  /** 是否有效 */
  isValid: Ref<boolean>;
  /** 是否正在提交 */
  isSubmitting: Ref<boolean>;
  /** 提交次数 */
  submitCount: Ref<number>;
  /** 设置字段值 */
  setFieldValue: <K extends keyof T>(name: K, value: T[K]) => void;
  /** 设置字段错误 */
  setFieldError: (name: keyof T, error: string) => void;
  /** 设置字段触摸 */
  setFieldTouched: (name: keyof T, touched?: boolean) => void;
  /** 设置所有值 */
  setValues: (values: Partial<T>) => void;
  /** 设置所有错误 */
  setErrors: (errors: Partial<Record<keyof T, string>>) => void;
  /** 重置表单 */
  reset: (newInitialValues?: T) => void;
  /** 提交表单 */
  submit: () => Promise<void>;
  /** 验证表单 */
  validateForm: () => boolean;
  /** 获取字段属性 */
  getFieldProps: <K extends keyof T>(name: K) => FieldProps<T[K]>;
}

export interface FieldProps<V> {
  modelValue: V;
  "onUpdate:modelValue": (value: V) => void;
  onBlur: () => void;
  error: boolean;
  errorMessages: string[];
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 表单状态管理
 * @example
 * const form = useForm({
 *   initialValues: { name: '', email: '' },
 *   onSubmit: async (values) => {
 *     await api.createUser(values)
 *   },
 *   validate: (values) => {
 *     const errors: Record<string, string> = {}
 *     if (!values.name) errors.name = '名称必填'
 *     if (!values.email) errors.email = '邮箱必填'
 *     return Object.keys(errors).length > 0 ? errors : null
 *   }
 * })
 */
export function useForm<T extends Record<string, unknown>>(options: UseFormOptions<T>): UseFormReturn<T> {
  const { initialValues: initialValuesOption, onSubmit, validate, onChange, onSuccess, onError } = options;

  // 存储初始值的副本
  const initialValues = reactive({ ...initialValuesOption }) as T;

  // 表单值
  const values = reactive({ ...initialValuesOption }) as T;

  // 错误信息
  const errors = reactive<Partial<Record<keyof T, string>>>({});

  // 触摸状态
  const touched = reactive<Partial<Record<keyof T, boolean>>>({});

  // 提交状态
  const isSubmitting = ref(false);
  const submitCount = ref(0);

  // 是否已修改
  const isDirty = computed(() => {
    return Object.keys(values).some((key) => {
      const k = key as keyof T;
      return JSON.stringify(values[k]) !== JSON.stringify(initialValues[k]);
    });
  });

  // 是否有效
  const isValid = computed(() => {
    return Object.keys(errors).length === 0;
  });

  // 设置字段值
  const setFieldValue = <K extends keyof T>(name: K, value: T[K]) => {
    (values as Record<K, T[K]>)[name] = value;
    onChange?.(name, value);
  };

  // 设置字段错误
  const setFieldError = (name: keyof T, error: string) => {
    if (error) {
      errors[name] = error;
    } else {
      delete errors[name];
    }
  };

  // 设置字段触摸
  const setFieldTouched = (name: keyof T, isTouched = true) => {
    touched[name] = isTouched;
  };

  // 设置所有值
  const setValues = (newValues: Partial<T>) => {
    Object.assign(values, newValues);
  };

  // 设置所有错误
  const setErrors = (newErrors: Partial<Record<keyof T, string>>) => {
    // 清除现有错误
    Object.keys(errors).forEach((key) => {
      delete errors[key as keyof T];
    });
    // 设置新错误
    Object.assign(errors, newErrors);
  };

  // 验证表单
  const validateForm = (): boolean => {
    if (!validate) return true;

    const validationErrors = validate(values as T);
    if (validationErrors) {
      setErrors(validationErrors);
      return false;
    }

    setErrors({});
    return true;
  };

  // 重置表单
  const reset = (newInitialValues?: T) => {
    if (newInitialValues) {
      Object.assign(initialValues, newInitialValues);
    }

    Object.assign(values, initialValues);

    // 清除错误和触摸状态
    Object.keys(errors).forEach((key) => {
      delete errors[key as keyof T];
    });
    Object.keys(touched).forEach((key) => {
      delete touched[key as keyof T];
    });

    submitCount.value = 0;
  };

  // 提交表单
  const submit = async (): Promise<void> => {
    submitCount.value++;

    // 标记所有字段为已触摸
    Object.keys(values).forEach((key) => {
      touched[key as keyof T] = true;
    });

    // 验证
    if (!validateForm()) {
      return;
    }

    if (!onSubmit) return;

    isSubmitting.value = true;

    try {
      await onSubmit(values as T);
      onSuccess?.();
    } catch (e) {
      const error = e instanceof Error ? e : new Error(String(e));
      onError?.(error);
      throw error;
    } finally {
      isSubmitting.value = false;
    }
  };

  // 获取字段属性（用于 v-bind）
  const getFieldProps = <K extends keyof T>(name: K): FieldProps<T[K]> => ({
    modelValue: values[name],
    "onUpdate:modelValue": (value: T[K]) => setFieldValue(name, value),
    onBlur: () => setFieldTouched(name),
    error: !!(touched[name] && errors[name]),
    errorMessages: touched[name] && errors[name] ? [errors[name]!] : [],
  });

  return {
    values: values as T,
    initialValues: initialValues as T,
    errors,
    touched,
    isDirty,
    isValid,
    isSubmitting,
    submitCount,
    setFieldValue,
    setFieldError,
    setFieldTouched,
    setValues,
    setErrors,
    reset,
    submit,
    validateForm,
    getFieldProps,
  };
}

// ============================================================================
// 表单脏状态警告
// ============================================================================

export interface UseFormDirtyGuardOptions {
  /** 是否脏 */
  isDirty: Ref<boolean>;
  /** 警告消息 */
  message?: string;
  /** 离开前回调（返回 true 允许离开） */
  onBeforeLeave?: () => boolean | Promise<boolean>;
}

/**
 * 表单脏状态离开警告
 * @example
 * const { isDirty } = useForm({ ... })
 * useFormDirtyGuard({ isDirty })
 */
export function useFormDirtyGuard(options: UseFormDirtyGuardOptions) {
  const { isDirty, message = "您有未保存的更改，确定要离开吗？", onBeforeLeave } = options;

  // 浏览器关闭/刷新警告
  const handleBeforeUnload = (e: BeforeUnloadEvent) => {
    if (isDirty.value) {
      e.preventDefault();
      e.returnValue = message;
      return message;
    }
  };

  // 设置/移除事件监听
  watch(
    isDirty,
    (dirty) => {
      if (dirty) {
        window.addEventListener("beforeunload", handleBeforeUnload);
      } else {
        window.removeEventListener("beforeunload", handleBeforeUnload);
      }
    },
    { immediate: true }
  );

  // 检查是否可以离开
  const canLeave = async (): Promise<boolean> => {
    if (!isDirty.value) return true;

    if (onBeforeLeave) {
      return await onBeforeLeave();
    }

    return window.confirm(message);
  };

  return {
    canLeave,
  };
}

// ============================================================================
// 字段数组
// ============================================================================

/**
 * 字段数组（用于动态表单字段）
 * @example
 * const { fields, append, remove, move } = useFieldArray<string>([])
 */
export function useFieldArray<T>(initialValue: T[] = []) {
  const fields = ref<T[]>([...initialValue]) as Ref<T[]>;

  // 添加项
  const append = (value: T) => {
    fields.value.push(value);
  };

  // 前置添加
  const prepend = (value: T) => {
    fields.value.unshift(value);
  };

  // 插入项
  const insert = (index: number, value: T) => {
    fields.value.splice(index, 0, value);
  };

  // 移除项
  const remove = (index: number) => {
    fields.value.splice(index, 1);
  };

  // 移动项
  const move = (fromIndex: number, toIndex: number) => {
    const item = fields.value[fromIndex];
    fields.value.splice(fromIndex, 1);
    fields.value.splice(toIndex, 0, item);
  };

  // 交换项
  const swap = (indexA: number, indexB: number) => {
    const temp = fields.value[indexA];
    fields.value[indexA] = fields.value[indexB];
    fields.value[indexB] = temp;
  };

  // 替换项
  const replace = (index: number, value: T) => {
    fields.value[index] = value;
  };

  // 更新项
  const update = (index: number, updater: (item: T) => T) => {
    fields.value[index] = updater(fields.value[index]);
  };

  // 清空
  const clear = () => {
    fields.value = [];
  };

  // 重置
  const reset = (newValue?: T[]) => {
    fields.value = newValue ? [...newValue] : [...initialValue];
  };

  return {
    fields,
    append,
    prepend,
    insert,
    remove,
    move,
    swap,
    replace,
    update,
    clear,
    reset,
  };
}
