/**
 * 验证工具函数
 */

/**
 * 验证邮箱格式
 * @param email - 邮箱地址
 * @returns 是否有效
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * 验证用户名格式
 * @param username - 用户名
 * @returns 是否有效
 */
export function isValidUsername(username: string): boolean {
  // 3-50 个字符，只能包含字母、数字和下划线
  const usernameRegex = /^[a-zA-Z0-9_]{3,50}$/;
  return usernameRegex.test(username);
}

/**
 * 验证密码强度
 * @param password - 密码
 * @returns 强度等级 (weak, medium, strong)
 */
export function checkPasswordStrength(
  password: string
): "weak" | "medium" | "strong" {
  if (password.length < 6) {
    return "weak";
  }

  let strength = 0;

  // 包含小写字母
  if (/[a-z]/.test(password)) strength++;
  // 包含大写字母
  if (/[A-Z]/.test(password)) strength++;
  // 包含数字
  if (/\d/.test(password)) strength++;
  // 包含特殊字符
  if (/[^a-zA-Z0-9]/.test(password)) strength++;
  // 长度大于 8
  if (password.length >= 8) strength++;

  if (strength <= 2) return "weak";
  if (strength <= 3) return "medium";
  return "strong";
}
