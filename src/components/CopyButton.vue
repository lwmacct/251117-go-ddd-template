<script setup lang="ts">
import { useClipboard } from "@vueuse/core";

/**
 * 复制按钮组件
 * 点击后复制指定文本到剪贴板，并显示成功反馈
 */

interface Props {
  /** 要复制的文本 */
  text: string;
  /** 按钮大小 */
  size?: "x-small" | "small" | "default" | "large" | "x-large";
  /** 是否仅显示图标 */
  iconOnly?: boolean;
  /** 成功提示文本 */
  successText?: string;
  /** 按钮颜色 */
  color?: string;
}

const props = withDefaults(defineProps<Props>(), {
  size: "small",
  iconOnly: true,
  successText: "已复制",
  color: undefined,
});

const { copied, copy } = useClipboard();

const handleClick = () => {
  copy(props.text);
};
</script>

<template>
  <v-tooltip :text="copied ? successText : '复制'" location="top">
    <template #activator="{ props: tooltipProps }">
      <v-btn
        v-bind="tooltipProps"
        :icon="iconOnly"
        :size="size"
        :color="copied ? 'success' : color"
        variant="text"
        @click="handleClick"
      >
        <v-icon :icon="copied ? 'mdi-check' : 'mdi-content-copy'" />
        <span v-if="!iconOnly" class="ml-1">{{ copied ? successText : "复制" }}</span>
      </v-btn>
    </template>
  </v-tooltip>
</template>
