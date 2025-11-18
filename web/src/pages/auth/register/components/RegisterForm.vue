<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useRegister } from "../composables";

const router = useRouter();
const registerStore = useRegister();

// 定义事件
const emit = defineEmits<{
  registerSuccess: [];
  requiresVerification: [];
}>();

// 处理注册
const handleRegister = async () => {
  const result = await registerStore.register();

  if (result.success) {
    // 注册成功，需要邮箱验证
    emit("requiresVerification");
  }
};

// 组件挂载时获取验证码
onMounted(() => {
  registerStore.fetchCaptcha();
});
</script>

<template>
  <div class="register-form-wrapper">
    <v-card class="register-card" elevation="12">
      <!-- 返回首页按钮 -->
      <div class="pa-4 pb-0">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" size="small" @click="router.push('/')"> 返回首页 </v-btn>
      </div>

      <v-card-title class="text-h4 text-center pt-4 pb-6">
        <v-icon size="large" color="primary" class="mr-3">mdi-account-plus</v-icon>
        用户注册
      </v-card-title>

      <v-card-text class="px-8 pb-8">
        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="registerStore.email.value"
            label="邮箱"
            prepend-inner-icon="mdi-email"
            variant="outlined"
            type="email"
            name="email"
            autocomplete="email"
            required
            class="mb-4"
            :disabled="registerStore.isLoading.value"
            :error="registerStore.email.value.length > 0 && !registerStore.email.value.includes('@')"
            :error-messages="registerStore.email.value.length > 0 && !registerStore.email.value.includes('@') ? '请输入有效的邮箱地址' : ''"
            hint="用户名将自动从邮箱生成 (可修改)"
            persistent-hint
          ></v-text-field>

          <!-- 密码和确认密码在同一行 -->
          <v-row dense class="mb-3">
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="registerStore.password.value"
                label="密码"
                prepend-inner-icon="mdi-lock"
                variant="outlined"
                type="password"
                name="password"
                autocomplete="new-password"
                required
                :disabled="registerStore.isLoading.value"
                :error="registerStore.password.value.length > 0 && registerStore.password.value.length < 6"
                :error-messages="registerStore.password.value.length > 0 && registerStore.password.value.length < 6 ? '密码至少6个字符' : ''"
              ></v-text-field>
            </v-col>

            <v-col cols="12" sm="6">
              <v-text-field
                v-model="registerStore.confirmPassword.value"
                label="确认密码"
                prepend-inner-icon="mdi-lock-check"
                variant="outlined"
                type="password"
                name="confirm-password"
                autocomplete="new-password"
                required
                :disabled="registerStore.isLoading.value"
                :error="registerStore.confirmPassword.value.length > 0 && !registerStore.passwordMatch.value"
                :error-messages="registerStore.confirmPassword.value.length > 0 && !registerStore.passwordMatch.value ? '两次输入的密码不一致' : ''"
              ></v-text-field>
            </v-col>
          </v-row>

          <!-- 验证码 -->
          <div class="mb-4">
            <div class="d-flex align-start gap-3">
              <v-text-field
                v-model="registerStore.captchaCode.value"
                label="验证码"
                prepend-inner-icon="mdi-shield-check"
                variant="outlined"
                name="captcha"
                autocomplete="off"
                required
                :disabled="registerStore.isLoading.value"
                placeholder="请输入验证码(不区分大小写)"
                maxlength="4"
                style="flex: 1"
              ></v-text-field>

              <!-- 验证码图片 -->
              <div v-if="registerStore.captchaData" class="captcha-image" @click="registerStore.fetchCaptcha" title="点击刷新验证码">
                <img :src="registerStore.captchaImage.value" alt="验证码" />
              </div>

              <!-- 加载中显示 -->
              <div v-else class="captcha-loading">
                <v-progress-circular indeterminate color="primary" size="30"></v-progress-circular>
              </div>
            </div>
          </div>

          <v-btn color="primary" variant="elevated" prepend-icon="mdi-account-plus" block size="large" type="submit" :loading="registerStore.isLoading.value" :disabled="!registerStore.isFormValid.value" class="mt-2">
            {{ registerStore.isLoading.value ? "注册中..." : "注册" }}
          </v-btn>

          <!-- 消息提示区域 - 预留固定空间 -->
          <div class="message-area mt-3">
            <v-fade-transition>
              <v-alert v-if="registerStore.errorMessage.value" type="error" density="compact" variant="tonal" closable @click:close="registerStore.errorMessage.value = ''">
                {{ registerStore.errorMessage.value }}
              </v-alert>
            </v-fade-transition>

            <v-fade-transition>
              <v-alert v-if="registerStore.successMessage.value" type="success" density="compact" variant="tonal" closable @click:close="registerStore.successMessage.value = ''">
                {{ registerStore.successMessage.value }}
              </v-alert>
            </v-fade-transition>
          </div>
        </v-form>

        <!-- 用户协议提示 -->
        <v-divider class="mt-4"></v-divider>
        <div class="text-center mt-4 text-caption text-medium-emphasis">
          <div>注册即表示您同意服务条款和隐私政策</div>
        </div>
      </v-card-text>

      <v-card-actions class="justify-center pb-6 flex-column gap-2">
        <v-btn @click="router.push('/auth/login')" variant="text" prepend-icon="mdi-login"> 已有账号？去登录 </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<style scoped>
.register-form-wrapper {
  width: 100%;
  margin: 0 auto;
}

.register-card {
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

/* 消息提示区域 - 预留固定最小高度，避免布局跳动 */
.message-area {
  min-height: 48px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-start;
}
</style>
