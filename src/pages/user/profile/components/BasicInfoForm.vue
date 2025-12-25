<script setup lang="ts">
import { ref, watch } from "vue";
import { userProfileApi, type HandlerUpdateProfileRequest } from "@/api";
import AvatarUploader from "@/components/AvatarUploader.vue";

// Props - using any for flexible user object from API response

const props = defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  user: any;
}>();

// Emits
const emit = defineEmits<{
  "update:success": [];
}>();

// 表单数据
const formData = ref<HandlerUpdateProfileRequest>({
  full_name: "",
  avatar: "",
  bio: "",
});

// 状态
const loading = ref(false);
const dialogVisible = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function resetFormFromUser(user?: any) {
  const source = user ?? props.user;
  formData.value = {
    full_name: source?.full_name || "",
    avatar: source?.avatar || "",
    bio: source?.bio || "",
  };
}

// 监听 user 变化，初始化表单数据
watch(
  () => props.user,
  (newUser) => {
    resetFormFromUser(newUser);
  },
  { immediate: true },
);

/**
 * 提交更新
 */
async function handleSubmit() {
  try {
    loading.value = true;
    errorMessage.value = "";

    await userProfileApi.apiUserProfilePut(formData.value);

    successMessage.value = "个人资料更新成功！";
    dialogVisible.value = false;

    // 通知父组件刷新数据
    emit("update:success");
  } catch (error) {
    errorMessage.value = (error as Error).message || "更新失败";
  } finally {
    loading.value = false;
  }
}

function openDialog() {
  resetFormFromUser();
  successMessage.value = "";
  dialogVisible.value = true;
}

function handleCloseDialog() {
  dialogVisible.value = false;
}

watch(dialogVisible, (isOpen, wasOpen) => {
  if (!isOpen && wasOpen) {
    resetFormFromUser();
    errorMessage.value = "";
  }
});
</script>

<template>
  <div>
    <!-- 头像显示区 -->
    <div class="d-flex align-center mb-4">
      <v-avatar :size="80" class="mr-4">
        <v-img v-if="user.avatar" :src="user.avatar" cover />
        <v-icon v-else size="48" color="grey">mdi-account-circle</v-icon>
      </v-avatar>
      <div>
        <div class="text-h6">{{ user.full_name || user.username }}</div>
        <div class="text-body-2 text-grey">{{ user.bio || "这个人很懒，什么都没写" }}</div>
      </div>
    </div>

    <v-divider class="mb-4" />

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
          <v-icon>mdi-text</v-icon>
        </template>
        <v-list-item-title>个人简介</v-list-item-title>
        <v-list-item-subtitle>{{ user.bio || "未设置" }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>

    <div class="d-flex align-center gap-2 mt-4">
      <v-btn color="primary" prepend-icon="mdi-pencil" @click="openDialog"> 编辑资料 </v-btn>
    </div>

    <v-alert v-if="successMessage" type="success" density="compact" class="mt-4" closable @click:close="successMessage = ''">
      {{ successMessage }}
    </v-alert>

    <v-dialog v-model="dialogVisible" max-width="600">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon class="mr-2">mdi-account-edit</v-icon>
          编辑个人资料
        </v-card-title>
        <v-card-text>
          <v-alert v-if="errorMessage" type="error" density="compact" class="mb-4" closable @click:close="errorMessage = ''">
            {{ errorMessage }}
          </v-alert>
          <v-form @submit.prevent="handleSubmit">
            <!-- 头像上传 -->
            <div class="d-flex justify-center mb-6">
              <AvatarUploader v-model="formData.avatar" :size="120" :max-size="2" />
            </div>

            <v-text-field
              v-model="formData.full_name"
              label="姓名"
              variant="outlined"
              placeholder="请输入您的姓名"
              class="mb-4"
            />

            <v-textarea
              v-model="formData.bio"
              label="个人简介"
              variant="outlined"
              placeholder="介绍一下自己..."
              rows="3"
              class="mb-4"
            />

            <v-divider class="my-4" />

            <div class="d-flex justify-end gap-2 mt-2">
              <v-btn type="submit" color="primary" :loading="loading" prepend-icon="mdi-check"> 保存 </v-btn>
              <v-btn variant="text" prepend-icon="mdi-close" @click="handleCloseDialog"> 取消 </v-btn>
            </div>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>
