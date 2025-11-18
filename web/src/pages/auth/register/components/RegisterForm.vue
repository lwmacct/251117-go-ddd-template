<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  usernameRules,
  emailRules,
  passwordRules,
  fullNameRules,
  createConfirmPasswordRules,
} from '@/composables/auth'
import { useRegister } from '../composables'

/** 认证路由路径 */
const AUTH_ROUTES = {
  LOGIN: '/auth/login',
} as const

const router = useRouter()
const {
  formData,
  confirmPassword,
  loading,
  errorMessage,
  handleRegister,
  clearError,
} = useRegister()

// 表单状态
const showPassword = ref(false)
const showConfirmPassword = ref(false)

// 表单引用
const form = ref()

// 确认密码验证规则（动态生成）
const confirmPasswordRules = computed(() =>
  createConfirmPasswordRules(formData.value.password)
)

/**
 * 处理表单提交
 */
const onSubmit = async () => {
  const { valid } = await form.value.validate()
  if (!valid) return
  await handleRegister()
}

/**
 * 跳转到登录页
 */
const goToLogin = () => {
  router.push(AUTH_ROUTES.LOGIN)
}
</script>

<template>
  <v-card flat max-width="450" width="100%" class="pa-4 pa-sm-8">
    <!-- 标题 -->
    <div class="text-center mb-6">
      <h2 class="text-h4 font-weight-bold mb-2">注册</h2>
      <p class="text-body-1 text-medium-emphasis">
        创建一个新账户
      </p>
    </div>

    <!-- 错误提示 -->
    <v-alert v-if="errorMessage" type="error" variant="tonal" closable class="mb-4"
      @click:close="clearError">
      {{ errorMessage }}
    </v-alert>

    <!-- 注册表单 -->
    <v-form ref="form" @submit.prevent="onSubmit">
      <v-text-field v-model="formData.username" :rules="usernameRules" label="用户名"
        prepend-inner-icon="mdi-account" variant="outlined" class="mb-3" :disabled="loading"
        hint="只能包含字母、数字和下划线" />

      <v-text-field v-model="formData.email" :rules="emailRules" label="邮箱" type="email"
        prepend-inner-icon="mdi-email" variant="outlined" class="mb-3" :disabled="loading" />

      <v-text-field v-model="formData.full_name" :rules="fullNameRules" label="姓名（可选）"
        prepend-inner-icon="mdi-card-account-details" variant="outlined" class="mb-3" :disabled="loading" />

      <v-text-field v-model="formData.password" :rules="passwordRules"
        :type="showPassword ? 'text' : 'password'" label="密码" prepend-inner-icon="mdi-lock"
        :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'" variant="outlined" class="mb-3"
        :disabled="loading" @click:append-inner="showPassword = !showPassword" />

      <v-text-field v-model="confirmPassword" :rules="confirmPasswordRules"
        :type="showConfirmPassword ? 'text' : 'password'" label="确认密码" prepend-inner-icon="mdi-lock-check"
        :append-inner-icon="showConfirmPassword ? 'mdi-eye' : 'mdi-eye-off'" variant="outlined" class="mb-6"
        :disabled="loading" @click:append-inner="showConfirmPassword = !showConfirmPassword" />

      <!-- 注册按钮 -->
      <v-btn type="submit" color="primary" size="large" block :loading="loading" class="mb-4">
        注册
      </v-btn>

      <!-- 分割线 -->
      <v-divider class="my-6" />

      <!-- 登录提示 -->
      <div class="text-center">
        <span class="text-body-2 text-medium-emphasis">
          已有账户?
        </span>
        <v-btn variant="text" color="primary" class="text-none" @click="goToLogin">
          立即登录
        </v-btn>
      </div>
    </v-form>
  </v-card>
</template>
