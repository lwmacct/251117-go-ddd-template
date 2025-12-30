<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useAuditLogs } from "./composables/useAuditLogs";
import { ITEMS_PER_PAGE_OPTIONS } from "@/composables";
import type { AuditAuditLogDTO } from "@models";

/**
 * 审计日志管理页面
 * 用于查看系统操作日志和审计记录
 */

// 使用 composable
const {
  logs,
  loading,
  exporting,
  filters,
  pagination,
  actionOptions,
  categoryOptions,
  fetchFilterOptions,
  applyFilters,
  resetFilters,
  onTableOptionsUpdate,
  exportLogs,
} = useAuditLogs();

// 页面初始化
onMounted(() => {
  fetchFilterOptions();
});

// 对话框状态
const detailDialog = ref(false);
const selectedLog = ref<AuditAuditLogDTO | null>(null);

// 表头配置
const headers = [
  { title: "ID", key: "id", sortable: false, width: "80px" },
  { title: "用户ID", key: "user_id", sortable: false, width: "100px" },
  { title: "操作", key: "action", sortable: false, width: "120px" },
  { title: "资源", key: "resource", sortable: false, width: "120px" },
  { title: "状态", key: "status", sortable: false, width: "100px" },
  { title: "IP地址", key: "ip_address", sortable: false },
  { title: "时间", key: "created_at", sortable: false, width: "180px" },
  { title: "操作", key: "actions", sortable: false, width: "100px" },
];

// 操作类型选项（带"全部"）
const actionSelectOptions = computed(() => [{ action: "", label: "全部" }, ...actionOptions.value]);

// 资源分类选项（带"全部"）
const categorySelectOptions = computed(() => [{ value: "", label: "全部" }, ...categoryOptions.value]);

// 状态选项（固定值）
const statusOptions = [
  { value: "", text: "全部" },
  { value: "success", text: "成功" },
  { value: "failure", text: "失败" },
];

// 构建 action 和 category 的 label 映射
const actionLabelMap = computed(() => {
  const map = new Map<string, string>();
  for (const opt of actionOptions.value) {
    map.set(opt.action ?? "", opt.label ?? opt.action ?? "");
  }
  return map;
});

const categoryLabelMap = computed(() => {
  const map = new Map<string, string>();
  for (const opt of categoryOptions.value) {
    map.set(opt.value ?? "", opt.label ?? opt.value ?? "");
  }
  return map;
});

// 查看详情
const viewDetail = (log: AuditAuditLogDTO) => {
  selectedLog.value = log;
  detailDialog.value = true;
};

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
};

// 格式化操作类型
const formatAction = (action?: string) => {
  if (!action) return "-";
  return actionLabelMap.value.get(action) || action;
};

// 格式化资源类型
const formatResource = (resource?: string) => {
  if (!resource) return "-";
  return categoryLabelMap.value.get(resource) || resource;
};

// 状态颜色
const getStatusColor = (status?: string) => {
  return status === "success" ? "success" : "error";
};

// 状态文本
const getStatusText = (status?: string) => {
  return status === "success" ? "成功" : "失败";
};
</script>

