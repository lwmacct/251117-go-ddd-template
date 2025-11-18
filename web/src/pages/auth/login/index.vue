<script setup lang="ts">
import LoginForm from './components/LoginForm.vue'
</script>

<template>
  <v-container fluid class="fill-height pa-0">
    <v-row no-gutters class="fill-height">
      <!-- 左侧背景装饰 -->
      <v-col cols="12" md="6" class="d-none d-md-flex">
        <v-sheet color="primary" class="d-flex align-center justify-center flex-column pa-12" height="100%"
          width="100%">
          <v-icon size="120" color="white" class="mb-8">
            mdi-account-circle
          </v-icon>
          <h1 class="text-h3 text-white font-weight-bold mb-4">
            欢迎回来
          </h1>
          <p class="text-h6 text-white text-center">
            登录您的账户以继续使用我们的服务
          </p>
        </v-sheet>
      </v-col>

      <!-- 右侧登录表单 -->
      <v-col cols="12" md="6">
        <v-sheet class="d-flex align-center justify-center" height="100%" width="100%">
          <LoginForm />
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.fill-height {
  height: 100vh;
}
</style>
