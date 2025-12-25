<!--
  UserLayout.vue - 用户中心布局
  基于 BaseLayout，配置用户中心菜单
-->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { userMenuItems } from "@/config/menus";
</script>

<template>
  <BaseLayout :menu-items="userMenuItems" />
</template>
