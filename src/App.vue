<script setup lang="ts">
// App.vue - 应用根组件
// 使用 router-view 渲染路由匹配的组件
</script>

<template>
  <v-app>
    <!-- 路由加载进度条 -->
    <LoadingBar />
    <!-- 路由视图 - 渲染匹配的路由组件 -->
    <router-view />
  </v-app>
</template>
