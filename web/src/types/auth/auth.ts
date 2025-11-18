/**
 * 认证相关类型定义
 */
import type { User } from "./user";

/** 登录请求 (基础)  */
export interface LoginRequest {
  login: string; // 用户名或邮箱
  password: string;
}

/** 登录请求 (平台版，带验证码)  */
export interface PlatformLoginRequest {
  account: string; // 手机号/用户名/邮箱
  password: string;
  captcha_id: string;
  captcha: string;
}

/** 注册请求 (基础)  */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name?: string;
}

/** 注册请求 (平台版，带验证码)  */
export interface PlatformRegisterRequest {
  email: string;
  password: string;
  captcha_id: string;
  captcha: string;
}

/** 验证码数据 */
export interface CaptchaData {
  id: string;
  image: string; // Base64 编码的图片
}

/** 认证响应 */
export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

/** 登录结果 (支持 2FA)  */
export interface LoginResult {
  success: boolean;
  message?: string;
  requiresTwoFactor: boolean;
  sessionToken?: string;
}

/** 平台 API 响应 */
export interface PlatformApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}
