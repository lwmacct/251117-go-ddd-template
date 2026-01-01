<!--
  OrgTeamSelector.vue - 组织/团队选择器组件

  仅在组织上下文中显示，允许用户快速切换组织和团队
  使用 CSS 媒体查询控制响应式显示
-->
<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { userOrgApi } from "@/api";
import { extractData } from "@/api/helpers";
import type { OrgUserOrgDTO, OrgUserTeamDTO } from "@models";

const route = useRoute();
const router = useRouter();

// ==================== 状态 ====================

const loading = ref(false);
const orgs = ref<OrgUserOrgDTO[]>([]);
const teams = ref<OrgUserTeamDTO[]>([]);

// ==================== 计算属性 ====================

/**
 * 判断是否处于组织上下文中
 */
const isInOrgContext = computed(() => {
  return /^\/org\/\d+\/teams\/\d+/.test(route.path);
});

/**
 * 当前组织 ID（从路由参数获取）
 * Vue Router 可能将 params 作为数组返回，需要处理
 */
const currentOrgId = computed<number | undefined>(() => {
  const orgId = route.params.orgId;
  if (!orgId) return undefined;
  // 处理数组情况：取第一个元素
  const id = Array.isArray(orgId) ? orgId[0] : orgId;
  const num = Number(id);
  return Number.isFinite(num) ? num : undefined;
});

/**
 * 当前团队 ID（从路由参数获取）
 * Vue Router 可能将 params 作为数组返回，需要处理
 */
const currentTeamId = computed<number | undefined>(() => {
  const teamId = route.params.teamId;
  if (!teamId) return undefined;
  // 处理数组情况：取第一个元素
  const id = Array.isArray(teamId) ? teamId[0] : teamId;
  const num = Number(id);
  return Number.isFinite(num) ? num : undefined;
});

/**
 * 当前组织对象
 */
const currentOrg = computed(() => {
  if (!currentOrgId.value) return undefined;
  return orgs.value.find((o) => o.id === currentOrgId.value);
});

/**
 * 当前团队对象
 */
const currentTeam = computed(() => {
  if (!currentTeamId.value) return undefined;
  return teams.value.find((t) => t.id === currentTeamId.value);
});

// ==================== 方法 ====================

/**
 * 获取用户所属组织列表
 */
