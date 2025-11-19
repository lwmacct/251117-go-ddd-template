/**
 * 平台认证 API (支持验证码和 2FA)
 */

import { apiClient } from "./client";
import type { PlatformLoginRequest, PlatformRegisterRequest, PlatformApiResponse, CaptchaData } from "@/types/auth";

/**
 * 平台认证 API 类
 */
export class PlatformAuthAPI {
  /**
   * 获取验证码
   */
  static async getCaptcha(): Promise<PlatformApiResponse<CaptchaData>> {
    try {
      const { data } = await apiClient.get<PlatformApiResponse<CaptchaData>>("/auth/captcha");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "获取验证码失败");
    }
  }

  /**
   * 登录
   */
  static async login(req: PlatformLoginRequest): Promise<PlatformApiResponse<any>> {
    try {
      const { data } = await apiClient.post<PlatformApiResponse<any>>("/auth/login", req);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "登录失败");
    }
  }

  /**
   * 注册
   */
  static async register(req: PlatformRegisterRequest): Promise<
    PlatformApiResponse<{
      session_token?: string;
      message?: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        PlatformApiResponse<{
          session_token?: string;
          message?: string;
        }>
      >("/auth/register", req);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "注册失败");
    }
  }

  /**
   * 验证邮箱
   */
  static async verifyEmail(params: { session_token?: string; email?: string; code: string }): Promise<PlatformApiResponse<any>> {
    try {
      const { data } = await apiClient.post<PlatformApiResponse<any>>("/auth/verify-email", params);
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "邮箱验证失败");
    }
  }

  /**
   * 重新发送验证码
   */
  static async resendVerificationCode(sessionToken: string): Promise<PlatformApiResponse<any>> {
    try {
      const { data } = await apiClient.post<PlatformApiResponse<any>>("/auth/resend-code", {
        session_token: sessionToken,
      });
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "发送验证码失败");
    }
  }

  /**
   * 验证 2FA (双因素认证)
   * 注意：2FA 验证实际上是第二次登录，使用相同的 /auth/login 端点
   */
  static async verify2FA(params: { session_token: string; code: string }): Promise<PlatformApiResponse<any>> {
    try {
      const { data } = await apiClient.post<PlatformApiResponse<any>>("/auth/login", {
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
    PlatformApiResponse<{
      secret: string;
      qrcode_url: string;
      qrcode_img: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        PlatformApiResponse<{
          secret: string;
          qrcode_url: string;
          qrcode_img: string;
        }>
      >("/auth/2fa/setup");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "设置 2FA 失败");
    }
  }

  /**
   * 验证并启用 2FA
   */
  static async enable2FA(code: string): Promise<
    PlatformApiResponse<{
      recovery_codes: string[];
      message: string;
    }>
  > {
    try {
      const { data } = await apiClient.post<
        PlatformApiResponse<{
          recovery_codes: string[];
          message: string;
        }>
      >("/auth/2fa/verify", { code });
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "启用 2FA 失败");
    }
  }

  /**
   * 禁用 2FA
   */
  static async disable2FA(): Promise<PlatformApiResponse<any>> {
    try {
      const { data } = await apiClient.post<PlatformApiResponse<any>>("/auth/2fa/disable");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "禁用 2FA 失败");
    }
  }

  /**
   * 获取 2FA 状态
   */
  static async get2FAStatus(): Promise<
    PlatformApiResponse<{
      enabled: boolean;
      recovery_codes_count: number;
    }>
  > {
    try {
      const { data } = await apiClient.get<
        PlatformApiResponse<{
          enabled: boolean;
          recovery_codes_count: number;
        }>
      >("/auth/2fa/status");
      return data;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || "获取 2FA 状态失败");
    }
  }
}
