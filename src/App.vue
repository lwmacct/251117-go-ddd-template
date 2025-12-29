<script setup lang="ts">
/**
 * App.vue - 应用根组件
 *
 * 职责：
 * - 路由视图渲染
 * - 路由切换时记录访问历史
 * - 全局消息提示（Snackbar）
 */
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useNavbarStore, useSnackbarStore } from "@/stores";

const router = useRouter();
const navbarStore = useNavbarStore();
const snackbarStore = useSnackbarStore();

onMounted(() => {
  // 路由切换时记录访问历史
  router.afterEach((to) => {
    // 只记录有标题的页面
    const title = (to.meta?.title as string) || "";
    if (title) {
      navbarStore.recordAccess({
        path: to.path,
        title,
        icon: (to.meta?.icon as string) || "mdi-file",
        category: to.meta?.category as string,
      });
    }
  });
});
</script>

<template>
  <v-app>
    <!-- 路由加载进度条 -->
    <LoadingBar />
    <!-- 路由视图 - 渲染匹配的路由组件 -->
    <router-view />
    <!-- 全局消息提示队列 -->
    <v-snackbar-queue v-model="snackbarStore.queue" location="bottom" timer />
  </v-app>
</template>
