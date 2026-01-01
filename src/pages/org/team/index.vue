<!--
  组织团队工作台页面
-->
<script setup lang="ts">
interface Props {
  orgId: number;
  teamId: number;
}

defineProps<Props>();
</script>

<template>
  <v-container class="fill-height">
    <v-card class="mx-auto text-center pa-8" max-width="600">
      <v-icon size="64" color="primary" class="mb-4">mdi-domain</v-icon>
      <v-card-title class="text-h5">组织工作台</v-card-title>
      <v-card-text>
        <p>组织 ID: {{ orgId }}</p>
        <p>团队 ID: {{ teamId }}</p>
        <p class="text-medium-emphasis mt-4">功能开发中...</p>
      </v-card-text>
    </v-card>
  </v-container>
</template>
