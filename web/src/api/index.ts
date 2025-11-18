/**
 * API 统一导出
 */
export * from './auth'

// 同时导出常用的类型定义
export type {
  CaptchaData,
  LoginRequest,
  RegisterRequest,
  PlatformLoginRequest,
  PlatformRegisterRequest,
  AuthResponse,
  LoginResult,
  User,
} from '@/types'
