<script setup lang="ts">
/**
 * Login 页面
 *
 * 仅包含表单切换逻辑，布局由父组件 auth/index.vue 提供
 */

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
  // 获取重定向目标 (从 query 参数)
  const redirectTo = (router.currentRoute.value.query.redirect as string) || null;

  if (redirectTo && redirectTo !== "/auth/login" && redirectTo !== "/auth/register") {
    // 有重定向目标，跳转回去
    await router.replace(redirectTo);
  } else {
    // 没有重定向目标，跳转到管理后台首页
    await router.replace("/admin/overview");
  }
};
</script>

<template>
  <div class="login-wrapper">
    <LoginForm v-if="!showTwoFactor" @login-success="handleLoginSuccess" @requires-two-factor="showTwoFactor = true" />
    <TwoFactorForm v-else @verified="handleLoginSuccess" @back="showTwoFactor = false" />
  </div>
</template>

<style scoped>
.login-wrapper {
  width: 100%;
}
</style>
