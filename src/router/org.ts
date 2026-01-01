import type { RouteRecordRaw } from "vue-router";

/**
 * Org 组织路由配置
 *
 * 架构说明：
 * - 组织相关功能路由
 * - 使用动态路由参数 :orgId 和 :teamId
 * - 路由参数通过 props 传递给组件，确保类型安全
 */

export const orgRoutes: RouteRecordRaw = {
  path: "/org/:orgId/team/:teamId",
  component: () => import("@/layout/OrgLayout.vue"),
  meta: {
    title: "组织工作台",
    requiresAuth: true,
  },
  children: [
    {
      path: "",
      name: "OrgTeam",
      component: () => import("@/pages/org/team/index.vue"),
      props: (route) => ({
        orgId: Number(route.params.orgId),
        teamId: Number(route.params.teamId),
      }),
    },
  ],
};
