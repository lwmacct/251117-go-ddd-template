/**
 * 认证状态管理 Store (Pinia)
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  login as apiLogin,
  register as apiRegister,
  logout as apiLogout,
  getCurrentUser,
} from '@/api/auth'
import { getAccessToken, clearAuthTokens } from '@/utils/auth'
import type { LoginRequest, RegisterRequest, User } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const currentUser = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!currentUser.value)
  const hasToken = computed(() => !!getAccessToken())

  /**
   * 初始化认证状态
   * 检查 localStorage 中的 token，如果存在则获取用户信息
   */
  async function initAuth() {
    const token = getAccessToken()
    if (!token) {
      currentUser.value = null
      return
    }

    try {
      isLoading.value = true
      error.value = null
      const user = await getCurrentUser()
      currentUser.value = user
    } catch (err: any) {
      // token 无效，清除并重置状态
      clearAuthTokens()
      currentUser.value = null
      error.value = err.message || 'Failed to initialize auth'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 登录
   */
  async function login(credentials: LoginRequest) {
    try {
      isLoading.value = true
      error.value = null
      const response = await apiLogin(credentials)
      currentUser.value = response.user
      return response
    } catch (err: any) {
      error.value = err.response?.data?.error || err.message || 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 注册
   */
  async function register(data: RegisterRequest) {
    try {
      isLoading.value = true
      error.value = null
      const response = await apiRegister(data)
      currentUser.value = response.user
      return response
    } catch (err: any) {
      error.value =
        err.response?.data?.error || err.message || 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 登出
   */
  function logout() {
    apiLogout()
    currentUser.value = null
    error.value = null
  }

  /**
   * 清除错误
   */
  function clearError() {
    error.value = null
  }

  /**
   * 更新用户信息
   */
  function updateUser(user: User) {
    currentUser.value = user
  }

  return {
    // 状态
    currentUser,
    isLoading,
    error,

    // 计算属性
    isAuthenticated,
    hasToken,

    // 方法
    initAuth,
    login,
    register,
    logout,
    clearError,
    updateUser,
  }
})
