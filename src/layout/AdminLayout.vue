<!--
  AdminLayout.vue - 管理后台布局
  基于 BaseLayout，使用动态菜单（基于角色过滤）
-->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { useMenus } from "@/composables";

const { adminMenus } = useMenus();
</script>

<template>
  <BaseLayout :menu-items="adminMenus" />
</template>
