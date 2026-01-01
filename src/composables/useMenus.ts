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
   * 标准化动态路由路径
   * 将 /org/1/teams/1/tasks 转换为 /org/:orgId/teams/:teamId/tasks
   * 用于收藏和历史记录中的统一标识
   */
  function normalizeDynamicPath(path: string): string {
    return path.replace(/\/\d+(?=\/|$)/g, "/:id").replace(/\/:id\/:id/g, "/:orgId/teams/:teamId");
  }

  /**
   * 检查是否为有效的动态路由路径
   * 如果包含所有必需参数则返回 true
   */
  function isValidDynamicPath(path: string): boolean {
    // 检查路径是否包含实际的组织/团队 ID（数字）
    return /\/org\/\d+\/teams\/\d+/.test(path);
  }

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

  /**
   * 组织工作台菜单
   * 使用模板路径作为标识，点击时由导航逻辑处理参数
   */
  const orgMenus = computed<MenuItem[]>(() => {
    const parentRoute = router.options.routes.find((r) => r.path === "/org/:orgId/teams/:teamId");
    if (!parentRoute?.children) return [];

    // 获取当前路由参数（如果用户在组织页面中）
    const currentParams = router.currentRoute.value.params;
    const hasOrgParams = currentParams.orgId && currentParams.teamId;

    return parentRoute.children
      .filter((r) => {
        if (r.meta?.menuVisible !== true) return false;
        return hasRole(r.meta?.roles as string[] | undefined);
      })
      .map((r) => {
        // 如果有当前组织参数，使用实际路径；否则使用模板路径
        let path: string;
        if (hasOrgParams) {
          path = `/org/${currentParams.orgId}/teams/${currentParams.teamId}/${r.path}`;
        } else {
          // 使用模板路径作为标识，导航时需要处理
          path = `/org/:orgId/teams/:teamId/${r.path}`;
        }

        return {
          path,
          title: (r.meta?.title as string) || "",
          icon: r.meta?.icon as string | undefined,
        };
      })
      .sort((a, b) => {
        const routeA = parentRoute.children?.find((r) => {
          const lastSegment = a.path.split("/").pop();
          return r.path === lastSegment || r.path === lastSegment?.replace(/:\w+/g, "");
        });
        const routeB = parentRoute.children?.find((r) => {
          const lastSegment = b.path.split("/").pop();
          return r.path === lastSegment || r.path === lastSegment?.replace(/:\w+/g, "");
        });
        const orderA = (routeA?.meta?.menuOrder as number) || 999;
        const orderB = (routeB?.meta?.menuOrder as number) || 999;
        return orderA - orderB;
      });
  });

  return {
    adminMenus,
    userMenus,
    orgMenus,
    normalizeDynamicPath,
    isValidDynamicPath,
  };
}
