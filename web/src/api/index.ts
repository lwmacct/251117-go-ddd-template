/**
 * API 统一导出
 */
export * from "./auth";
export * from "./admin";
export * from "./user";

// 同时导出常用的类型定义
export type {
  CaptchaData,
  LoginRequest,
  RegisterRequest,
  BasicLoginRequest,
  BasicRegisterRequest,
  AuthResponse,
  LoginResult,
  User,
} from "@/types";
