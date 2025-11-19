<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useOverview } from "./composables/useOverview";

const { stats, loading, errorMessage, fetchStats } = useOverview();

// 统计卡片配置
const statsCards = computed(() => [
  {
    title: "总用户数",
    value: stats.value?.total_users || 0,
    icon: "mdi-account-group",
    color: "primary",
  },
  {
    title: "活跃用户",
    value: stats.value?.active_users || 0,
    icon: "mdi-account-check",
    color: "success",
  },
  {
    title: "总角色数",
    value: stats.value?.total_roles || 0,
    icon: "mdi-shield-account",
    color: "info",
  },
  {
    title: "总权限数",
    value: stats.value?.total_permissions || 0,
    icon: "mdi-key",
    color: "warning",
  },
  {
    title: "禁用用户",
    value: stats.value?.inactive_users || 0,
    icon: "mdi-account-off",
    color: "warning",
  },
  {
    title: "封禁用户",
    value: stats.value?.banned_users || 0,
    icon: "mdi-account-cancel",
    color: "error",
  },
  {
    title: "菜单数量",
    value: stats.value?.total_menus || 0,
    icon: "mdi-menu",
    color: "purple",
  },
]);

onMounted(() => {
  fetchStats();
});
</script>

<template>
  <div class="overview-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">系统概览</h1>
        <p class="text-body-2 text-medium-emphasis mb-6">查看系统整体运行状况和关键指标</p>
      </v-col>
    </v-row>

    <v-row v-if="errorMessage">
      <v-col cols="12">
        <v-alert type="error" closable @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>

    <!-- 统计卡片 -->
    <v-row>
      <v-col v-for="stat in statsCards" :key="stat.title" cols="12" sm="6" md="3">
        <v-card hover>
          <v-card-text>
            <div class="d-flex align-center">
              <v-avatar :color="stat.color" size="56" class="mr-4">
                <v-icon :color="'white'" size="32">{{ stat.icon }}</v-icon>
              </v-avatar>
              <div>
                <div class="text-h4 font-weight-bold">
                  {{ stat.value }}
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  {{ stat.title }}
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 快速操作 -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon start>mdi-lightning-bolt</v-icon>
            快速操作
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="3">
                <v-btn block color="primary" to="/admin/users" variant="tonal">
                  <v-icon start>mdi-account-plus</v-icon>
                  管理用户
                </v-btn>
              </v-col>
              <v-col cols="12" md="3">
                <v-btn block color="info" to="/admin/roles" variant="tonal">
                  <v-icon start>mdi-shield-plus</v-icon>
                  管理角色
                </v-btn>
              </v-col>
              <v-col cols="12" md="3">
                <v-btn block color="warning" to="/admin/menus" variant="tonal">
                  <v-icon start>mdi-menu</v-icon>
                  管理菜单
                </v-btn>
              </v-col>
              <v-col cols="12" md="3">
                <v-btn block color="success" to="/admin/audit-logs" variant="tonal">
                  <v-icon start>mdi-file-document</v-icon>
                  审计日志
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 系统信息 -->
    <v-row class="mt-4">
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon start>mdi-information</v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title>系统名称</v-list-item-title>
                <v-list-item-subtitle>Go DDD Template</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>版本</v-list-item-title>
                <v-list-item-subtitle>v1.0.0</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>架构模式</v-list-item-title>
                <v-list-item-subtitle>DDD (领域驱动设计)</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon start>mdi-chart-line</v-icon>
            系统状态
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title>数据库连接</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip size="small" color="success">正常</v-chip>
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Redis 连接</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip size="small" color="success">正常</v-chip>
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>服务状态</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip size="small" color="success">运行中</v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.overview-page {
  width: 100%;
}
</style>
