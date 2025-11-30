<script setup lang="ts">
/**
 * 空状态组件
 * 用于显示列表为空、搜索无结果等场景
 */

interface Props {
  /** 图标名称（MDI 图标） */
  icon?: string;
  /** 标题文本 */
  title?: string;
  /** 描述文本 */
  description?: string;
  /** 图标大小 */
  iconSize?: string | number;
  /** 图标颜色 */
  iconColor?: string;
  /** 操作按钮文本 */
  actionText?: string;
  /** 是否显示操作按钮 */
  showAction?: boolean;
  /** 预设类型 */
  type?: "empty" | "search" | "error" | "no-permission" | "custom";
}

interface Emits {
  (e: "action"): void;
}

const props = withDefaults(defineProps<Props>(), {
  icon: "",
  title: "",
  description: "",
  iconSize: 80,
  iconColor: "grey-lighten-1",
  actionText: "",
  showAction: false,
  type: "custom",
});

const emit = defineEmits<Emits>();

/**
 * 预设配置
 */
const presets = {
  empty: {
    icon: "mdi-inbox-outline",
    title: "暂无数据",
    description: "当前列表为空",
    iconColor: "grey-lighten-1",
  },
  search: {
    icon: "mdi-file-search-outline",
    title: "未找到结果",
    description: "尝试修改搜索条件或清空筛选",
    iconColor: "grey-lighten-1",
  },
  error: {
    icon: "mdi-alert-circle-outline",
    title: "加载失败",
    description: "数据加载出错，请稍后重试",
    iconColor: "error",
  },
  "no-permission": {
    icon: "mdi-lock-outline",
    title: "无权限访问",
    description: "您没有权限查看此内容",
    iconColor: "warning",
  },
  custom: {
    icon: "mdi-information-outline",
    title: "",
    description: "",
    iconColor: "grey-lighten-1",
  },
};

/**
 * 获取当前显示的配置
 */
const currentIcon = props.icon || presets[props.type].icon;
const currentTitle = props.title || presets[props.type].title;
const currentDescription = props.description || presets[props.type].description;
const currentIconColor = props.iconColor || presets[props.type].iconColor;

const handleAction = () => {
  emit("action");
};
</script>

<template>
  <div class="empty-state d-flex flex-column align-center justify-center py-12">
    <!-- 图标 -->
    <v-icon :icon="currentIcon" :size="iconSize" :color="currentIconColor" class="mb-4" />

    <!-- 标题 -->
    <h3 v-if="currentTitle" class="text-h6 text-medium-emphasis mb-2">
      {{ currentTitle }}
    </h3>

    <!-- 描述 -->
    <p v-if="currentDescription" class="text-body-2 text-disabled text-center" style="max-width: 300px">
      {{ currentDescription }}
    </p>

    <!-- 自定义内容插槽 -->
    <slot />

    <!-- 操作按钮 -->
    <v-btn v-if="showAction && actionText" variant="outlined" color="primary" class="mt-4" @click="handleAction">
      {{ actionText }}
    </v-btn>

    <!-- 操作按钮插槽 -->
    <div v-if="$slots.action" class="mt-4">
      <slot name="action" />
    </div>
  </div>
</template>

<style scoped>
.empty-state {
  min-height: 200px;
}
</style>
