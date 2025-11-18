/**
 * 认证相关类型定义
 */
import type { User } from './user'

/** 登录请求 */
export interface LoginRequest {
  login: string // 用户名或邮箱
  password: string
}

/** 注册请求 */
export interface RegisterRequest {
  username: string
  email: string
  password: string
  full_name?: string
}

/** 认证响应 */
export interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}
