/**
 * 认证存储工具
 */

/** Token 存储键名 */
const TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const TOKEN_EXPIRY_KEY = "token_expiry";

/**
 * 保存访问令牌
 */
export const saveAccessToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token);
};

/**
 * 获取访问令牌
 */
export const getAccessToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

/**
 * 保存刷新令牌
 */
export const saveRefreshToken = (token: string): void => {
  localStorage.setItem(REFRESH_TOKEN_KEY, token);
};

/**
 * 获取刷新令牌
 */
export const getRefreshToken = (): string | null => {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
};

/**
 * 保存令牌过期时间
 */
export const saveTokenExpiry = (expiry: number): void => {
  localStorage.setItem(TOKEN_EXPIRY_KEY, expiry.toString());
};

/**
 * 获取令牌过期时间
 */
export const getTokenExpiry = (): number | null => {
  const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
  return expiry ? parseInt(expiry, 10) : null;
};

/**
 * 清除所有认证令牌
 */
export const clearAuthTokens = (): void => {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(TOKEN_EXPIRY_KEY);
};

/**
 * 检查是否有访问令牌
 */
export const hasAccessToken = (): boolean => {
  return !!getAccessToken();
};
