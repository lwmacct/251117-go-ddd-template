/**
 * 登录页面专用 Composable
 */
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/api/auth'
import { formatAuthError } from '@/utils/auth'
import type { LoginRequest } from '@/types/auth'

/** 登录成功后的默认跳转路径 */
const DEFAULT_REDIRECT_PATH = '/admin/overview'

/**
 * 登录逻辑
 */
export function useLogin() {
  const router = useRouter()

  // 表单数据
  const formData = ref<LoginRequest>({
    login: '',
    password: '',
  })

  // 表单状态
  const loading = ref(false)
  const errorMessage = ref('')

  /**
   * 执行登录
   */
  const handleLogin = async (): Promise<boolean> => {
    loading.value = true
    errorMessage.value = ''

    try {
      await login(formData.value)
      // 登录成功，跳转到管理后台
      router.push(DEFAULT_REDIRECT_PATH)
      return true
    } catch (error: any) {
      errorMessage.value = formatAuthError(error)
      console.error('Login error:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  /**
   * 清除错误信息
   */
  const clearError = () => {
    errorMessage.value = ''
  }

  /**
   * 重置表单
   */
  const resetForm = () => {
    formData.value = {
      login: '',
      password: '',
    }
    errorMessage.value = ''
  }

  return {
    formData,
    loading,
    errorMessage,
    handleLogin,
    clearError,
    resetForm,
  }
}
