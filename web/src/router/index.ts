import { createRouter, createWebHashHistory } from "vue-router";
import { adminRoutes } from "./admin";
import { authRoutes } from "./auth";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/auth/login",
    },
    authRoutes,
    adminRoutes,
  ],
});

export default router;
