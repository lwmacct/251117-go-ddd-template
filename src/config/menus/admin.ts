/**
 * 管理后台菜单配置
 *
 * 菜单项顺序决定了侧边栏的显示顺序
 * 后续可扩展权限过滤逻辑
 */
import type { MenuItem } from "@/views/Navigation/types";

export const adminMenuItems: MenuItem[] = [
  {
    title: "数据概览",
    path: "/admin/overview",
    icon: "mdi-speedometer",
  },
  {
    title: "角色管理",
    path: "/admin/roles",
    icon: "mdi-account-group",
  },
  {
    title: "用户管理",
    path: "/admin/users",
    icon: "mdi-account",
  },
  {
    title: "菜单管理",
    path: "/admin/menus",
    icon: "mdi-menu",
  },
  {
    title: "审计日志",
    path: "/admin/auditlogs",
    icon: "mdi-file-document-outline",
  },
  {
    title: "系统设置",
    path: "/admin/settings",
    icon: "mdi-cog",
  },
];
