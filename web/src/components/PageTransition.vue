<script setup lang="ts">
/**
 * 页面过渡动画组件
 * 用于 router-view 的过渡效果
 */

interface Props {
  /** 过渡动画名称 */
  name?: "fade" | "slide-left" | "slide-right" | "slide-up" | "slide-down" | "scale" | "none";
  /** 过渡模式 */
  mode?: "default" | "out-in" | "in-out";
  /** 是否启用 */
  enabled?: boolean;
}

withDefaults(defineProps<Props>(), {
  name: "fade",
  mode: "out-in",
  enabled: true,
});
</script>

<template>
  <Transition v-if="enabled" :name="`page-${name}`" :mode="mode">
    <slot />
  </Transition>
  <slot v-else />
</template>

<style>
/* 淡入淡出 */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

/* 左滑 */
.page-slide-left-enter-active,
.page-slide-left-leave-active {
  transition: all 0.25s ease-out;
}

.page-slide-left-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-slide-left-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 右滑 */
.page-slide-right-enter-active,
.page-slide-right-leave-active {
  transition: all 0.25s ease-out;
}

.page-slide-right-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.page-slide-right-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* 上滑 */
.page-slide-up-enter-active,
.page-slide-up-leave-active {
  transition: all 0.25s ease-out;
}

.page-slide-up-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-slide-up-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 下滑 */
.page-slide-down-enter-active,
.page-slide-down-leave-active {
  transition: all 0.25s ease-out;
}

.page-slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.page-slide-down-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* 缩放 */
.page-scale-enter-active,
.page-scale-leave-active {
  transition: all 0.2s ease;
}

.page-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.page-scale-leave-to {
  opacity: 0;
  transform: scale(1.05);
}

/* 无动画 */
.page-none-enter-active,
.page-none-leave-active {
  transition: none;
}
</style>
