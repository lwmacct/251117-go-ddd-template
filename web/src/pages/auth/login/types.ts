/**
 * 登录页面专用类型定义
 */

/**
 * 双因素认证请求
 * (预留，用于未来实现 2FA)
 */
export interface TwoFactorRequest {
  code: string
  user_id: number
}

/**
 * 双因素认证响应
 */
export interface TwoFactorResponse {
  success: boolean
  message?: string
}
