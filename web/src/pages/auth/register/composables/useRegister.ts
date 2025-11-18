/**
 * 注册页面专用 Composable
 */
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/auth'
import { formatAuthError } from '@/utils/auth'
import type { RegisterRequest } from '@/types/auth'

/** 注册成功后的默认跳转路径 */
const DEFAULT_REDIRECT_PATH = '/admin/overview'

/**
 * 注册逻辑
 */
export function useRegister() {
  const router = useRouter()

  // 表单数据
  const formData = ref<RegisterRequest>({
    username: '',
    email: '',
    password: '',
    full_name: '',
  })

  // 确认密码
  const confirmPassword = ref('')

  // 表单状态
  const loading = ref(false)
  const errorMessage = ref('')

  /**
   * 执行注册
   */
  const handleRegister = async (): Promise<boolean> => {
    loading.value = true
    errorMessage.value = ''

    try {
      await register(formData.value)
      // 注册成功，跳转到管理后台
      router.push(DEFAULT_REDIRECT_PATH)
      return true
    } catch (error: any) {
      errorMessage.value = formatAuthError(error)
      console.error('Register error:', error)
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
      username: '',
      email: '',
      password: '',
      full_name: '',
    }
    confirmPassword.value = ''
    errorMessage.value = ''
  }

  return {
    formData,
    confirmPassword,
    loading,
    errorMessage,
    handleRegister,
    clearError,
    resetForm,
  }
}
