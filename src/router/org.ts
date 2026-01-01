import type { RouteRecordRaw } from "vue-router";

/**
 * Org 组织路由配置
 *
 * 架构说明：
 * - 组织相关功能路由（团队任务等）
 * - 使用动态路由参数 :orgId 和 :teamId
 * - 路由参数通过 props 传递给组件，确保类型安全
 */

export const orgRoutes: RouteRecordRaw = {
  path: "/org/:orgId/teams/:teamId",
  name: "OrgLayout",
  component: () => import("@/layout/OrgLayout.vue"),
  redirect: (to) => ({ ...to, name: "OrgTeamTasks" }),
  meta: {
    title: "组织工作台",
    requiresAuth: true,
  },
  children: [
    {
      path: "tasks",
      name: "OrgTeamTasks",
      component: () => import("@/pages/org/tasks/index.vue"),
      meta: {
        title: "团队任务",
        icon: "mdi-check-circle",
        menuVisible: true,
        menuOrder: 1,
      },
      // 将路由参数转换为 props
      props: (route) => ({
        orgId: Number(route.params.orgId),
        teamId: Number(route.params.teamId),
      }),
    },
  ],
};
