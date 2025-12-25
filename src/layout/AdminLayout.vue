<!--
  AdminLayout.vue - 管理后台布局
  基于 BaseLayout，配置管理后台菜单
-->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { adminMenuItems } from "@/config/menus";
</script>

<template>
  <BaseLayout :menu-items="adminMenuItems" />
</template>
