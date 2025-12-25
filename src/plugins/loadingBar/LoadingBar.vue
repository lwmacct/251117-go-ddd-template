<!--
  LoadingBar.vue - 路由加载进度条组件
  在页面顶部显示类似 GitHub 的线性进度条
-->
<script setup lang="ts">
import { useLoadingBar } from "./composable";

const { isLoading, progress, color } = useLoadingBar();
</script>

<template>
  <v-progress-linear :active="isLoading" :model-value="progress" :color="color" :height="3" absolute style="z-index: 9999" />
</template>
