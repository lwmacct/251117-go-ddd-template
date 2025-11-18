<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useLogin } from "../composables";
import { useAuthStore } from "@/stores/auth";
import { PlatformAuthAPI } from "@/api";

const router = useRouter();
const loginStore = useLogin();
const authStore = useAuthStore();

// 定义事件
const emit = defineEmits<{
  verified: [];
}>();

// 状态
const twoFactorCode = ref("");
const isLoading = ref(false);
const errorMessage = ref("");
const showRecoveryCodeHelp = ref(false);

// 计算属性
const isCodeValid = computed(() => {
  // TOTP验证码是6位数字，恢复码通常是更长的字符串
  return twoFactorCode.value.length >= 6;
});

// 监听验证码输入，自动提交 (如果长度足够)
watch(twoFactorCode, (newVal) => {
  if (newVal.length === 6) {
    // 自动验证TOTP码 (如果是6位数字)
    const totpPattern = /^\d{6}$/;
    if (totpPattern.test(newVal)) {
      handleVerify();
    }
  }
});

/**
 * 验证2FA验证码
 */
async function handleVerify() {
  if (!isCodeValid.value) {
    errorMessage.value = "请输入6位验证码或恢复码";
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    // 调用2FA验证API
    const response = await PlatformAuthAPI.verify2FA({
      session_token: loginStore.sessionToken.value,
      code: twoFactorCode.value,
    });

    if (response.code === 200) {
      // 2FA验证成功，发出事件
      emit("verified");
    } else {
      errorMessage.value = response.message || "验证码错误，请重试";
      twoFactorCode.value = "";
    }
  } catch (error: any) {
    errorMessage.value = error.message || "验证失败，请检查验证码";
    twoFactorCode.value = "";
  } finally {
    isLoading.value = false;
  }
}

/**
 * 返回登录页面
 */
function goBack() {
  router.push("/auth/login");
}
</script>

<template>
  <div class="two-factor-form-wrapper">
    <v-card class="two-factor-card" elevation="12">
      <!-- 返回按钮 -->
      <div class="pa-4 pb-0">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" size="small" @click="goBack"> 返回登录 </v-btn>
      </div>

      <v-card-title class="text-h4 text-center pt-4 pb-6">
        <v-icon size="large" color="primary" class="mr-3">mdi-shield-lock</v-icon>
        双因素认证验证
      </v-card-title>

      <v-card-text class="px-8 pb-8">
        <v-alert type="info" variant="tonal" density="compact" class="mb-6">
          <div class="text-body-2">
            <div class="font-weight-bold mb-2">请输入您的双因素认证验证码：</div>
            <div>• 6位数字验证码 (TOTP)</div>
            <div>• 或使用恢复码 (Recovery Code)</div>
          </div>
        </v-alert>

        <v-form @submit.prevent="handleVerify">
          <v-text-field
            v-model="twoFactorCode"
            label="验证码"
            prepend-inner-icon="mdi-shield-check"
            variant="outlined"
            type="text"
            required
            :maxlength="20"
            :disabled="isLoading"
            placeholder="请输入6位验证码或恢复码"
            autocomplete="one-time-code"
            autofocus
            class="mb-4"
            :error="errorMessage.length > 0"
            :error-messages="errorMessage"
            hint="输入您手机应用中的6位验证码，或使用恢复码"
            persistent-hint
          >
            <template #append-inner>
              <v-btn icon="mdi-help-circle-outline" variant="text" size="small" @click="showRecoveryCodeHelp = !showRecoveryCodeHelp"></v-btn>
            </template>
          </v-text-field>

          <!-- 恢复码帮助 -->
          <v-expand-transition>
            <v-alert v-if="showRecoveryCodeHelp" type="warning" variant="tonal" density="compact" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">恢复码说明：</div>
                <div>如果您无法访问认证器应用，可以使用恢复码。恢复码是在设置2FA时生成的备用验证码，请妥善保管。</div>
              </div>
            </v-alert>
          </v-expand-transition>

          <v-btn color="primary" variant="elevated" prepend-icon="mdi-check" block size="large" type="submit" :loading="isLoading" :disabled="!isCodeValid || isLoading" class="mt-2">
            {{ isLoading ? "验证中..." : "验证" }}
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.two-factor-form-wrapper {
  width: 100%;
  margin: 0 auto;
}

.two-factor-card {
  border-radius: 16px;
  backdrop-filter: blur(10px);
}
</style>
