<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useLogin } from "../composables";

const router = useRouter();
const loginStore = useLogin();

// 定义事件
const emit = defineEmits<{
  loginSuccess: [];
  requiresTwoFactor: [];
}>();

// 处理登录
const handleLogin = async () => {
  const result = await loginStore.login();

  // 检查是否需要2FA验证（优先检查，无论success状态）
  if (result.requiresTwoFactor) {
    emit("requiresTwoFactor");
    return;
  }

  // 登录成功
  if (result.success) {
    emit("loginSuccess");
  }
};

// 组件挂载时获取验证码
onMounted(() => {
  loginStore.fetchCaptcha();
});
</script>

<template>
  <div class="login-form-wrapper">
    <v-card class="login-card" elevation="12">
      <!-- 返回首页按钮 -->
      <div class="pa-4 pb-0">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" size="small" @click="router.push('/')"> 返回首页 </v-btn>
      </div>

      <v-card-title class="text-h4 text-center pt-4 pb-6">
        <v-icon size="large" color="primary" class="mr-3">mdi-login</v-icon>
        用户登录
      </v-card-title>

      <v-card-text class="px-8 pb-8">
        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="loginStore.account.value"
            label="邮箱 / 手机号 / 用户名"
            prepend-inner-icon="mdi-email"
            variant="outlined"
            required
            class="mb-4"
            :disabled="loginStore.isLoading.value"
            hint="支持手机号、用户名或邮箱登录"
            persistent-hint
          ></v-text-field>

          <v-text-field v-model="loginStore.password.value" label="密码" prepend-inner-icon="mdi-lock" variant="outlined" type="password" required class="mb-4" :disabled="loginStore.isLoading.value"></v-text-field>

          <!-- 验证码 -->
          <div class="mb-4">
            <div class="d-flex align-start gap-3">
              <v-text-field
                v-model="loginStore.captchaCode.value"
                label="验证码"
                prepend-inner-icon="mdi-shield-check"
                variant="outlined"
                required
                :disabled="loginStore.isLoading.value"
                placeholder="请输入验证码(不区分大小写)"
                maxlength="4"
                style="flex: 1"
              ></v-text-field>

              <!-- 验证码图片 -->
              <div v-if="loginStore.captchaData" class="captcha-image" @click="loginStore.fetchCaptcha" title="点击刷新验证码">
                <img :src="loginStore.captchaImage.value" alt="验证码" />
              </div>

              <!-- 加载中显示 -->
              <div v-else class="captcha-loading">
                <v-progress-circular indeterminate color="primary" size="30"></v-progress-circular>
              </div>
            </div>
          </div>

          <v-btn color="primary" variant="elevated" prepend-icon="mdi-login" block size="large" type="submit" :loading="loginStore.isLoading.value" :disabled="!loginStore.isFormValid.value" class="mt-2">
            {{ loginStore.isLoading.value ? "登录中..." : "登录" }}
          </v-btn>

          <!-- 错误提示区域 - 预留固定空间 -->
          <div class="error-message-area mt-3">
            <v-fade-transition>
              <v-alert v-if="loginStore.errorMessage.value" type="error" density="compact" variant="tonal" closable @click:close="loginStore.errorMessage.value = ''">
                {{ loginStore.errorMessage.value }}
              </v-alert>
            </v-fade-transition>
          </div>
        </v-form>

        <!-- 测试账号提示 -->
        <v-divider class="mt-4"></v-divider>
        <div class="text-center mt-4 text-caption text-medium-emphasis">
          <div>测试账号: root / Root@123456 或 admin / admin123</div>
        </div>
      </v-card-text>

      <v-card-actions class="justify-center pb-6 flex-column gap-2">
        <v-btn @click="router.push('/auth/register')" variant="text" prepend-icon="mdi-account-plus"> 没有账号？去注册 </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<style scoped>
.login-form-wrapper {
  width: 100%;
  margin: 0 auto;
}

.login-card {
  border-radius: 16px;
  backdrop-filter: blur(10px);
}

.captcha-image {
  width: 140px;
  height: 56px;
  min-width: 140px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  transition: all 0.3s ease;
  cursor: pointer;
  flex-shrink: 0;
  padding: 0 4px;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.captcha-image:hover {
  border-color: #1976d2;
  box-shadow: 0 2px 4px rgba(25, 118, 210, 0.2);
}

.captcha-loading {
  width: 140px;
  height: 56px;
  min-width: 140px;
  border: 1px dashed #e0e0e0;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
  flex-shrink: 0;
}

.gap-3 {
  gap: 12px;
}

.gap-2 {
  gap: 8px;
}

/* 错误提示区域 - 预留固定最小高度，避免布局跳动 */
.error-message-area {
  min-height: 48px;
  display: flex;
  align-items: flex-start;
}
</style>
