<script setup lang="ts">
import { ref, onMounted } from "vue";
import { userProfileApi, extractData, type UserUserWithRolesDTO } from "@/api";
import BasicInfoForm from "./components/BasicInfoForm.vue";

// 用户信息 - using any for flexible user object from API response
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const user = ref<any>(null);
const loading = ref(false);
const errorMessage = ref("");

/**
 * 获取用户信息
 */
async function fetchUserInfo() {
  try {
    loading.value = true;
    errorMessage.value = "";
    const response = await userProfileApi.apiUserProfileGet();
    user.value = extractData<UserUserWithRolesDTO>(response.data);
  } catch (error) {
    console.error("获取用户信息失败:", error);
    errorMessage.value = (error as Error).message || "获取用户信息失败";
  } finally {
    loading.value = false;
  }
}

/**
 * 处理信息更新成功
 */
function handleUpdateSuccess() {
  fetchUserInfo();
}

// 组件挂载时获取用户信息
onMounted(() => {
  fetchUserInfo();
});
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">
          <v-icon class="mr-2">mdi-account-circle</v-icon>
          个人资料
        </h1>
      </v-col>
    </v-row>

    <!-- 加载中 -->
    <v-row v-if="loading">
      <v-col cols="12">
        <v-skeleton-loader type="card"></v-skeleton-loader>
      </v-col>
    </v-row>

    <!-- 错误消息 -->
    <v-row v-else-if="errorMessage">
      <v-col cols="12">
        <v-alert type="error" variant="tonal">
          {{ errorMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <!-- 用户信息 -->
    <v-row v-else-if="user">
      <v-col cols="12" md="8">
        <v-card>
          <v-card-title class="text-h6"> 基本信息 </v-card-title>
          <v-card-text>
            <BasicInfoForm :user="user" @update:success="handleUpdateSuccess" />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 侧边栏：用户统计信息 -->
      <v-col cols="12" md="4">
        <v-card class="mb-4">
          <v-card-title class="text-h6"> 账号信息 </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <template #prepend>
                  <v-icon>mdi-account</v-icon>
                </template>
                <v-list-item-title>用户名</v-list-item-title>
                <v-list-item-subtitle>{{ user.username }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template #prepend>
                  <v-icon>mdi-email</v-icon>
                </template>
                <v-list-item-title>邮箱</v-list-item-title>
                <v-list-item-subtitle>{{ user.email }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template #prepend>
                  <v-icon>mdi-calendar</v-icon>
                </template>
                <v-list-item-title>注册时间</v-list-item-title>
                <v-list-item-subtitle>{{ new Date(user.created_at).toLocaleDateString() }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="user.roles && user.roles.length > 0">
                <template #prepend>
                  <v-icon>mdi-shield-account</v-icon>
                </template>
                <v-list-item-title>角色</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip v-for="role in user.roles" :key="role.id" size="small" class="mr-1">
                    {{ role.name }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
/* 可根据需要添加自定义样式 */
</style>
