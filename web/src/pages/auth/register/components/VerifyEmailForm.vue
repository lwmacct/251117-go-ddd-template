<script setup lang="ts">
/**
 * VerifyEmailForm 子组件
 * 邮箱验证表单，支持从 URL 参数初始化（用于独立访问场景）
 */

import { ref, computed, watch, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useRegister } from "../composables";
import { PlatformAuthAPI } from "@/api";

const route = useRoute();
const router = useRouter();
const registerStore = useRegister();

// 定义事件
const emit = defineEmits<{
  verified: [];
  "go-back": [];
}>();

// 状态
const verificationCode = ref("");
const isLoading = ref(false);
const errorMessage = ref("");
const isVerified = ref(false);
const successMessage = ref("");

// 计算属性
const isCodeValid = computed(() => {
  return verificationCode.value.length === 6;
});

// 从 URL 参数获取邮箱和验证码（独立访问场景）
onMounted(() => {
  const emailParam = route.query.email as string;
  const codeParam = route.query.code as string;

  // 如果 URL 中有邮箱，设置到 store（用于显示）
  if (emailParam && !registerStore.email.value) {
    registerStore.email.value = emailParam;
  }

  // 如果 URL 中有验证码，自动填充
  if (codeParam && codeParam.length === 6) {
    verificationCode.value = codeParam;
    // 如果是6位数字，自动验证
    const codePattern = /^\d{6}$/;
    if (codePattern.test(codeParam)) {
      handleVerify();
    }
  }
});

// 监听验证码输入，自动提交（如果长度足够）
watch(verificationCode, (newVal) => {
  if (newVal.length === 6 && !isVerified.value) {
    // 自动验证（如果是6位数字）
    const codePattern = /^\d{6}$/;
    if (codePattern.test(newVal)) {
      handleVerify();
    }
  }
});

/**
 * 验证邮箱
 */
async function handleVerify() {
  if (!isCodeValid.value) {
    errorMessage.value = "请输入6位验证码";
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";
  successMessage.value = "";

  try {
    // 优先使用 session token（注册流程），如果没有则使用邮箱（独立访问）
    if (registerStore.sessionToken.value) {
      // 注册流程：使用 session token
      await PlatformAuthAPI.verifyEmail({
        session_token: registerStore.sessionToken.value,
        code: verificationCode.value,
      });
    } else if (registerStore.email.value) {
      // 独立访问：使用邮箱
      await PlatformAuthAPI.verifyEmail({
        email: registerStore.email.value,
        code: verificationCode.value,
      });
    } else {
      throw new Error("缺少验证信息，请重新注册");
    }

    // 验证成功
    isVerified.value = true;
    successMessage.value = "邮箱验证成功！您现在可以登录使用";

    // 3秒后触发 verified 事件（如果是注册流程）或直接跳转（独立访问）
    setTimeout(() => {
      if (registerStore.sessionToken.value) {
        // 注册流程：触发事件让父组件处理
        emit("verified");
      } else {
        // 独立访问：直接跳转到登录页
        router.push("/auth/login");
      }
    }, 3000);
  } catch (error: any) {
    errorMessage.value = error.message || "验证失败，请检查验证码是否正确";
    verificationCode.value = "";
  } finally {
    isLoading.value = false;
  }
}

/**
 * 返回注册表单
 */
function goBack() {
  verificationCode.value = "";
  errorMessage.value = "";
  emit("go-back");
}
</script>

<template>
  <div class="verify-email-form-wrapper">
    <v-card class="verify-email-card" elevation="12">
      <!-- 返回按钮 -->
      <div class="pa-4 pb-0">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" size="small" @click="goBack"> 返回注册 </v-btn>
      </div>

      <v-card-title class="text-h4 text-center pt-4 pb-6">
        <v-icon size="large" color="primary" class="mr-3">mdi-email-check</v-icon>
        验证邮箱
      </v-card-title>

      <v-card-text class="px-8 pb-8">
        <!-- 成功状态 -->
        <template v-if="isVerified">
          <v-alert type="success" variant="tonal" class="mb-4">
            {{ successMessage }}
          </v-alert>
          <div class="text-body-1 text-medium-emphasis mb-4 text-center">
            {{ registerStore.sessionToken.value ? "正在跳转到登录页面..." : "您可以立即登录使用" }}
          </div>
          <v-btn v-if="!registerStore.sessionToken.value" color="primary" variant="elevated" block @click="router.push('/auth/login')"> 立即登录 </v-btn>
        </template>

        <!-- 验证表单 -->
        <template v-else>
          <v-alert type="info" variant="tonal" density="compact" class="mb-6">
            <div class="text-body-2">
              <div class="font-weight-bold mb-2">我们已向您的邮箱发送了验证码：</div>
              <div class="text-medium-emphasis">
                {{ registerStore.email.value || "您的邮箱" }}
              </div>
              <div class="mt-2">请查收邮件并输入6位验证码完成注册验证。</div>
            </div>
          </v-alert>

          <v-form @submit.prevent="handleVerify">
            <v-text-field
              v-model="verificationCode"
              label="验证码"
              prepend-inner-icon="mdi-shield-check"
              variant="outlined"
              type="text"
              required
              :maxlength="6"
              :disabled="isLoading"
              placeholder="请输入6位验证码"
              autocomplete="one-time-code"
              autofocus
              class="mb-4"
              :error="errorMessage.length > 0"
              :error-messages="errorMessage"
              hint="请输入邮件中收到的6位数字验证码"
              persistent-hint
            ></v-text-field>

            <v-btn color="primary" variant="elevated" prepend-icon="mdi-check" block size="large" type="submit" :loading="isLoading" :disabled="!isCodeValid || isLoading" class="mt-2">
              {{ isLoading ? "验证中..." : "验证邮箱" }}
            </v-btn>
          </v-form>

          <v-divider class="mt-4"></v-divider>
          <div class="text-center mt-4 text-caption text-medium-emphasis">
            <div>没有收到邮件？请检查垃圾邮件文件夹，或稍后重试</div>
          </div>
        </template>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.verify-email-form-wrapper {
  width: 100%;
  margin: 0 auto;
}

.verify-email-card {
  border-radius: 16px;
  backdrop-filter: blur(10px);
}
</style>
