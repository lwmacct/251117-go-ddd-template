/**
 * Register 页面状态管理 Composable
 * 
 * 管理注册表单状态，包括：
 * - 注册表单数据（邮箱、密码、确认密码、验证码）
 * - 验证码获取和显示
 * - 错误和成功消息管理
 * - Session token（用于邮箱验证流程）
 */

import { ref, computed } from 'vue'
import { PlatformAuthAPI } from '@/api'
import type { CaptchaData, PlatformRegisterRequest } from '@/api'

/**
 * Register 页面状态管理 Composable
 * 
 * 管理注册表单状态（支持邮箱+密码+验证码注册）
 */
export function useRegister() {
  // === 状态 ===
  const email = ref('') // 邮箱
  const password = ref('') // 密码
  const confirmPassword = ref('') // 确认密码
  const errorMessage = ref('') // 错误消息
  const successMessage = ref('') // 成功消息
  const isLoading = ref(false) // 加载状态
  const captchaData = ref<CaptchaData | null>(null) // 验证码数据
  const captchaCode = ref('') // 验证码输入
  const loadingCaptcha = ref(false) // 验证码加载状态
  const sessionToken = ref<string>('') // 临时会话token（用于邮箱验证流程）

  // 错误提示自动消失定时器
  let errorTimer: ReturnType<typeof setTimeout> | null = null
  let successTimer: ReturnType<typeof setTimeout> | null = null

  // === 计算属性 ===
  const isFormValid = computed(() => {
    return (
      email.value.includes('@') &&
      password.value.length >= 6 &&
      confirmPassword.value === password.value &&
      captchaCode.value.length > 0 &&
      captchaData.value !== null
    )
  })

  const passwordMatch = computed(() => {
    if (!confirmPassword.value) return true
    return password.value === confirmPassword.value
  })

  // 验证码图片（用于模板，避免类型错误）
  const captchaImage = computed(() => captchaData.value?.image || '')

  // === 方法 ===

  /**
   * 显示错误消息（5秒后自动消失）
   */
  const showErrorMessage = (message: string, duration = 5000) => {
    if (errorTimer) {
      clearTimeout(errorTimer)
    }

    errorMessage.value = message
    successMessage.value = ''

    if (duration > 0) {
      errorTimer = setTimeout(() => {
        errorMessage.value = ''
      }, duration)
    }
  }

  /**
   * 显示成功消息
   */
  const showSuccessMessage = (message: string, duration = 3000) => {
    if (successTimer) {
      clearTimeout(successTimer)
    }

    successMessage.value = message
    errorMessage.value = ''

    if (duration > 0) {
      successTimer = setTimeout(() => {
        successMessage.value = ''
      }, duration)
    }
  }

  /**
   * 清除错误消息
   */
  const clearError = () => {
    if (errorTimer) {
      clearTimeout(errorTimer)
      errorTimer = null
    }
    errorMessage.value = ''
  }

  /**
   * 获取验证码
   */
  const fetchCaptcha = async () => {
    try {
      loadingCaptcha.value = true
      const response = await PlatformAuthAPI.getCaptcha()
      if (response.data) {
        captchaData.value = response.data
        captchaCode.value = ''
      }
    } catch (error) {
      showErrorMessage(error instanceof Error ? error.message : '获取验证码失败')
    } finally {
      loadingCaptcha.value = false
    }
  }

  /**
   * 注册
   */
  const register = async () => {
    if (!isFormValid.value) {
      showErrorMessage('请填写完整的注册信息')
      return { success: false, message: '请填写完整的注册信息' }
    }

    if (!passwordMatch.value) {
      showErrorMessage('两次输入的密码不一致')
      return { success: false, message: '两次输入的密码不一致' }
    }

    if (!captchaData.value) {
      showErrorMessage('请先获取验证码')
      return { success: false, message: '请先获取验证码' }
    }

    clearError()
    isLoading.value = true

    try {
      const requestData: PlatformRegisterRequest = {
        email: email.value,
        password: password.value,
        captcha_id: captchaData.value.id,
        captcha: captchaCode.value,
      }

      const response = await PlatformAuthAPI.register(requestData)

      if (response.code === 200) {
        // 保存session token（用于邮箱验证）
        if (response.data?.session_token) {
          sessionToken.value = response.data.session_token
        }
        showSuccessMessage('注册成功！请查收邮箱验证码完成注册')
        // 不清空表单，保留邮箱以便显示在验证页面
        return {
          success: true,
          message: response.message,
        }
      }

      showErrorMessage(response.message || '注册失败')
      await fetchCaptcha()
      return {
        success: false,
        message: response.message || '注册失败',
      }
    } catch (error: any) {
      const message = error.message || '注册失败'
      showErrorMessage(message)
      await fetchCaptcha()
      return {
        success: false,
        message,
      }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 清空表单
   */
  const clearForm = () => {
    email.value = ''
    password.value = ''
    confirmPassword.value = ''
    captchaCode.value = ''
    clearError()
    captchaData.value = null
  }

  return {
    // 状态
    email,
    password,
    confirmPassword,
    errorMessage,
    successMessage,
    isLoading,
    captchaData,
    captchaCode,
    loadingCaptcha,
    sessionToken,

    // 计算属性
    isFormValid,
    passwordMatch,
    captchaImage,

    // 方法
    register,
    clearForm,
    fetchCaptcha,
  }
}

