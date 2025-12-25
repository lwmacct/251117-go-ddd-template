/**
 * Token 管理工具
 */

/**
 * 检查 token 是否即将过期
 * @param expiresIn - 过期时间 (秒)
 * @param threshold - 阈值 (秒) ，默认 5 分钟
 * @returns 是否即将过期
 */
export function isTokenExpiringSoon(expiresIn: number, threshold: number = 300): boolean {
  return expiresIn <= threshold;
}

/**
 * 从 JWT token 中解析 payload
 * @param token - JWT token
 * @returns payload 对象
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function parseJwtToken(token: string): any {
  try {
    const parts = token.split(".");
    if (parts.length !== 3 || !parts[1]) {
      throw new Error("Invalid JWT token format");
    }

    const base64Url = parts[1]!;
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join(""),
    );
    return JSON.parse(jsonPayload);
  } catch (error) {
    console.error("Failed to parse JWT token:", error);
    return null;
  }
}
