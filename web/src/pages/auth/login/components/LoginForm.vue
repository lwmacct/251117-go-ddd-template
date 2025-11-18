<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { loginFieldRules, passwordRules } from '@/composables/auth'
import { useLogin } from '../composables'

/** 认证路由路径 */
const AUTH_ROUTES = {
  REGISTER: '/auth/register',
} as const

const router = useRouter()
const { formData, loading, errorMessage, handleLogin, clearError } = useLogin()

// 表单状态
const showPassword = ref(false)

// 表单引用
const form = ref()

/**
 * 处理表单提交
 */
const onSubmit = async () => {
  const { valid } = await form.value.validate()
  if (!valid) return
  await handleLogin()
}

/**
 * 跳转到注册页
 */
const goToRegister = () => {
  router.push(AUTH_ROUTES.REGISTER)
}
</script>

<template>
  <v-card flat max-width="450" width="100%" class="pa-4 pa-sm-8">
    <!-- 标题 -->
    <div class="text-center mb-8">
      <h2 class="text-h4 font-weight-bold mb-2">登录</h2>
      <p class="text-body-1 text-medium-emphasis">
        使用您的账户登录
      </p>
    </div>

    <!-- 错误提示 -->
    <v-alert v-if="errorMessage" type="error" variant="tonal" closable class="mb-4"
      @click:close="clearError">
      {{ errorMessage }}
    </v-alert>

    <!-- 登录表单 -->
    <v-form ref="form" @submit.prevent="onSubmit">
      <v-text-field v-model="formData.login" :rules="loginFieldRules" label="用户名或邮箱"
        prepend-inner-icon="mdi-account" variant="outlined" class="mb-4" :disabled="loading" />

      <v-text-field v-model="formData.password" :rules="passwordRules"
        :type="showPassword ? 'text' : 'password'" label="密码" prepend-inner-icon="mdi-lock"
        :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'" variant="outlined" class="mb-2"
        :disabled="loading" @click:append-inner="showPassword = !showPassword" />

      <!-- 忘记密码链接 -->
      <div class="text-right mb-6">
        <v-btn variant="text" size="small" color="primary" class="text-none">
          忘记密码?
        </v-btn>
      </div>

      <!-- 登录按钮 -->
      <v-btn type="submit" color="primary" size="large" block :loading="loading" class="mb-4">
        登录
      </v-btn>

      <!-- 分割线 -->
      <v-divider class="my-6" />

      <!-- 注册提示 -->
      <div class="text-center">
        <span class="text-body-2 text-medium-emphasis">
          还没有账户?
        </span>
        <v-btn variant="text" color="primary" class="text-none" @click="goToRegister">
          立即注册
        </v-btn>
      </div>
    </v-form>
  </v-card>
</template>
