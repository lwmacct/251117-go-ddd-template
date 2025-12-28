<script setup lang="ts">
/**
 * 收藏按钮组件
 *
 * 消除 MenuListSection、RecentPagesPanel、ProductsPanel 中的重复代码
 */

interface Props {
  /** 路径标识 */
  path: string;
  /** 是否已收藏 */
  isFavorite?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  isFavorite: false,
});

const emit = defineEmits<{
  toggle: [path: string];
}>();

function handleClick(event: Event) {
  event.stopPropagation();
  emit("toggle", props.path);
}
</script>

<template>
  <v-btn
    :icon="isFavorite ? 'mdi-star' : 'mdi-star-outline'"
    :color="isFavorite ? 'warning' : undefined"
    size="x-small"
    variant="text"
    @click="handleClick"
  />
</template>
