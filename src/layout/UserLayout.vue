<!--
  UserLayout.vue - 用户中心布局
  基于 BaseLayout，使用动态菜单（基于角色过滤）
-->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { useMenus } from "@/composables";

const { userMenus } = useMenus();
</script>

<template>
  <BaseLayout :menu-items="userMenus" />
</template>
