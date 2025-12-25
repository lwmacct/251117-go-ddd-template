/**
 * 用户中心菜单配置
 */
import type { MenuItem } from "@/views/Navigation/types";

export const userMenuItems: MenuItem[] = [
  {
    title: "个人资料",
    path: "/user/profile",
    icon: "mdi-account-circle",
  },
  {
    title: "安全设置",
    path: "/user/security",
    icon: "mdi-shield-lock",
  },
  {
    title: "访问令牌",
    path: "/user/tokens",
    icon: "mdi-key-variant",
  },
];
