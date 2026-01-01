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
   * @param useCurrentParams 是否使用当前路由参数（用于动态路由）
   */
  function extractMenusFromRoute(parentPath: string, useCurrentParams = false): MenuItem[] {
    const parentRoute = router.options.routes.find((r) => r.path === parentPath);
    if (!parentRoute?.children) return [];

    // 获取当前路由参数（用于动态路由）
    const currentParams = useCurrentParams ? router.currentRoute.value.params : {};

    return parentRoute.children
      .filter((r) => {
        // 只显示设置了 menuVisible: true 的路由
        if (r.meta?.menuVisible !== true) return false;
        // 按角色过滤
        return hasRole(r.meta?.roles as string[] | undefined);
      })
      .map((r) => {
        // 构建路径：对于动态路由，使用当前路由参数
        let path = `${parentPath}/${r.path}`;

        // 替换路径中的动态参数
        if (useCurrentParams && Object.keys(currentParams).length > 0) {
          Object.entries(currentParams).forEach(([key, value]) => {
            path = path.replace(`:${key}`, String(value));
          });
        }

        return {
          path,
          title: (r.meta?.title as string) || "",
          icon: r.meta?.icon as string | undefined,
        };
      })
      .sort((a, b) => {
        // 按 menuOrder 排序（需要从原始路由获取）
        const routeA = parentRoute.children?.find((r) => r.path === a.path.split("/").pop());
        const routeB = parentRoute.children?.find((r) => r.path === b.path.split("/").pop());
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
