<script setup lang="ts">
/**
 * 加载遮罩组件
 * 用于显示全屏或区域加载状态
 */

interface Props {
  /** 是否显示加载 */
  loading?: boolean;
  /** 加载提示文本 */
  text?: string;
  /** 是否绝对定位（覆盖父容器） */
  absolute?: boolean;
  /** 遮罩透明度 */
  opacity?: number;
  /** 加载指示器大小 */
  size?: string | number;
  /** 加载指示器颜色 */
  color?: string;
  /** 遮罩颜色 */
  overlayColor?: string;
}

withDefaults(defineProps<Props>(), {
  loading: false,
  text: "",
  absolute: true,
  opacity: 0.8,
  size: 48,
  color: "primary",
  overlayColor: "white",
});
</script>

<template>
  <div class="loading-overlay-container" :class="{ 'position-relative': absolute }">
    <!-- 内容插槽 -->
    <slot />

    <!-- 加载遮罩 -->
    <Transition name="fade">
      <div
        v-if="loading"
        class="loading-overlay d-flex flex-column align-center justify-center"
        :class="{ absolute: absolute }"
        :style="{
          '--overlay-opacity': opacity,
          '--overlay-color': `rgb(var(--v-theme-${overlayColor}))`,
        }"
      >
        <v-progress-circular :size="size" :color="color" indeterminate />
        <span v-if="text" class="text-body-2 text-medium-emphasis mt-3">
          {{ text }}
        </span>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.loading-overlay-container {
  width: 100%;
  height: 100%;
}

.loading-overlay {
  z-index: 10;
  background-color: var(--overlay-color);
  opacity: var(--overlay-opacity);
}

.loading-overlay.absolute {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