<template>
  <div class="auditlogs-page">
    <!-- 标题 -->
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">审计日志</h1>
      </v-col>
    </v-row>

    <!-- 日志列表卡片 -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row>
              <v-col cols="12">
                <h3 class="text-h6">过滤条件</h3>
              </v-col>
            </v-row>
          </v-card-title>

          <!-- 过滤条件 -->
          <v-card-text>
            <v-row dense>
              <v-col cols="12" md="3">
                <v-text-field
                  v-model.number="filters.user_id"
                  label="用户ID"
                  type="number"
                  variant="outlined"
                  density="compact"
                  hide-details
                  clearable
                ></v-text-field>
              </v-col>

              <v-col cols="12" md="2">
                <v-select
                  v-model="filters.action"
                  :items="actionSelectOptions"
                  item-title="label"
                  item-value="action"
                  label="操作类型"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-select>
              </v-col>

              <v-col cols="12" md="2">
                <v-select
                  v-model="filters.resource"
                  :items="categorySelectOptions"
                  item-title="label"
                  item-value="value"
                  label="资源分类"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-select>
              </v-col>

              <v-col cols="12" md="2">
                <v-select
                  v-model="filters.status"
                  :items="statusOptions"
                  item-title="text"
                  item-value="value"
                  label="状态"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-select>
              </v-col>

              <v-col cols="12" md="3" class="d-flex gap-2 flex-wrap">
                <v-btn color="primary" :loading="loading" @click="applyFilters">
                  <v-icon start>mdi-magnify</v-icon>
                  查询
                </v-btn>
                <v-btn variant="outlined" @click="resetFilters">
                  <v-icon start>mdi-refresh</v-icon>
                  重置
                </v-btn>
                <v-btn color="success" variant="outlined" :loading="exporting" :disabled="loading" @click="exportLogs">
                  <v-icon start>mdi-download</v-icon>
                  导出
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :items-per-page-options="ITEMS_PER_PAGE_OPTIONS"
              :page="pagination.page"
              :headers="headers"
              :items="logs"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无审计日志"
              @update:options="onTableOptionsUpdate"
            >
              <!-- 操作列 -->
              <template #item.action="{ item }">
                <v-chip size="small" color="primary" variant="flat">
                  {{ formatAction(item.action) }}
                </v-chip>
              </template>

              <!-- 资源列 -->
              <template #item.resource="{ item }">
                <v-chip size="small" color="info" variant="flat">
                  {{ formatResource(item.resource) }}
                </v-chip>
              </template>

              <!-- 状态列 -->
              <template #item.status="{ item }">
                <v-chip :color="getStatusColor(item.status)" size="small">
                  {{ getStatusText(item.status) }}
                </v-chip>
              </template>

              <!-- IP地址列 -->
              <template #item.ip_address="{ item }">
                <span class="text-body-2">{{ item.ip_address || "-" }}</span>
              </template>

              <!-- 时间列 -->
              <template #item.created_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
              </template>

              <!-- 操作列 -->
              <template #item.actions="{ item }">
                <v-tooltip text="查看详情">
                  <template #activator="{ props }">
                    <v-btn icon="mdi-eye" size="small" variant="text" v-bind="props" @click="viewDetail(item)"></v-btn>
                  </template>
                </v-tooltip>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 详情对话框 -->
    <v-dialog v-model="detailDialog" max-width="800">
      <v-card v-if="selectedLog">
        <v-card-title class="text-h5 d-flex align-center">
          <v-icon start color="primary">mdi-file-document-outline</v-icon>
          日志详情
        </v-card-title>

        <v-divider></v-divider>

        <v-card-text class="pt-4">
          <v-row dense>
            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">日志ID</div>
              <div class="text-body-1">{{ selectedLog.id }}</div>
            </v-col>

            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">用户ID</div>
              <div class="text-body-1">{{ selectedLog.user_id }}</div>
            </v-col>

            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">操作类型</div>
              <div class="text-body-1">
                <v-chip size="small" color="primary">{{ formatAction(selectedLog.action) }}</v-chip>
              </div>
            </v-col>

            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">资源类型</div>
              <div class="text-body-1">
                <v-chip size="small" color="info">{{ formatResource(selectedLog.resource) }}</v-chip>
              </div>
            </v-col>

            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">状态</div>
              <div class="text-body-1">
                <v-chip :color="getStatusColor(selectedLog.status)" size="small">
                  {{ getStatusText(selectedLog.status) }}
                </v-chip>
              </div>
            </v-col>

            <v-col cols="6">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">IP地址</div>
              <div class="text-body-1">{{ selectedLog.ip_address || "-" }}</div>
            </v-col>

            <v-col cols="12">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">UserUserWithRolesDTO Agent</div>
              <div class="text-body-2" style="word-break: break-all">{{ selectedLog.user_agent || "-" }}</div>
            </v-col>

            <v-col v-if="selectedLog.details" cols="12">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">详细信息</div>
              <v-card variant="outlined" class="pa-3">
                <pre class="text-body-2" style="white-space: pre-wrap; word-break: break-word">{{ selectedLog.details }}</pre>
              </v-card>
            </v-col>

            <v-col cols="12">
              <div class="text-subtitle-2 text-medium-emphasis mb-1">创建时间</div>
              <div class="text-body-1">{{ formatDate(selectedLog.created_at) }}</div>
            </v-col>
          </v-row>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="detailDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.auditlogs-page {
  width: 100%;
}

.gap-2 {
  gap: 8px;
}
</style>
