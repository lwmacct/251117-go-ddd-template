/**
 * 认证 Composable
 * 封装 Pinia store 的使用，提供更简洁的 API
 *
 * 注意：根据 Pinia 最佳实践，不应该直接解构 store
 * 应该返回 store 实例或使用 storeToRefs() 保持响应式
 */
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'

/**
 * 认证 Composable
 * 提供认证相关的状态和方法
 */
export function useAuth() {
  const router = useRouter()
  const authStore = useAuthStore()

  // 使用 storeToRefs 保持响应式（只用于 state 和 getters）
  const { currentUser, isLoading, error, isAuthenticated, hasToken } = storeToRefs(authStore)

  // actions 可以直接解构
  const { initAuth, login, register, logout, clearError, updateUser } = authStore

  /**
   * 登出并跳转到登录页
   */
  const logoutAndRedirect = () => {
    logout()
    router.push('/auth/login')
  }

  /**
   * 检查认证状态并重定向
   * 如果已登录则跳转到指定页面
   */
  const checkAuthAndRedirect = (redirectTo: string = '/admin/overview') => {
    if (isAuthenticated.value) {
      router.push(redirectTo)
      return true
    }
    return false
  }

  return {
    // 响应式状态（通过 storeToRefs）
    currentUser,
    isLoading,
    error,
    isAuthenticated,
    hasToken,

    // actions
    initAuth,
    login,
    register,
    logout,
    clearError,
    updateUser,

    // 额外的便捷方法
    logoutAndRedirect,
    checkAuthAndRedirect,
  }
}
