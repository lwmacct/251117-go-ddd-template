import { createRouter, createWebHashHistory } from "vue-router";
import { adminRoutes } from "./admin";
import { authRoutes } from "./auth";
import { userRoutes } from "./user";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/auth/login",
    },
    authRoutes,
    adminRoutes,
    userRoutes,
  ],
});

export default router;
