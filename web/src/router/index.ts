import { createRouter, createWebHashHistory } from "vue-router";
import { adminRoutes } from "@/pages/admin/router";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/admin/overview",
    },
    adminRoutes,
  ],
});

export default router;