async function fetchOrgs() {
  loading.value = true;
  try {
    const response = await userOrgApi.apiUserOrgsGet();
    orgs.value = extractData(response.data) || [];
  } catch {
    orgs.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * 获取指定组织的团队列表
 */
async function fetchTeams(orgId: number) {
  loading.value = true;
  try {
    const response = await userOrgApi.apiUserTeamsGet(orgId);
    teams.value = extractData(response.data) || [];
  } catch {
    teams.value = [];
  } finally {
    loading.value = false;
  }
}

/**
 * 切换组织
 * 需要先获取目标组织的团队列表，然后导航到第一个团队
 */
async function selectOrg(orgId: number) {
  // 如果目标组织的团队已加载，直接使用
  const existingTeams = teams.value.filter((t) => t.org_id === orgId);
  if (existingTeams.length > 0) {
    const firstTeam = existingTeams[0]!;
    router.push(`/org/${orgId}/teams/${firstTeam.id}/tasks`);
    return;
  }

  // 否则先获取团队列表，再导航
  loading.value = true;
  try {
    const response = await userOrgApi.apiUserTeamsGet(orgId);
    const fetchedTeams = extractData(response.data) || [];

    // 将新团队的团队合并到列表中
    fetchedTeams.forEach((team) => {
      if (!teams.value.find((t) => t.id === team.id)) {
        teams.value.push(team);
      }
    });

    if (fetchedTeams.length > 0) {
      const firstTeam = fetchedTeams[0]!;
      router.push(`/org/${orgId}/teams/${firstTeam.id}/tasks`);
    }
  } finally {
    loading.value = false;
  }
}

/**
 * 切换团队
 */
function selectTeam(teamId: number) {
  if (currentOrgId.value) {
    router.push(`/org/${currentOrgId.value}/teams/${teamId}/tasks`);
  }
}

// ==================== 生命周期与监听 ====================

// 组件挂载时获取组织列表
onMounted(() => {
  fetchOrgs();
});

// 监听当前组织变化，重新获取团队列表
watch(
  currentOrgId,
  (newOrgId) => {
    // 只有在 orgId 是有效数字时才获取团队列表
    if (newOrgId !== undefined) {
      fetchTeams(newOrgId);
    }
  },
  { immediate: true },
);
</script>

<template>
  <div v-if="isInOrgContext" class="org-team-selector">
    <!-- 加载中状态 -->
    <div v-if="!currentOrg" class="selector-container selector-loading">
      <div class="skeleton-avatar"></div>
      <v-divider vertical class="mx-1" />
      <div class="skeleton-avatar"></div>
    </div>

    <!-- 组织/团队选择器容器 -->
    <div v-else class="selector-container">
      <!-- 组织选择器 -->
      <v-menu>
        <template #activator="{ props: menuProps }">
          <v-btn variant="text" density="compact" class="selector-btn" v-bind="menuProps">
            <v-avatar v-if="currentOrg.avatar" :image="currentOrg.avatar" size="20" />
            <v-icon v-else size="20">mdi-domain</v-icon>
            <span class="selector-text">{{ currentOrg.display_name || currentOrg.name }}</span>
            <v-icon size="16" class="dropdown-arrow">mdi-menu-down</v-icon>
          </v-btn>
        </template>
        <v-list density="compact">
          <v-list-item v-for="org in orgs" :key="org.id" :active="org.id === currentOrgId" @click="selectOrg(org.id!)">
            <template #prepend>
              <v-avatar v-if="org.avatar" :image="org.avatar" size="24" />
              <v-icon v-else>mdi-domain</v-icon>
            </template>
            <v-list-item-title>{{ org.display_name || org.name }}</v-list-item-title>
            <v-list-item-subtitle v-if="org.role">{{ org.role }}</v-list-item-subtitle>
          </v-list-item>
          <v-divider v-if="orgs.length === 0" />
          <v-list-item v-if="orgs.length === 0" disabled>
            <v-list-item-title>暂无组织</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>

      <!-- 分隔符 -->
      <v-divider vertical class="mx-1" />

      <!-- 团队选择器 -->
      <v-menu :disabled="!currentOrgId">
        <template #activator="{ props: menuProps }">
          <v-btn variant="text" density="compact" class="selector-btn" :disabled="!currentOrgId" v-bind="menuProps">
            <v-avatar v-if="currentTeam?.avatar" :image="currentTeam.avatar" size="20" />
            <v-icon v-else size="20">mdi-account-group</v-icon>
            <span class="selector-text">{{ currentTeam?.display_name || currentTeam?.name || "选择团队" }}</span>
            <v-icon size="16" class="dropdown-arrow">mdi-menu-down</v-icon>
          </v-btn>
        </template>
        <v-list density="compact">
          <v-list-item
            v-for="team in teams.filter((t) => t.org_id === currentOrgId)"
            :key="team.id"
            :active="team.id === currentTeamId"
            @click="selectTeam(team.id!)"
          >
            <template #prepend>
              <v-avatar v-if="team.avatar" :image="team.avatar" size="24" />
              <v-icon v-else>mdi-account-group</v-icon>
            </template>
            <v-list-item-title>{{ team.display_name || team.name }}</v-list-item-title>
            <v-list-item-subtitle v-if="team.role">{{ team.role }}</v-list-item-subtitle>
          </v-list-item>
          <v-divider v-if="teams.filter((t) => t.org_id === currentOrgId).length === 0" />
          <v-list-item v-if="teams.filter((t) => t.org_id === currentOrgId).length === 0" disabled>
            <v-list-item-title>该组织暂无团队</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<style scoped>
.org-team-selector {
  display: flex;
  align-items: center;
  margin-left: 8px;
}

/* 容器样式：参考 SearchBox，统一高度和边框 */
.selector-container {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px 10px;
  border-radius: 6px;
  background-color: rgba(var(--v-theme-surface), 1);
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  transition: all 0.2s ease;
  /* 限制高度与 SearchBox 一致 */
  min-height: 36px;
}

.selector-container:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
  border-color: rgb(var(--v-theme-primary));
}

.selector-btn {
  padding: 0 4px;
  min-width: 28px;
  height: auto;
  border-radius: 4px;
  background-color: transparent;
  transition: background-color 0.2s;
  color: rgb(var(--v-theme-on-surface));
}

.selector-btn:hover {
  background-color: rgba(var(--v-theme-on-surface), 0.1);
}

.dropdown-arrow {
  margin-left: 2px;
  opacity: 0.7;
  color: rgb(var(--v-theme-on-surface));
}

.selector-text {
  display: none;
  /* 默认隐藏名称 */
  font-size: 0.875rem;
  max-width: 150px;
  margin-left: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(var(--v-theme-on-surface));
}

/* 宽屏时显示名称 */
@media (min-width: 900px) {
  .selector-text {
    display: inline;
  }
}

/* 分隔符颜色 */
.selector-container .v-divider {
  border-color: rgba(var(--v-border-color), var(--v-border-opacity));
  opacity: 0.6;
}

/* 加载状态骨架屏 */
.skeleton-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: rgba(var(--v-theme-on-surface), 0.12);
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.4;
  }
}
</style>
