<script setup lang="ts">
/**
 * 动态设置字段组件
 * 根据 setting.input_type 动态渲染对应的 Vuetify 控件
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
  change: [value: unknown]; // 用于 switch/select 即时保存
  blur: []; // 用于 text 类失焦保存
}>();

// 处理即时变更（switch/select）
const handleChange = (value: unknown) => {
  emit("update:modelValue", value);
  emit("change", value);
};

// 处理失焦（text 类控件）
const handleBlur = () => {
  emit("blur");
};

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

// 控件类型（从顶层字段获取，非 ui_config）
const inputType = computed(() => props.setting.input_type || "text");

// 值类型转换（用于双向绑定 + 即时变更）
const numericValue = computed({
  get: () => Number(props.modelValue) || 0,
  set: (val: number) => emit("update:modelValue", val),
});

const booleanValue = computed({
  get: () => Boolean(props.modelValue),
  set: (val: boolean) => handleChange(val),
});

const stringValue = computed({
  get: () => String(props.modelValue ?? ""),
  set: (val: string) => emit("update:modelValue", val),
});

// Select 专用（需要即时触发 change）
const selectValue = computed({
  get: () => String(props.modelValue ?? ""),
  set: (val: string) => handleChange(val),
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
      v-model="selectValue"
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
      @blur="handleBlur"
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
      @blur="handleBlur"
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
      @blur="handleBlur"
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
      @blur="handleBlur"
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
      @blur="handleBlur"
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
      @blur="handleBlur"
    />
  </div>
</template>

<style scoped>
.dynamic-setting-field {
  width: 100%;
}
</style>
