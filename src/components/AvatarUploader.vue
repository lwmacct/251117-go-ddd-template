<script setup lang="ts">
import { ref, computed, watch } from "vue";

/**
 * 头像上传组件
 * 支持文件选择、拖拽上传、预览和 Base64 转换
 */

// Props
const props = withDefaults(
  defineProps<{
    /** 当前头像 URL 或 Base64 */
    modelValue?: string;
    /** 头像大小（像素） */
    size?: number;
    /** 最大文件大小（MB） */
    maxSize?: number;
    /** 是否禁用 */
    disabled?: boolean;
  }>(),
  {
    modelValue: "",
    size: 120,
    maxSize: 2,
    disabled: false,
  },
);

// Emits
const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

// 状态
const isDragging = ref(false);
const errorMessage = ref("");
const fileInput = ref<HTMLInputElement | null>(null);

// 支持的图片格式
const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/gif", "image/webp"];
const ACCEPTED_EXTENSIONS = ".jpg,.jpeg,.png,.gif,.webp";

// 计算显示的头像
const displayAvatar = computed(() => props.modelValue || "");

// 是否有头像
const hasAvatar = computed(() => !!props.modelValue);

/**
 * 验证文件
 */
function validateFile(file: File): boolean {
  errorMessage.value = "";

  // 检查文件类型
  if (!ACCEPTED_TYPES.includes(file.type)) {
    errorMessage.value = "仅支持 JPG、PNG、GIF、WebP 格式的图片";
    return false;
  }

  // 检查文件大小
  const maxBytes = props.maxSize * 1024 * 1024;
  if (file.size > maxBytes) {
    errorMessage.value = `图片大小不能超过 ${props.maxSize}MB`;
    return false;
  }

  return true;
}

/**
 * 读取文件为 Base64
 */
function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = () => reject(new Error("读取文件失败"));
    reader.readAsDataURL(file);
  });
}

/**
 * 处理文件
 */
async function handleFile(file: File) {
  if (!validateFile(file)) {
    return;
  }

  try {
    const base64 = await readFileAsBase64(file);
    emit("update:modelValue", base64);
  } catch {
    errorMessage.value = "读取图片失败，请重试";
  }
}

/**
 * 点击上传
 */
function handleClick() {
  if (props.disabled) return;
  fileInput.value?.click();
}

/**
 * 文件选择
 */
function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    handleFile(file);
  }
  // 重置 input，允许选择相同文件
  input.value = "";
}

/**
 * 拖拽进入
 */
function handleDragEnter(event: DragEvent) {
  event.preventDefault();
  if (!props.disabled) {
    isDragging.value = true;
  }
}

/**
 * 拖拽离开
 */
function handleDragLeave(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;
}

/**
 * 拖拽悬停
 */
function handleDragOver(event: DragEvent) {
  event.preventDefault();
}

/**
 * 拖拽放下
 */
function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (props.disabled) return;

  const file = event.dataTransfer?.files[0];
  if (file) {
    handleFile(file);
  }
}

/**
 * 移除头像
 */
function handleRemove() {
  emit("update:modelValue", "");
  errorMessage.value = "";
}

// 清除错误消息
watch(
  () => props.modelValue,
  () => {
    errorMessage.value = "";
  },
);
</script>

<template>
  <div class="avatar-uploader">
    <!-- 上传区域 -->
    <div
      class="upload-area"
      :class="{
        'is-dragging': isDragging,
        'has-avatar': hasAvatar,
        'is-disabled': disabled,
      }"
      :style="{ width: `${size}px`, height: `${size}px` }"
      @click="handleClick"
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @dragover="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 已有头像 -->
      <template v-if="hasAvatar">
        <v-avatar :size="size - 8" class="avatar-preview">
          <v-img :src="displayAvatar" cover />
        </v-avatar>
        <div class="overlay">
          <v-icon color="white">mdi-camera</v-icon>
        </div>
      </template>

      <!-- 无头像 -->
      <template v-else>
        <div class="placeholder">
          <v-icon size="32" color="grey">mdi-cloud-upload</v-icon>
          <span class="text-caption text-grey mt-1">点击或拖拽上传</span>
        </div>
      </template>
    </div>

    <!-- 隐藏的文件输入 -->
    <input ref="fileInput" type="file" :accept="ACCEPTED_EXTENSIONS" class="d-none" @change="handleFileChange" />

    <!-- 操作按钮 -->
    <div v-if="hasAvatar && !disabled" class="actions mt-2">
      <v-btn size="x-small" variant="text" color="error" @click.stop="handleRemove">
        <v-icon size="small">mdi-delete</v-icon>
        移除
      </v-btn>
    </div>

    <!-- 错误消息 -->
    <v-alert v-if="errorMessage" type="error" density="compact" variant="tonal" class="mt-2" style="max-width: 200px">
      {{ errorMessage }}
    </v-alert>

    <!-- 提示信息 -->
    <div class="text-caption text-grey mt-1" style="max-width: 200px">支持 JPG、PNG、GIF，最大 {{ maxSize }}MB</div>
  </div>
</template>

<style scoped>
.avatar-uploader {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
}

.upload-area {
  position: relative;
  border: 2px dashed #ccc;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
}

.upload-area:hover:not(.is-disabled) {
  border-color: #1976d2;
}

.upload-area.is-dragging {
  border-color: #1976d2;
  background-color: rgba(25, 118, 210, 0.1);
}

.upload-area.has-avatar {
  border-style: solid;
  border-color: transparent;
}

.upload-area.is-disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.avatar-preview {
  transition: transform 0.3s ease;
}

.overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s ease;
  border-radius: 50%;
}

.upload-area:hover:not(.is-disabled) .overlay {
  opacity: 1;
}

.actions {
  display: flex;
  gap: 4px;
}
</style>
