<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useTokens } from "./composables/useTokens";
import TokenDialog from "./components/TokenDialog.vue";
import TokenDisplay from "./components/TokenDisplay.vue";
import type { PatCreateTokenDTO } from "@models";

const {
  tokens,
  loading,
  errorMessage,
  successMessage,
  fetchTokens,
  createToken,
  deleteToken,
  disableToken,
  enableToken,
  clearMessages,
} = useTokens();

const tokenDialog = ref(false);
const tokenDisplayDialog = ref(false);
const deleteDialog = ref(false);
const newToken = ref("");
const newTokenName = ref("");
const tokenToOperate = ref<number | null>(null);

const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "Token 名称", key: "name", sortable: true },
  { title: "Token 前缀", key: "token_prefix" },
  { title: "状态", key: "status" },
  { title: "最后使用", key: "last_used_at" },
  { title: "过期时间", key: "expires_at" },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

onMounted(() => {
  fetchTokens();
});

const openCreateDialog = () => {
  tokenDialog.value = true;
};

const handleCreateToken = async (data: PatCreateTokenDTO) => {
  const response = await createToken(data);
  if (response) {
    newToken.value = response.plain_token ?? "";
    newTokenName.value = response.token?.name ?? "";
    tokenDisplayDialog.value = true;
  }
};

const openDeleteDialog = (tokenId?: number) => {
  if (tokenId === undefined) return;
  tokenToOperate.value = tokenId;
  deleteDialog.value = true;
};

const confirmDelete = async () => {
  if (tokenToOperate.value === null) return;

  const success = await deleteToken(tokenToOperate.value);
  if (success) {
    deleteDialog.value = false;
    tokenToOperate.value = null;
  }
};

const handleDisable = async (tokenId?: number) => {
  if (tokenId === undefined) return;
  await disableToken(tokenId);
};

const handleEnable = async (tokenId?: number) => {
  if (tokenId === undefined) return;
  await enableToken(tokenId);
};

const formatDate = (dateString?: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const getStatusColor = (status?: string) => {
  const colors: Record<string, string> = {
    active: "success",
    disabled: "warning",
    expired: "warning",
  };
  return colors[status ?? ""] || "default";
};

const getStatusText = (status?: string) => {
  const texts: Record<string, string> = {
    active: "正常",
    disabled: "已禁用",
    expired: "已过期",
  };
  return texts[status ?? ""] || (status ?? "-");
};

const isTokenExpired = (expiresAt?: string) => {
  if (!expiresAt) return false;
  return new Date(expiresAt) < new Date();
};
</script>

<template>
  <div class="tokens-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">Personal Access Tokens</h1>
        <p class="text-body-2 text-medium-emphasis">Personal Access Tokens 可用于通过 API 访问系统资源</p>
      </v-col>
    </v-row>

    <v-row v-if="errorMessage || successMessage">
      <v-col cols="12">
        <v-alert v-if="errorMessage" type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
        </v-alert>
        <v-alert v-if="successMessage" type="success" closable @click:close="clearMessages">
          {{ successMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <span class="text-h6">我的 Tokens</span>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  创建 Token
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table
              :headers="headers"
              :items="tokens"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无 Token"
            >
              <template #item.token_prefix="{ item }">
                <code>{{ item.token_prefix }}...</code>
              </template>

              <template #item.status="{ item }">
                <v-chip :color="getStatusColor(item.status)" size="small">
                  {{ getStatusText(item.status) }}
                </v-chip>
                <v-chip v-if="isTokenExpired(item.expires_at)" color="warning" size="small" class="ml-1"> 已过期 </v-chip>
              </template>

              <template #item.last_used_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.last_used_at) }}</span>
              </template>

              <template #item.expires_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.expires_at) || "永不过期" }}</span>
              </template>

              <template #item.created_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
              </template>

              <template #item.actions="{ item }">
                <div class="d-flex align-center ga-2">
                  <v-tooltip v-if="item.status === 'active'" text="禁用">
                    <template #activator="{ props }">
                      <v-btn
                        icon="mdi-pause-circle"
                        size="small"
                        variant="text"
                        color="warning"
                        v-bind="props"
                        @click="handleDisable(item.id)"
                      ></v-btn>
                    </template>
                  </v-tooltip>
                  <v-tooltip v-else-if="item.status === 'disabled'" text="启用">
                    <template #activator="{ props }">
                      <v-btn
                        icon="mdi-play-circle"
                        size="small"
                        variant="text"
                        color="success"
                        v-bind="props"
                        @click="handleEnable(item.id)"
                      ></v-btn>
                    </template>
                  </v-tooltip>
                  <v-tooltip text="删除">
                    <template #activator="{ props }">
                      <v-btn
                        icon="mdi-delete"
                        size="small"
                        variant="text"
                        color="error"
                        v-bind="props"
                        @click="openDeleteDialog(item.id)"
                      ></v-btn>
                    </template>
                  </v-tooltip>
                </div>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <TokenDialog v-model="tokenDialog" @save="handleCreateToken" />

    <TokenDisplay v-model="tokenDisplayDialog" :token="newToken" :token-name="newTokenName" />

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text> 确定要删除此 Token 吗？删除后将无法恢复，使用此 Token 的 API 请求将失败。 </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="elevated" :loading="loading" @click="confirmDelete">删除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.tokens-page {
  width: 100%;
}
</style>
