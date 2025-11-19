<script setup lang="ts">
import { ref, watch } from "vue";
import { updateProfile, type UpdateProfileRequest } from "@/api/auth/user";

// Props
const props = defineProps<{
  user: any;
}>();

// Emits
const emit = defineEmits<{
  "update:success": [];
}>();

// 表单数据
const formData = ref<UpdateProfileRequest>({
  full_name: "",
  avatar: "",
  bio: "",
});

// 状态
const loading = ref(false);
const editMode = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

// 监听 user 变化，初始化表单数据
watch(
  () => props.user,
  (newUser) => {
    if (newUser) {
      formData.value = {
        full_name: newUser.full_name || "",
        avatar: newUser.avatar || "",
        bio: newUser.bio || "",
      };
    }
  },
  { immediate: true }
);

/**
 * 提交更新
 */
async function handleSubmit() {
  try {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    await updateProfile(formData.value);

    successMessage.value = "个人资料更新成功！";
    editMode.value = false;

    // 通知父组件刷新数据
    emit("update:success");
  } catch (error: any) {
    errorMessage.value = error.message || "更新失败";
  } finally {
    loading.value = false;
  }
}

/**
 * 取消编辑
 */
function cancelEdit() {
  editMode.value = false;
  // 重置表单数据
  if (props.user) {
    formData.value = {
      full_name: props.user.full_name || "",
      avatar: props.user.avatar || "",
      bio: props.user.bio || "",
    };
  }
  errorMessage.value = "";
  successMessage.value = "";
}
</script>

<template>
  <div>
    <!-- 查看模式 -->
    <div v-if="!editMode">
      <v-list density="compact">
        <v-list-item>
          <template #prepend>
            <v-icon>mdi-account-box</v-icon>
          </template>
          <v-list-item-title>姓名</v-list-item-title>
          <v-list-item-subtitle>{{ user.full_name || "未设置" }}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-image-account</v-icon>
          </template>
          <v-list-item-title>头像 URL</v-list-item-title>
          <v-list-item-subtitle>{{ user.avatar || "未设置" }}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-text</v-icon>
          </template>
          <v-list-item-title>个人简介</v-list-item-title>
          <v-list-item-subtitle>{{ user.bio || "未设置" }}</v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <v-btn color="primary" class="mt-4" prepend-icon="mdi-pencil" @click="editMode = true"> 编辑资料 </v-btn>
    </div>

    <!-- 编辑模式 -->
    <v-form v-else @submit.prevent="handleSubmit">
      <v-text-field v-model="formData.full_name" label="姓名" variant="outlined" placeholder="请输入您的姓名" class="mb-4" />

      <v-text-field v-model="formData.avatar" label="头像 URL" variant="outlined" placeholder="https://example.com/avatar.jpg" class="mb-4" />

      <v-textarea v-model="formData.bio" label="个人简介" variant="outlined" placeholder="介绍一下自己..." rows="3" class="mb-4" />

      <!-- 操作按钮 -->
      <div class="d-flex gap-2">
        <v-btn type="submit" color="primary" :loading="loading" prepend-icon="mdi-check"> 保存 </v-btn>
        <v-btn variant="outlined" @click="cancelEdit"> 取消 </v-btn>
      </div>
    </v-form>

    <!-- 成功/错误消息 -->
    <v-alert v-if="successMessage" type="success" density="compact" class="mt-4" closable @click:close="successMessage = ''">
      {{ successMessage }}
    </v-alert>

    <v-alert v-if="errorMessage" type="error" density="compact" class="mt-4" closable @click:close="errorMessage = ''">
      {{ errorMessage }}
    </v-alert>
  </div>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>
