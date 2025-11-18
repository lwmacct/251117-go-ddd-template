<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import LoginForm from "./components/LoginForm.vue";
import TwoFactorForm from "./components/TwoFactorForm.vue";

const router = useRouter();

// 是否显示2FA验证页面
const showTwoFactor = ref(false);

/**
 * 处理登录成功后的路由跳转
 */
const handleLoginSuccess = async () => {
  console.log("🎉 登录成功，正在跳转...");

  // 获取重定向目标（从 query 参数）
  const redirectTo = (router.currentRoute.value.query.redirect as string) || null;

  if (redirectTo && redirectTo !== "/auth/login" && redirectTo !== "/auth/register") {
    // 有重定向目标，跳转回去
    console.log("📍 返回来源页面:", redirectTo);
    await router.replace(redirectTo);
  } else {
    // 没有重定向目标，跳转到管理后台首页
    console.log("📍 跳转到管理后台");
    await router.replace("/admin/overview");
  }
};
</script>

<template>
  <!-- 登录页面不显示头部 -->

  <!-- 主要内容区域 -->
  <v-main class="d-flex align-center justify-center" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh">
    <v-container>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="6" lg="4">
          <LoginForm v-if="!showTwoFactor" @login-success="handleLoginSuccess" @requires-two-factor="showTwoFactor = true" />
          <TwoFactorForm v-else @verified="handleLoginSuccess" />
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<style scoped>
/* 登录页面样式 */
</style>
