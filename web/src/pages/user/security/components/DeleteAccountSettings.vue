<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { deleteAccount } from "@/api/auth/user";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();

// 表单数据
const password = ref("");
const confirmText = ref("");

// 状态
const loading = ref(false);
const showPassword = ref(false);
const showConfirmDialog = ref(false);

// 消息
const errorMessage = ref("");

// 确认文本
const CONFIRM_TEXT = "DELETE";

/**
 * 打开确认对话框
 */
function openConfirmDialog() {
  if (!password.value) {
    errorMessage.value = "请输入当前密码";
    return;
  }
  errorMessage.value = "";
  showConfirmDialog.value = true;
}

/**
 * 关闭确认对话框
 */
function closeConfirmDialog() {
  showConfirmDialog.value = false;
  confirmText.value = "";
}

/**
 * 执行删除账户
 */
async function handleDelete() {
  if (confirmText.value !== CONFIRM_TEXT) {
    errorMessage.value = `请输入 "${CONFIRM_TEXT}" 确认删除`;
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = "";

    await deleteAccount({ password: password.value });

    // 删除成功，清除认证状态并跳转
    closeConfirmDialog();
    await authStore.logout();
    router.push("/auth/login");
  } catch (error) {
    errorMessage.value = (error as Error).message || "删除账户失败";
    closeConfirmDialog();
  } finally {
    loading.value = false;
  }
}

/**
 * 重置表单
 */
function resetForm() {
  password.value = "";
  confirmText.value = "";
  errorMessage.value = "";
}
</script>

<template>
  <div>
    <!-- 警告说明 -->
    <v-alert type="error" variant="tonal" density="compact" class="mb-4">
      <div class="text-body-2">
        <div class="font-weight-bold mb-2">
          <v-icon size="small" class="mr-1">mdi-alert</v-icon>
          危险操作警告
        </div>
        <ul class="ml-4">
          <li>删除账户后，您的所有数据将被永久删除</li>
          <li>此操作不可撤销，请谨慎操作</li>
          <li>删除后您将无法使用此账户登录</li>
        </ul>
      </div>
    </v-alert>

    <!-- 删除账户表单 -->
    <v-form @submit.prevent="openConfirmDialog">
      <!-- 当前密码 -->
      <v-text-field
        v-model="password"
        label="输入当前密码以确认身份"
        :type="showPassword ? 'text' : 'password'"
        :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
        variant="outlined"
        required
        class="mb-4"
        @click:append-inner="showPassword = !showPassword"
      />

      <!-- 操作按钮 -->
      <div class="d-flex gap-2">
        <v-btn type="submit" color="error" :loading="loading" prepend-icon="mdi-delete-forever"> 删除我的账户 </v-btn>
        <v-btn variant="outlined" @click="resetForm"> 重置 </v-btn>
      </div>
    </v-form>

    <!-- 错误消息 -->
    <v-alert v-if="errorMessage" type="error" density="compact" class="mt-4" closable @click:close="errorMessage = ''">
      {{ errorMessage }}
    </v-alert>

    <!-- 确认对话框 -->
    <v-dialog v-model="showConfirmDialog" max-width="500" persistent>
      <v-card>
        <v-card-title class="bg-error text-white">
          <v-icon class="mr-2">mdi-alert-circle</v-icon>
          确认删除账户
        </v-card-title>

        <v-card-text class="pt-4">
          <v-alert type="warning" variant="tonal" density="compact" class="mb-4">
            您即将永久删除您的账户。此操作无法撤销。
          </v-alert>

          <p class="mb-4">
            请输入 <strong class="text-error">{{ CONFIRM_TEXT }}</strong> 确认删除：
          </p>

          <v-text-field
            v-model="confirmText"
            label="确认文本"
            variant="outlined"
            :placeholder="`请输入 ${CONFIRM_TEXT}`"
            autofocus
          />
        </v-card-text>

        <v-card-actions>
          <v-spacer />
          <v-btn variant="outlined" @click="closeConfirmDialog"> 取消 </v-btn>
          <v-btn color="error" :loading="loading" :disabled="confirmText !== CONFIRM_TEXT" @click="handleDelete">
            确认删除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>
