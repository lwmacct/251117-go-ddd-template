/**
 * 动态菜单 Composable
 *
 * 从路由配置派生菜单数据，根据用户角色过滤可见菜单
 * 实现 "路由即菜单" 设计模式，单一数据源
 */
import { computed } from "vue";
import { useRouter } from "vue-router";
import { usePermission } from "./usePermission";
import type { MenuItem } from "@/views/Navigation/types";

export function useMenus() {
  const router = useRouter();
  const { hasRole } = usePermission();

  /**
   * 从路由配置提取菜单项
   * @param parentPath 父路由路径
   */
  function extractMenusFromRoute(parentPath: string): MenuItem[] {
    const parentRoute = router.options.routes.find((r) => r.path === parentPath);
    if (!parentRoute?.children) return [];

    return parentRoute.children
      .filter((r) => {
        // 只显示设置了 menuVisible: true 的路由
        if (r.meta?.menuVisible !== true) return false;
        // 按角色过滤
        return hasRole(r.meta?.roles as string[] | undefined);
      })
      .map((r) => ({
        path: `${parentPath}/${r.path}`,
        title: (r.meta?.title as string) || "",
        icon: r.meta?.icon as string | undefined,
      }))
      .sort((a, b) => {
        // 按 menuOrder 排序（需要从原始路由获取）
        const routeA = parentRoute.children?.find((r) => `${parentPath}/${r.path}` === a.path);
        const routeB = parentRoute.children?.find((r) => `${parentPath}/${r.path}` === b.path);
        const orderA = (routeA?.meta?.menuOrder as number) || 999;
        const orderB = (routeB?.meta?.menuOrder as number) || 999;
        return orderA - orderB;
      });
  }

  /**
   * 管理后台菜单
   */
  const adminMenus = computed<MenuItem[]>(() => extractMenusFromRoute("/admin"));

  /**
   * 用户中心菜单
   */
  const userMenus = computed<MenuItem[]>(() => extractMenusFromRoute("/user"));

  return {
    adminMenus,
    userMenus,
  };
}
