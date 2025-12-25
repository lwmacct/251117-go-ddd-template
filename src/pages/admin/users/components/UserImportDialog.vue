<script setup lang="ts">
import { ref, computed } from "vue";
import { parseUserCSV, readFileAsText, generateUserCSVTemplate, type ParsedUser, type ParseError } from "@/utils/import";
import { adminUserApi, type UserBatchCreateUserResultDTO } from "@/api";

/**
 * 用户批量导入对话框
 */

interface Props {
  modelValue: boolean;
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "imported", result: { success: number; failed: number }): void;
}

const _props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 状态
const step = ref<"upload" | "preview" | "result">("upload");
const loading = ref(false);
const errorMessage = ref("");

// 文件上传
const selectedFile = ref<File | null>(null);
const isDragging = ref(false);

// 解析结果
const parsedData = ref<ParsedUser[]>([]);
const parseErrors = ref<ParseError[]>([]);
const totalRows = ref(0);
const validRows = ref(0);

// 导入结果
const importResult = ref<{
  success: number;
  failed: number;
  errors: Array<{ index: number; username: string; error: string }>;
} | null>(null);

// 预览表格 headers
const previewHeaders = [
  { title: "用户名", key: "username" },
  { title: "邮箱", key: "email" },
  { title: "密码", key: "password" },
  { title: "姓名", key: "full_name" },
  { title: "状态", key: "status" },
];

// 计算属性
const canProceed = computed(() => parsedData.value.length > 0 && parseErrors.value.length === 0);

const hasParseErrors = computed(() => parseErrors.value.length > 0);

// 关闭对话框
const closeDialog = () => {
  emit("update:modelValue", false);
  resetState();
};

// 重置状态
const resetState = () => {
  step.value = "upload";
  selectedFile.value = null;
  parsedData.value = [];
  parseErrors.value = [];
  totalRows.value = 0;
  validRows.value = 0;
  importResult.value = null;
  errorMessage.value = "";
  loading.value = false;
};

// 处理文件选择
const handleFileSelect = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    await processFile(file);
  }
};

// 处理拖拽
const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = true;
};

const handleDragLeave = () => {
  isDragging.value = false;
};

const handleDrop = async (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;
  const file = event.dataTransfer?.files?.[0];
  if (file) {
    await processFile(file);
  }
};

// 处理文件
const processFile = async (file: File) => {
  if (!file.name.endsWith(".csv")) {
    errorMessage.value = "请选择 CSV 文件";
    return;
  }

  selectedFile.value = file;
  errorMessage.value = "";
  loading.value = true;

  try {
    const content = await readFileAsText(file);
    const result = parseUserCSV(content);

    parsedData.value = result.data;
    parseErrors.value = result.errors;
    totalRows.value = result.totalRows;
    validRows.value = result.validRows;

    step.value = "preview";
  } catch (error) {
    errorMessage.value = (error as Error).message || "文件解析失败";
  } finally {
    loading.value = false;
  }
};

