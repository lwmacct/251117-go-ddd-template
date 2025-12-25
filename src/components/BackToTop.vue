<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useEventListener } from "@vueuse/core";

/**
 * 回到顶部按钮组件
 * 当页面滚动超过一定距离时显示
 */

interface Props {
  /** 显示阈值（滚动距离） */
  threshold?: number;
  /** 滚动目标元素选择器，默认为 window */
  target?: string;
  /** 按钮位置 */
  position?: "bottom-right" | "bottom-left";
  /** 距离边缘的距离 */
  offset?: number;
  /** 滚动动画时长（毫秒） */
  duration?: number;
}

const props = withDefaults(defineProps<Props>(), {
  threshold: 200,
  target: "",
  position: "bottom-right",
  offset: 24,
  duration: 300,
});

const visible = ref(false);
const targetElement = ref<Element | Window | null>(null);

/**
 * 处理滚动事件
 */
const handleScroll = () => {
  if (targetElement.value === window) {
    visible.value = window.scrollY > props.threshold;
  } else if (targetElement.value instanceof Element) {
    visible.value = targetElement.value.scrollTop > props.threshold;
  }
};

/**
 * 滚动到顶部
 */
const scrollToTop = () => {
  if (targetElement.value === window) {
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  } else if (targetElement.value instanceof Element) {
    targetElement.value.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }
};

onMounted(() => {
  if (props.target) {
    targetElement.value = document.querySelector(props.target);
  } else {
    targetElement.value = window;
  }

  // 使用 useEventListener 自动管理事件监听器生命周期
  if (targetElement.value) {
    useEventListener(targetElement.value, "scroll", handleScroll, { passive: true });
  }
});

/**
 * 计算按钮位置样式
 */
const positionStyle = computed(() => ({
  [props.position === "bottom-right" ? "right" : "left"]: `${props.offset}px`,
  bottom: `${props.offset}px`,
}));
</script>

<template>
  <Transition name="fade-slide">
    <v-btn
      v-show="visible"
      icon="mdi-chevron-up"
      color="primary"
      class="back-to-top"
      :style="positionStyle"
      elevation="4"
      size="large"
      @click="scrollToTop"
    >
      <v-tooltip activator="parent" location="left">回到顶部</v-tooltip>
    </v-btn>
  </Transition>
</template>

<style scoped>
.back-to-top {
  position: fixed;
  z-index: 1000;
}

/* 过渡动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
