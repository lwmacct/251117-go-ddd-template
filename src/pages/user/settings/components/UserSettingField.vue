<script setup lang="ts">
/**
 * 用户设置字段组件
 * 包装 DynamicSettingField，添加重置到默认值功能
 *
 * 支持三种状态：
 * - 可编辑（scope=user）：显示重置按钮
 * - 只读（scope=system）：显示锁定图标，不可编辑
 * - 暂不生效（系统依赖未满足）：可编辑，显示警告锁图标
 */
import { computed } from "vue";
import DynamicSettingField from "@/pages/admin/settings/components/DynamicSettingField.vue";
import type { SettingSettingsItemDTO, SettingSelectOptionDTO } from "@models";

const props = defineProps<{
  setting: SettingSettingsItemDTO;
  modelValue: unknown;
  disabled?: boolean;
  readonly?: boolean; // 只读模式（scope=system 的公开设置）
  systemDependencyInactive?: boolean; // 因系统级依赖而暂不生效（可编辑但无效）
  systemDependencyLabel?: string; // 依赖项的 label（用于 Tooltip）
  showResetButton?: boolean; // 是否显示重置按钮（由父组件判断）
  hint?: string;
  errorMessages?: string[];
  resetting?: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: unknown];
  change: [value: unknown]; // 透传即时变更
  blur: []; // 透传失焦
  reset: [];
}>();

// 格式化默认值用于 tooltip 显示
const formatDefaultValue = computed(() => {
  const defaultVal = props.setting.default_value;
  const valueType = props.setting.value_type;

  if (defaultVal === null || defaultVal === undefined) return "无";

  switch (valueType) {
    case "boolean":
      return defaultVal ? "开启" : "关闭";
    case "select": {
      // 尝试从 options 找到对应的 label
      // default_value 是 JSONB，可能是任意类型，需要转换为 unknown 进行比较
      const option = props.setting.ui_config?.options?.find((o: SettingSelectOptionDTO) => o.value === (defaultVal as unknown));
      return option?.label ?? String(defaultVal);
    }
    default:
      return String(defaultVal);
  }
});

// 将 UserSettingsItemDTO 适配为 SettingSettingsItemDTO 格式
// 两者结构兼容，只是 UserSchema 多了 value 和 is_customized 字段
const adaptedSetting = computed(() => ({
  key: props.setting.key,
  label: props.setting.label,
  value_type: props.setting.value_type,
  input_type: props.setting.input_type,
  ui_config: props.setting.ui_config,
  order: props.setting.order,
  // 使用 value 作为默认值显示（DynamicSettingField 不使用 default_value）
}));
</script>

<template>
  <div class="user-setting-field">
    <div class="d-flex align-start">
      <div class="flex-grow-1">
        <DynamicSettingField
          :setting="adaptedSetting"
          :model-value="modelValue"
          :disabled="disabled || resetting"
          :hint="hint"
          :error-messages="errorMessages"
          @update:model-value="emit('update:modelValue', $event)"
          @change="emit('change', $event)"
          @blur="emit('blur')"
        />
      </div>

      <!-- 只读标识（系统设置） -->
      <v-tooltip v-if="readonly" text="由管理员设置，仅供查看">
        <template #activator="{ props: tooltipProps }">
          <v-btn v-bind="tooltipProps" icon="mdi-lock" size="small" variant="text" color="grey" class="ml-2 mt-3" />
        </template>
      </v-tooltip>

      <!-- 系统依赖暂不生效提示（可配置但当前无效） -->
      <v-tooltip v-else-if="systemDependencyInactive" location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn v-bind="tooltipProps" icon="mdi-lock-outline" size="small" variant="text" color="warning" class="ml-2 mt-3" />
        </template>
        <div class="text-center">需要管理员开启「{{ systemDependencyLabel }}」后才会生效</div>
      </v-tooltip>

      <!-- 重置按钮（仅可编辑且已自定义时显示） -->
      <v-tooltip v-else-if="showResetButton" :text="`重置为默认值: ${formatDefaultValue}`">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            icon="mdi-restore"
            size="small"
            variant="text"
            color="info"
            class="ml-2 mt-3"
            :loading="resetting"
            :disabled="disabled"
            @click="emit('reset')"
          ></v-btn>
        </template>
      </v-tooltip>
    </div>
  </div>
</template>

<style scoped>
.user-setting-field {
  width: 100%;
}
</style>
