<script setup lang="ts">
/**
 * Register 页面
 *
 * 仅包含表单切换逻辑，布局由父组件 auth/index.vue 提供
 */

import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import RegisterForm from "./components/RegisterForm.vue";
import VerifyEmailForm from "./components/VerifyEmailForm.vue";

const route = useRoute();
const router = useRouter();

// 是否显示邮箱验证页面 (如果 URL 中有验证码，直接显示验证页面)
const showVerification = ref(false);

// 检查 URL 参数，如果是独立访问验证页面，直接显示验证表单
onMounted(() => {
  const emailParam = route.query.email as string;
  const codeParam = route.query.code as string;

  // 如果 URL 中有邮箱或验证码，直接显示验证页面 (独立访问场景)
  if (emailParam || codeParam) {
    showVerification.value = true;
  }
});

/**
 * 处理需要邮箱验证
 */
const handleRequiresVerification = () => {
  showVerification.value = true;
};

/**
 * 处理邮箱验证成功后的跳转
 */
const handleVerified = () => {
  // 验证成功，跳转到登录页
  router.push("/auth/login");
};

/**
 * 返回注册表单
 */
const handleGoBack = () => {
  showVerification.value = false;
};
</script>

<template>
  <div class="register-wrapper">
    <RegisterForm v-if="!showVerification" @requires-verification="handleRequiresVerification" />
    <VerifyEmailForm v-else @verified="handleVerified" @go-back="handleGoBack" />
  </div>
</template>

<style scoped>
.register-wrapper {
  width: 100%;
}
</style>
