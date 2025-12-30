/**
 * 角色权限检查 Composable
 *
 * 提供基于角色的权限判断，用于路由守卫和菜单过滤
 */
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";

export function usePermission() {
  const authStore = useAuthStore();

  /**
   * 当前用户的角色名称列表
   * 从 user.roles 中提取 name 字段
   */
  const roleNames = computed<string[]>(() => {
    const roles = authStore.currentUser?.roles;
    if (!roles) return [];
    return roles.map((r) => r.name).filter((name): name is string => !!name);
  });

  /**
   * 检查用户是否拥有任一指定角色
   * @param requiredRoles 允许访问的角色列表，为空表示所有角色可访问
   */
  function hasRole(requiredRoles?: string[]): boolean {
    // 未指定角色限制，所有已登录用户可访问
    if (!requiredRoles || requiredRoles.length === 0) return true;
    // 用户角色与要求角色有交集即可
    return requiredRoles.some((r) => roleNames.value.includes(r));
  }

  /**
   * 检查用户是否拥有所有指定角色
   * @param requiredRoles 必须拥有的角色列表
   */
  function hasAllRoles(requiredRoles: string[]): boolean {
    if (requiredRoles.length === 0) return true;
    return requiredRoles.every((r) => roleNames.value.includes(r));
  }

  /**
   * 是否为管理员角色
   */
  const isAdmin = computed(() => roleNames.value.includes("admin"));

  return {
    roleNames,
    hasRole,
    hasAllRoles,
    isAdmin,
  };
}
