/**
 * Register 页面相关的类型定义
 */

/**
 * 注册表单接口
 */
export interface RegisterForm {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  nickname?: string;
  captcha_id: string;
  captcha: string;
}

/**
 * 注册响应接口
 */
export interface RegisterResponse {
  success: boolean;
  message?: string;
  user?: {
    id: number;
    username: string;
    email: string;
    nickname: string;
  };
}

/**
 * Register 页面数据接口
 */
export interface RegisterPageData {
  pageTitle: string;
  pageIcon: string;
  backgroundGradient: string;
}
