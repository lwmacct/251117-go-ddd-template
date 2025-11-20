/**
 * 认证 API（标准版）
 * 支持验证码、2FA 等完整认证功能
 */

import { apiClient } from "./client";
import type { LoginRequest, RegisterRequest, CaptchaData } from "@/types/auth";
import type { ApiResponse } from "@/types/response";

/**
 * 认证 API
 * 提供完整的认证功能：登录、注册、验证码、2FA 等
 */
export class AuthAPI {
  /**
   * 获取验证码
   */
  static async getCaptcha(): Promise<ApiResponse<CaptchaData>> {
    try {
      const { data } = await apiClient.get<ApiResponse<CaptchaData>>("/api/auth/captcha");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "获取验证码失败");
    }
  }

  /**
   * 登录（带验证码）
   */
  static async login(req: LoginRequest): Promise<ApiResponse<any>> {
    try {
      const { data } = await apiClient.post<ApiResponse<any>>("/api/auth/login", req);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "登录失败");
    }
  }

  /**
   * 注册（带验证码）
   */
  static async register(req: RegisterRequest): Promise<
    ApiResponse<{
      session_token?: string;
      message?: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        ApiResponse<{
          session_token?: string;
          message?: string;
        }>
      >("/api/auth/register", req);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "注册失败");
    }
  }

  /**
   * 验证邮箱
   */
  static async verifyEmail(params: { session_token?: string; email?: string; code: string }): Promise<ApiResponse<any>> {
    try {
      const { data } = await apiClient.post<ApiResponse<any>>("/api/auth/verify-email", params);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "邮箱验证失败");
    }
  }

  /**
   * 重新发送验证码
   */
  static async resendVerificationCode(sessionToken: string): Promise<ApiResponse<any>> {
    try {
      const { data } = await apiClient.post<ApiResponse<any>>("/api/auth/resend-code", {
        session_token: sessionToken,
      });
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "发送验证码失败");
    }
  }

  /**
   * 验证 2FA (双因素认证)
   * 注意：2FA 验证实际上是第二次登录，使用相同的 /login 端点
   */
  static async verify2FA(params: { session_token: string; code: string }): Promise<ApiResponse<any>> {
    try {
      const { data } = await apiClient.post<ApiResponse<any>>("/api/auth/login", {
        session_token: params.session_token,
        two_factor_code: params.code,
      });
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "2FA 验证失败");
    }
  }

  /**
   * 设置 2FA（生成密钥和二维码）
   */
  static async setup2FA(): Promise<
    ApiResponse<{
      secret: string;
      qrcode_url: string;
      qrcode_img: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        ApiResponse<{
          secret: string;
          qrcode_url: string;
          qrcode_img: string;
        }>
      >("/api/auth/2fa/setup");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "设置 2FA 失败");
    }
  }

  /**
   * 验证并启用 2FA
   */
  static async enable2FA(code: string): Promise<
    ApiResponse<{
      recovery_codes: string[];
      message: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        ApiResponse<{
          recovery_codes: string[];
          message: string;
        }>
      >("/api/auth/2fa/verify", { code });
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "启用 2FA 失败");
    }
  }

  /**
   * 禁用 2FA
   */
  static async disable2FA(): Promise<ApiResponse<any>> {
    try {
      const { data } = await apiClient.post<ApiResponse<any>>("/api/auth/2fa/disable");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "禁用 2FA 失败");
    }
  }

  /**
   * 获取 2FA 状态
   */
  static async get2FAStatus(): Promise<
    ApiResponse<{
      enabled: boolean;
      recovery_codes_count: number;
    }>
  > {
    try {
      const { data } = await apiClient.get<
        ApiResponse<{
          enabled: boolean;
          recovery_codes_count: number;
        }>
      >("/api/auth/2fa/status");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "获取 2FA 状态失败");
    }
  }
}

/**
 * @deprecated 使用 AuthAPI 代替
 * 向后兼容的别名
 */
export const PlatformAuthAPI = AuthAPI;
