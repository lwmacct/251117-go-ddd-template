<script setup lang="ts">
import { ref, onMounted } from "vue";
import { userOrgApi, extractData } from "@/api";
import type { OrgUserOrgDTO } from "@models";
import { useSnackbar } from "@/composables";

const { error } = useSnackbar();
const orgs = ref<OrgUserOrgDTO[]>([]);
const loading = ref(false);

/**
 * 获取我的组织列表
 */
const fetchMyOrgs = async () => {
  loading.value = true;
  try {
    const response = await userOrgApi.apiUserOrgsGet();
    orgs.value = extractData<OrgUserOrgDTO[]>(response.data) || [];
  } catch (err) {
    error((err as Error).message || "获取组织列表失败");
  } finally {
    loading.value = false;
  }
};

/**
 * 获取角色文本
 */
const getRoleText = (role?: string) => {
  const map: Record<string, string> = {
    owner: "所有者",
    admin: "管理员",
    member: "成员",
  };
  return map[role as keyof typeof map] || role || "-";
};

/**
 * 获取角色颜色
 */
const getRoleColor = (role?: string) => {
  const map: Record<string, string> = {
    owner: "error",
    admin: "warning",
    member: "default",
  };
  return map[role as keyof typeof map] || "default";
};

/**
 * 格式化日期
 */
const formatDate = (dateStr?: string) => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
};

onMounted(fetchMyOrgs);
</script>

<template>
  <div class="my-orgs-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">我的组织</h1>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" />
        <p class="text-body-2 mt-4 text-medium-emphasis">加载中...</p>
      </v-col>
    </v-row>

    <v-row v-else-if="orgs.length > 0">
      <v-col v-for="org in orgs" :key="org.id" cols="12" md="6" lg="4">
        <v-card hover>
          <v-card-item>
            <template #prepend>
              <v-avatar color="primary" size="48">
                <span class="text-h5">{{ org.display_name?.[0] || org.name?.[0] || "?" }}</span>
              </v-avatar>
            </template>
            <v-card-title>{{ org.display_name || org.name }}</v-card-title>
            <v-card-subtitle>{{ org.name }}</v-card-subtitle>
          </v-card-item>

          <v-card-text>
            <div v-if="org.description" class="text-body-2 mb-4">
              {{ org.description }}
            </div>
            <div v-else class="text-body-2 text-medium-emphasis mb-4">暂无描述</div>

            <v-chip :color="getRoleColor(org.role)" size="small" class="mb-2">
              <v-icon start size="small">mdi-shield-account</v-icon>
              {{ getRoleText(org.role) }}
            </v-chip>

            <div class="d-flex align-center mt-2 text-body-2 text-medium-emphasis">
              <v-icon icon="mdi-calendar" size="small" class="mr-1" />
              <span>加入于 {{ formatDate(org.joined_at) }}</span>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col cols="12">
        <v-card>
          <v-card-text class="text-center py-8">
            <v-icon icon="mdi-office-building" size="64" color="disabled" class="mb-4" />
            <p class="text-h6 text-medium-emphasis">你还没有加入任何组织</p>
            <p class="text-body-2 text-medium-emphasis mt-2">联系组织管理员邀请你加入，或创建一个新的组织</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.my-orgs-page {
  width: 100%;
}
</style>
