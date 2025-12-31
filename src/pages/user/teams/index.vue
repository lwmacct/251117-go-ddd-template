<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { userOrgApi, extractData } from "@/api";
import type { OrgUserTeamDTO } from "@models";
import { useSnackbar } from "@/composables";

const { error } = useSnackbar();
const teams = ref<OrgUserTeamDTO[]>([]);
const loading = ref(false);
const orgFilter = ref<number | null>(null);

/**
 * 获取我的团队列表
 */
const fetchMyTeams = async () => {
  loading.value = true;
  try {
    const response = await userOrgApi.apiUserTeamsGet(orgFilter.value || undefined);
    teams.value = extractData<OrgUserTeamDTO[]>(response.data) || [];
  } catch (err) {
    error((err as Error).message || "获取团队列表失败");
  } finally {
    loading.value = false;
  }
};

/**
 * 按组织分组展示
 */
const groupedTeams = computed(() => {
  const groups: Record<string, OrgUserTeamDTO[]> = {};
  teams.value.forEach((team) => {
    const orgName = team.org_name || "未分类";
    if (!groups[orgName]) {
      groups[orgName] = [];
    }
    groups[orgName].push(team);
  });
  return groups;
});

/**
 * 可用的组织筛选选项
 */
const orgOptions = computed(() => {
  const orgSet = new Set(teams.value.map((t) => t.org_id));
  return Array.from(orgSet).map((orgId) => {
    const team = teams.value.find((t) => t.org_id === orgId);
    return {
      id: orgId,
      name: team?.org_name || `组织 ${orgId}`,
    };
  });
});

/**
 * 获取角色文本
 */
const getRoleText = (role?: string) => {
  const map: Record<string, string> = {
    lead: "负责人",
    member: "成员",
  };
  return map[role as keyof typeof map] || role || "-";
};

/**
 * 获取角色颜色
 */
const getRoleColor = (role?: string) => {
  const map: Record<string, string> = {
    lead: "warning",
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

onMounted(fetchMyTeams);
</script>

<template>
  <div class="my-teams-page">
    <v-row>
      <v-col cols="12">
        <div class="d-flex align-center mb-6">
          <h1 class="text-h4">我的团队</h1>
          <v-spacer />
          <v-select
            v-if="orgOptions.length > 1"
            v-model="orgFilter"
            label="筛选组织"
            :items="orgOptions"
            item-title="name"
            item-value="id"
            clearable
            variant="outlined"
            density="compact"
            hide-details
            style="max-width: 200px"
            @update:model-value="fetchMyTeams"
          >
            <template #prepend-inner>
              <v-icon icon="mdi-filter" />
            </template>
          </v-select>
        </div>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" />
        <p class="text-body-2 mt-4 text-medium-emphasis">加载中...</p>
      </v-col>
    </v-row>

    <template v-else-if="Object.keys(groupedTeams).length > 0">
      <v-row v-for="(teamList, orgName) in groupedTeams" :key="orgName">
        <v-col cols="12">
          <v-card>
            <v-card-title class="bg-surface-light">
              <v-icon icon="mdi-office-building" start />
              {{ orgName }}
            </v-card-title>
            <v-card-text class="pt-4">
              <v-row>
                <v-col v-for="team in teamList" :key="team.id" cols="12" sm="6" md="4">
                  <v-card variant="outlined" class="h-100">
                    <v-card-item>
                      <v-card-title>{{ team.display_name || team.name }}</v-card-title>
                      <v-card-subtitle>{{ team.name }}</v-card-subtitle>
                    </v-card-item>

                    <v-card-text>
                      <p v-if="team.description" class="text-body-2 mb-2">
                        {{ team.description }}
                      </p>
                      <p v-else class="text-body-2 text-medium-emphasis mb-2">暂无描述</p>

                      <v-chip :color="getRoleColor(team.role)" size="small">
                        <v-icon start size="small">mdi-shield-account</v-icon>
                        {{ getRoleText(team.role) }}
                      </v-chip>

                      <div class="d-flex align-center mt-2 text-body-2 text-medium-emphasis">
                        <v-icon icon="mdi-calendar" size="small" class="mr-1" />
                        <span>加入于 {{ formatDate(team.joined_at) }}</span>
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>

    <v-row v-else>
      <v-col cols="12">
        <v-card>
          <v-card-text class="text-center py-8">
            <v-icon icon="mdi-account-group" size="64" color="disabled" class="mb-4" />
            <p class="text-h6 text-medium-emphasis">你还没有加入任何团队</p>
            <p class="text-body-2 text-medium-emphasis mt-2">联系团队管理员邀请你加入</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.my-teams-page {
  width: 100%;
}
</style>
