/**
 * 表单验证规则 Composable
 * 提供可复用的表单验证规则
 */

/** 用户名验证规则 */
export const usernameRules = [
  (v: string) => !!v || "请输入用户名",
  (v: string) =>
    (v && v.length >= 3 && v.length <= 50) || "用户名长度为 3-50 个字符",
  (v: string) =>
    /^[a-zA-Z0-9_]+$/.test(v) || "用户名只能包含字母、数字和下划线",
];

/** 邮箱验证规则 */
export const emailRules = [
  (v: string) => !!v || "请输入邮箱",
  (v: string) => /.+@.+\..+/.test(v) || "请输入有效的邮箱地址",
];

/** 密码验证规则 */
export const passwordRules = [
  (v: string) => !!v || "请输入密码",
  (v: string) => (v && v.length >= 6) || "密码至少需要 6 个字符",
];

/** 姓名验证规则 */
export const fullNameRules = [
  (v: string) => !v || v.length <= 100 || "姓名最多 100 个字符",
];

/** 登录字段验证规则（用户名或邮箱） */
export const loginFieldRules = [
  (v: string) => !!v || "请输入用户名或邮箱",
  (v: string) => v.length >= 3 || "至少需要 3 个字符",
];

/**
 * 创建确认密码验证规则
 * @param password - 原始密码的 ref
 */
export const createConfirmPasswordRules = (password: string) => [
  (v: string) => !!v || "请确认密码",
  (v: string) => v === password || "两次输入的密码不一致",
];
