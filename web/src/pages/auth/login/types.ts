/**
 * Login 页面相关的类型定义
 */

/**
 * 登录表单接口
 */
export interface LoginForm {
  email: string
  password: string
  rememberMe: boolean
}

/**
 * 登录响应接口
 */
export interface LoginResponse {
  success: boolean
  token?: string
  message?: string
  user?: {
    id: string
    email: string
    name: string
  }
}

/**
 * Login 页面数据接口
 */
export interface LoginPageData {
  pageTitle: string
  pageIcon: string
  backgroundGradient: string
}
