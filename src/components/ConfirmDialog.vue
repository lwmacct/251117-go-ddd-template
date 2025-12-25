<script setup lang="ts">
import { computed } from "vue";

/**
 * 可复用确认对话框组件
 * 支持多种类型：删除、警告、信息确认
 */

interface Props {
  /** 对话框显示状态 */
  modelValue: boolean;
  /** 对话框标题 */
  title?: string;
  /** 确认消息内容 */
  message?: string;
  /** 对话框类型 */
  type?: "delete" | "warning" | "info";
  /** 确认按钮文本 */
  confirmText?: string;
  /** 取消按钮文本 */
  cancelText?: string;
  /** 是否显示加载状态 */
  loading?: boolean;
  /** 最大宽度 */
  maxWidth?: string | number;
  /** 是否持久化（点击遮罩不关闭） */
  persistent?: boolean;
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "confirm"): void;
  (e: "cancel"): void;
}

const props = withDefaults(defineProps<Props>(), {
  title: "确认操作",
  message: "确定要执行此操作吗？",
  type: "info",
  confirmText: "确认",
  cancelText: "取消",
  loading: false,
  maxWidth: 450,
  persistent: false,
});

const emit = defineEmits<Emits>();

// 类型配置
const typeConfig = computed(() => {
  switch (props.type) {
    case "delete":
      return {
        icon: "mdi-delete-alert",
        iconColor: "error",
        confirmColor: "error",
        title: props.title || "确认删除",
        confirmText: props.confirmText || "删除",
      };
    case "warning":
      return {
        icon: "mdi-alert",
        iconColor: "warning",
        confirmColor: "warning",
        title: props.title || "警告",
        confirmText: props.confirmText || "继续",
      };
    default:
      return {
        icon: "mdi-help-circle",
        iconColor: "primary",
        confirmColor: "primary",
        title: props.title || "确认操作",
        confirmText: props.confirmText || "确认",
      };
  }
});

const handleClose = () => {
  emit("update:modelValue", false);
  emit("cancel");
};

const handleConfirm = () => {
  emit("confirm");
};
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    :max-width="maxWidth"
    :persistent="persistent || loading"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon :color="typeConfig.iconColor" class="mr-2">{{ typeConfig.icon }}</v-icon>
        {{ typeConfig.title }}
      </v-card-title>

      <v-card-text>
        <slot>
          <span v-html="message"></span>
        </slot>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" :disabled="loading" @click="handleClose">
          {{ cancelText }}
        </v-btn>
        <v-btn :color="typeConfig.confirmColor" variant="elevated" :loading="loading" @click="handleConfirm">
          {{ typeConfig.confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
