<script setup lang="ts">
/**
 * 动态设置字段组件
 * 根据 setting.ui_config.input_type 动态渲染对应的 Vuetify 控件
 */
import { computed } from "vue";
import type { SettingSchemaSettingDTO } from "@models";

// 选项配置
interface SelectOption {
  value: string;
  label: string;
}

const props = defineProps<{
  setting: SettingSchemaSettingDTO;
  modelValue: unknown;
  disabled?: boolean;
  hint?: string;
  errorMessages?: string[];
}>();

const emit = defineEmits<{
  "update:modelValue": [value: unknown];
}>();

// 获取 UI 配置
const uiConfig = computed(() => props.setting.ui_config || {});

// 解析选项配置
const options = computed<SelectOption[]>(() => {
  const items = uiConfig.value.options;
  if (!items || !Array.isArray(items)) return [];
  return items
    .filter((item): item is { value: string; label?: string } => typeof item.value === "string")
    .map((item) => ({ value: item.value, label: item.label || item.value }));
});

// 转换为 Vuetify select items 格式
const selectItems = computed(() =>
  options.value.map((opt) => ({
    title: opt.label,
    value: opt.value,
  })),
);

// 控件类型
const inputType = computed(() => uiConfig.value.input_type || "text");

// 值类型转换
const numericValue = computed({
  get: () => Number(props.modelValue) || 0,
  set: (val: number) => emit("update:modelValue", val),
});

const booleanValue = computed({
  get: () => Boolean(props.modelValue),
  set: (val: boolean) => emit("update:modelValue", val),
});

const stringValue = computed({
  get: () => String(props.modelValue ?? ""),
  set: (val: string) => emit("update:modelValue", val),
});

// 最终显示的 hint（外部传入优先，否则使用 ui_config.hint）
const finalHint = computed(() => props.hint || uiConfig.value.hint);
</script>

<template>
  <div class="dynamic-setting-field">
    <!-- Switch 开关 -->
    <v-switch
      v-if="inputType === 'switch'"
      v-model="booleanValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      color="primary"
      persistent-hint
      class="mb-2"
    />

    <!-- Select 下拉选择 -->
    <v-select
      v-else-if="inputType === 'select'"
      v-model="stringValue"
      :label="setting.label"
      :items="selectItems"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      variant="outlined"
      persistent-hint
    />

    <!-- Number 数字输入 -->
    <v-text-field
      v-else-if="inputType === 'number'"
      v-model.number="numericValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      type="number"
      variant="outlined"
      persistent-hint
    />

    <!-- Textarea 多行文本 -->
    <v-textarea
      v-else-if="inputType === 'textarea'"
      v-model="stringValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      variant="outlined"
      persistent-hint
      rows="3"
    />

    <!-- Password 密码 -->
    <v-text-field
      v-else-if="inputType === 'password'"
      v-model="stringValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      type="password"
      variant="outlined"
      persistent-hint
    />

    <!-- Email 邮箱 -->
    <v-text-field
      v-else-if="inputType === 'email'"
      v-model="stringValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      type="email"
      variant="outlined"
      persistent-hint
    />

    <!-- URL -->
    <v-text-field
      v-else-if="inputType === 'url'"
      v-model="stringValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      type="url"
      variant="outlined"
      persistent-hint
      placeholder="https://example.com"
    />

    <!-- 默认 Text 文本输入 -->
    <v-text-field
      v-else
      v-model="stringValue"
      :label="setting.label"
      :hint="finalHint"
      :disabled="disabled"
      :error-messages="errorMessages"
      variant="outlined"
      persistent-hint
    />
  </div>
</template>

<style scoped>
.dynamic-setting-field {
  width: 100%;
}
</style>