// 下载模板
const downloadTemplate = () => {
  const content = generateUserCSVTemplate();
  const blob = new Blob([content], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = "user_import_template.csv";
  link.click();
  URL.revokeObjectURL(url);
};

// 执行导入
const handleImport = async () => {
  loading.value = true;
  errorMessage.value = "";

  try {
    const users = parsedData.value.map((user) => ({
      username: user.username,
      email: user.email,
      password: user.password,
      full_name: user.full_name,
      status: user.status as "active" | "inactive" | undefined,
    }));

    const response = await adminUserApi.apiAdminUsersBatchPost({ users });
    const result = response.data.data as UserBatchCreateUserResultDTO;
    importResult.value = {
      success: result.success ?? 0,
      failed: result.failed ?? 0,
      errors: (result.errors ?? []).map((e) => ({
        index: e.index ?? 0,
        username: e.username ?? "",
        error: e.error ?? "",
      })),
    };
    step.value = "result";

    // 通知父组件
    emit("imported", { success: importResult.value.success, failed: importResult.value.failed });
  } catch (error) {
    errorMessage.value = (error as Error).message || "导入失败";
  } finally {
    loading.value = false;
  }
};

// 返回上传步骤
const backToUpload = () => {
  step.value = "upload";
  selectedFile.value = null;
  parsedData.value = [];
  parseErrors.value = [];
  errorMessage.value = "";
};
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="900" persistent @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2">mdi-upload</v-icon>
        <span>批量导入用户</span>
        <v-spacer />
        <v-btn icon variant="text" @click="closeDialog">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-divider />

      <!-- 步骤指示器 -->
      <v-stepper :model-value="step === 'upload' ? 1 : step === 'preview' ? 2 : 3" alt-labels class="elevation-0">
        <v-stepper-header>
          <v-stepper-item title="上传文件" value="1" :complete="step !== 'upload'" />
          <v-divider />
          <v-stepper-item title="预览数据" value="2" :complete="step === 'result'" />
          <v-divider />
          <v-stepper-item title="导入结果" value="3" />
        </v-stepper-header>
      </v-stepper>

      <v-card-text>
        <!-- 错误提示 -->
        <v-alert v-if="errorMessage" type="error" density="compact" class="mb-4" closable @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>

        <!-- 步骤 1: 上传文件 -->
        <div v-if="step === 'upload'">
          <!-- 拖拽上传区域 -->
          <div
            class="upload-zone pa-8 text-center rounded-lg mb-4"
            :class="{ dragging: isDragging }"
            @dragover="handleDragOver"
            @dragleave="handleDragLeave"
            @drop="handleDrop"
          >
            <v-icon size="64" color="grey">mdi-file-upload-outline</v-icon>
            <div class="text-h6 mt-4">拖拽 CSV 文件到此处</div>
            <div class="text-body-2 text-grey mt-2">或点击下方按钮选择文件</div>

            <input ref="fileInput" type="file" accept=".csv" hidden @change="handleFileSelect" />

            <div class="mt-4">
              <v-btn
                color="primary"
                variant="elevated"
                :loading="loading"
                @click="($refs.fileInput as HTMLInputElement).click()"
              >
                <v-icon left>mdi-folder-open</v-icon>
                选择文件
              </v-btn>
            </div>
          </div>

          <!-- 模板下载 -->
          <v-card variant="outlined" class="pa-4">
            <div class="d-flex align-center">
              <v-icon color="info" class="mr-3">mdi-information-outline</v-icon>
              <div class="flex-grow-1">
                <div class="text-subtitle-2">CSV 格式要求</div>
                <div class="text-body-2 text-grey">
                  必需列: username, email, password<br />
                  可选列: full_name, status
                </div>
              </div>
              <v-btn variant="text" color="primary" @click="downloadTemplate">
                <v-icon left>mdi-download</v-icon>
                下载模板
              </v-btn>
            </div>
          </v-card>
        </div>

        <!-- 步骤 2: 预览数据 -->
        <div v-else-if="step === 'preview'">
          <!-- 统计信息 -->
          <v-row class="mb-4">
            <v-col cols="4">
              <v-card variant="outlined" class="pa-3 text-center">
                <div class="text-h4">{{ totalRows }}</div>
                <div class="text-body-2 text-grey">总记录数</div>
              </v-card>
            </v-col>
            <v-col cols="4">
              <v-card variant="outlined" class="pa-3 text-center">
                <div class="text-h4 text-success">{{ validRows }}</div>
                <div class="text-body-2 text-grey">有效记录</div>
              </v-card>
            </v-col>
            <v-col cols="4">
              <v-card variant="outlined" class="pa-3 text-center">
                <div class="text-h4" :class="parseErrors.length > 0 ? 'text-error' : ''">{{ parseErrors.length }}</div>
                <div class="text-body-2 text-grey">验证错误</div>
              </v-card>
            </v-col>
          </v-row>

          <!-- 解析错误列表 -->
          <v-alert v-if="hasParseErrors" type="warning" variant="tonal" class="mb-4">
            <div class="text-subtitle-2 mb-2">数据验证错误（以下记录将被跳过）：</div>
            <div v-for="(err, index) in parseErrors.slice(0, 10)" :key="index" class="text-body-2">
              第 {{ err.row }} 行 - {{ err.field }}: {{ err.message }}
            </div>
            <div v-if="parseErrors.length > 10" class="text-body-2 mt-2">... 还有 {{ parseErrors.length - 10 }} 个错误</div>
          </v-alert>

          <!-- 数据预览表格 -->
          <v-data-table
            :headers="previewHeaders"
            :items="parsedData"
            :items-per-page="10"
            density="compact"
            class="elevation-1"
          >
            <template #item.password>
              <span class="text-grey">••••••••</span>
            </template>
            <template #item.status="{ item }">
              <v-chip :color="item.status === 'active' ? 'success' : 'grey'" size="small">
                {{ item.status === "active" ? "启用" : "禁用" }}
              </v-chip>
            </template>
            <template #item.full_name="{ item }">
              {{ item.full_name || "-" }}
            </template>
          </v-data-table>
        </div>

        <!-- 步骤 3: 导入结果 -->
        <div v-else-if="step === 'result'">
          <div class="text-center pa-8">
            <v-icon :color="importResult?.failed === 0 ? 'success' : 'warning'" size="80">
              {{ importResult?.failed === 0 ? "mdi-check-circle" : "mdi-alert-circle" }}
            </v-icon>

            <div class="text-h5 mt-4">导入完成</div>

            <v-row class="mt-6 justify-center">
              <v-col cols="4">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4 text-success">{{ importResult?.success || 0 }}</div>
                  <div class="text-body-2 text-grey">成功导入</div>
                </v-card>
              </v-col>
              <v-col cols="4">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4" :class="(importResult?.failed || 0) > 0 ? 'text-error' : ''">
                    {{ importResult?.failed || 0 }}
                  </div>
                  <div class="text-body-2 text-grey">导入失败</div>
                </v-card>
              </v-col>
            </v-row>

            <!-- 失败详情 -->
            <v-alert
              v-if="importResult?.errors && importResult.errors.length > 0"
              type="error"
              variant="tonal"
              class="mt-6 text-left"
            >
              <div class="text-subtitle-2 mb-2">失败详情：</div>
              <div v-for="(err, index) in importResult.errors.slice(0, 10)" :key="index" class="text-body-2">
                用户 "{{ err.username }}": {{ err.error }}
              </div>
              <div v-if="importResult.errors.length > 10" class="text-body-2 mt-2">
                ... 还有 {{ importResult.errors.length - 10 }} 个错误
              </div>
            </v-alert>
          </div>
        </div>
      </v-card-text>

      <v-divider />

      <v-card-actions class="pa-4">
        <v-btn v-if="step === 'preview'" variant="text" @click="backToUpload">
          <v-icon left>mdi-arrow-left</v-icon>
          返回
        </v-btn>
        <v-spacer />
        <v-btn variant="text" @click="closeDialog">
          {{ step === "result" ? "关闭" : "取消" }}
        </v-btn>
        <v-btn
          v-if="step === 'preview'"
          color="primary"
          variant="elevated"
          :disabled="!canProceed"
          :loading="loading"
          @click="handleImport"
        >
          <v-icon left>mdi-upload</v-icon>
          开始导入 ({{ validRows }} 条)
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.upload-zone {
  border: 2px dashed #e0e0e0;
  background-color: #fafafa;
  transition: all 0.3s ease;
  cursor: pointer;
}

.upload-zone:hover,
.upload-zone.dragging {
  border-color: #1976d2;
  background-color: #e3f2fd;
}
</style>
