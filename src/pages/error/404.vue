<!--
  404.vue - 页面不存在
  当访问不存在的路由时显示此页面
-->
<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

function goBack() {
  router.back();
}

function goHome() {
  router.push("/");
}
</script>

<template>
  <v-container class="fill-height">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" class="text-center">
        <v-icon size="120" color="warning" class="mb-4">mdi-file-question</v-icon>

        <h1 class="text-h2 font-weight-bold mb-2">404</h1>

        <p class="text-h5 text-medium-emphasis mb-6">页面不存在</p>

        <p class="text-body-1 text-medium-emphasis mb-8">抱歉，您访问的页面不存在或已被移除。</p>

        <div class="d-flex justify-center gap-4">
          <v-btn variant="outlined" prepend-icon="mdi-arrow-left" @click="goBack"> 返回上一页 </v-btn>

          <v-btn color="primary" prepend-icon="mdi-home" @click="goHome"> 返回首页 </v-btn>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
