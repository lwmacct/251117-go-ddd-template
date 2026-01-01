<!--
  OrgLayout.vue - 组织页面布局
  包含导航栏 + 内容区域
-->
<script setup lang="ts">
import AppBars from "@/views/AppBars/index.vue";
</script>

<template>
  <!-- 导航栏 -->
  <AppBars />

  <!-- 主内容区域 -->
  <v-main>
    <router-view />
  </v-main>
</template>
