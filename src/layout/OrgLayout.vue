<!--
  OrgLayout.vue - 组织相关布局
  用于组织团队功能（任务、成员管理等）
-->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { useMenus } from "@/composables";

/**
 * 组织布局菜单
 * 目前只有团队任务，未来可扩展：团队设置、成员管理等
 */
const { orgMenus } = useMenus();
</script>

<template>
  <BaseLayout :menu-items="orgMenus" />
</template>
