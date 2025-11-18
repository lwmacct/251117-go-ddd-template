/**
 * 注册页面专用类型定义
 */

/**
 * 邮箱验证请求
 * (预留，用于未来实现邮箱验证)
 */
export interface VerifyEmailRequest {
  code: string
  email: string
}

/**
 * 邮箱验证响应
 */
export interface VerifyEmailResponse {
  success: boolean
  message?: string
}
