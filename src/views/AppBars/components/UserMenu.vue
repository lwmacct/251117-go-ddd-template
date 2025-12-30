<script setup lang="ts">
/**
 * 用户菜单组件
 *
 * 功能：
 * - 用户头像（显示首字母）
 * - 动态菜单项（用户中心菜单）
 * - 登出/登录按钮
 */
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore, useNavbarStore } from "@/stores";

const authStore = useAuthStore();
const navbarStore = useNavbarStore();
const router = useRouter();
const menuOpen = ref(false);

/** 显示名称 */
const displayName = computed(() => {
  return authStore.currentUser?.username || "未登录";
});

/** 用户头像 */
const userAvatar = computed(() => {
  return undefined; // 简化版不存储头像
});

/** 用户名首字母 */
const userInitial = computed(() => {
  const name = authStore.currentUser?.username || "?";
  return name.charAt(0).toUpperCase();
});

/** 用户中心菜单项 */
const menuItems = computed(() => {
  return navbarStore.allMenuItems
    .filter((item) => item.category === "用户中心")
    .map((item) => ({
      title: item.title,
      icon: item.icon,
      path: item.path,
    }));
});

/** 导航到指定页面 */
function navigateTo(path: string) {
  menuOpen.value = false;
  navbarStore.handleNavigation();
  router.push(path);
}

/** 退出登录 */
async function handleLogout() {
  menuOpen.value = false;
  authStore.logout();
  // 可选：清除访问历史
  navbarStore.clearHistory();
  router.push("/auth/login");
}
</script>

<template>
  <div class="user-menu">
    <v-menu v-model="menuOpen" :close-on-content-click="false" location="bottom end" offset="8">
      <template #activator="{ props }">
        <v-btn icon v-bind="props" class="user-avatar-btn">
          <v-avatar v-if="authStore.isAuthenticated" color="primary" size="36">
            <v-img v-if="userAvatar" :src="userAvatar" :alt="displayName" />
            <span v-else class="text-white font-weight-bold">{{ userInitial }}</span>
          </v-avatar>
          <v-icon v-else>mdi-account-circle</v-icon>
        </v-btn>
      </template>

      <v-card min-width="280" elevation="8">
        <!-- 用户信息头部 -->
        <v-card-item v-if="authStore.isAuthenticated">
          <template #prepend>
            <v-avatar color="primary" size="48">
              <v-img v-if="userAvatar" :src="userAvatar" :alt="displayName" />
              <span v-else class="text-h6 font-weight-bold">{{ userInitial }}</span>
            </v-avatar>
          </template>

          <v-card-title class="text-subtitle-1">{{ displayName }}</v-card-title>
        </v-card-item>

        <v-divider v-if="authStore.isAuthenticated" />

        <!-- 菜单列表 -->
        <v-list v-if="authStore.isAuthenticated" density="compact" nav>
          <v-list-item
            v-for="item in menuItems"
            :key="item.path"
            :prepend-icon="item.icon"
            :title="item.title"
            link
            @click="navigateTo(item.path)"
          />
        </v-list>

        <v-divider v-if="authStore.isAuthenticated" />

        <!-- 退出登录按钮 -->
        <v-card-actions v-if="authStore.isAuthenticated">
          <v-btn block color="error" variant="text" prepend-icon="mdi-logout" @click="handleLogout"> 退出登录 </v-btn>
        </v-card-actions>

        <!-- 未登录状态 -->
        <v-card-actions v-else>
          <v-btn block color="primary" variant="flat" @click="navigateTo('/auth/login')"> 立即登录 </v-btn>
        </v-card-actions>
      </v-card>
    </v-menu>
  </div>
</template>

<style scoped>
.user-menu {
  display: flex;
  align-items: center;
}

.user-avatar-btn {
  transition: transform 0.2s ease;
}

.user-avatar-btn:hover {
  transform: scale(1.05);
}
</style>
