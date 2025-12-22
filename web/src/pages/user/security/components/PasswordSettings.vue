<script setup lang="ts">
import { ref } from "vue";
import { changePassword } from "@/api/auth/user";
import PasswordStrengthIndicator from "@/components/PasswordStrengthIndicator.vue";

// 表单数据
const oldPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");

// 状态
const loading = ref(false);
const showOldPassword = ref(false);
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// 消息
const errorMessage = ref("");
const successMessage = ref("");

/**
 * 密码验证规则
 */
const passwordRules = [
  (v: string) => !!v || "请输入新密码",
  (v: string) => v.length >= 8 || "密码至少 8 个字符",
  (v: string) => /[a-z]/.test(v) || "密码必须包含小写字母",
  (v: string) => /[A-Z]/.test(v) || "密码必须包含大写字母",
  (v: string) => /\d/.test(v) || "密码必须包含数字",
];

const confirmPasswordRules = [
  (v: string) => !!v || "请确认新密码",
  (v: string) => v === newPassword.value || "两次输入的密码不一致",
];

/**
 * 提交修改密码
 */
async function handleSubmit() {
  // 验证表单
  if (!oldPassword.value) {
    errorMessage.value = "请输入当前密码";
    return;
  }

  if (!newPassword.value) {
    errorMessage.value = "请输入新密码";
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = "两次输入的密码不一致";
    return;
  }

  if (newPassword.value.length < 8) {
    errorMessage.value = "新密码至少 8 个字符";
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    await changePassword({
      old_password: oldPassword.value,
      new_password: newPassword.value,
    });

    successMessage.value = "密码修改成功！";

    // 清空表单
    oldPassword.value = "";
    newPassword.value = "";
    confirmPassword.value = "";
  } catch (error) {
    errorMessage.value = (error as Error).message || "密码修改失败";
  } finally {
    loading.value = false;
  }
}

/**
 * 重置表单
 */
function resetForm() {
  oldPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
  errorMessage.value = "";
  successMessage.value = "";
}
</script>

<template>
  <div>
    <!-- 说明 -->
    <v-alert type="info" variant="tonal" density="compact" class="mb-4">
      <div class="text-body-2">
        <div class="font-weight-bold mb-2">密码安全要求</div>
        <ul class="ml-4">
          <li>密码长度至少 8 个字符，建议 12 个字符以上</li>
          <li>必须包含大写字母、小写字母和数字</li>
          <li>建议包含特殊字符（如 !@#$%^&*）</li>
        </ul>
      </div>
    </v-alert>

    <!-- 密码修改表单 -->
    <v-form @submit.prevent="handleSubmit">
      <!-- 当前密码 -->
      <v-text-field
        v-model="oldPassword"
        label="当前密码"
        :type="showOldPassword ? 'text' : 'password'"
        :append-inner-icon="showOldPassword ? 'mdi-eye-off' : 'mdi-eye'"
        variant="outlined"
        required
        class="mb-4"
        @click:append-inner="showOldPassword = !showOldPassword"
      />

      <!-- 新密码 -->
      <v-text-field
        v-model="newPassword"
        label="新密码"
        :type="showNewPassword ? 'text' : 'password'"
        :append-inner-icon="showNewPassword ? 'mdi-eye-off' : 'mdi-eye'"
        :rules="passwordRules"
        variant="outlined"
        required
        class="mb-2"
        @click:append-inner="showNewPassword = !showNewPassword"
      />

      <!-- 密码强度指示器 -->
      <PasswordStrengthIndicator :password="newPassword" class="mb-4" />

      <!-- 确认新密码 -->
      <v-text-field
        v-model="confirmPassword"
        label="确认新密码"
        :type="showConfirmPassword ? 'text' : 'password'"
        :append-inner-icon="showConfirmPassword ? 'mdi-eye-off' : 'mdi-eye'"
        :rules="confirmPasswordRules"
        variant="outlined"
        required
        class="mb-4"
        @click:append-inner="showConfirmPassword = !showConfirmPassword"
      />

      <!-- 操作按钮 -->
      <div class="d-flex gap-2">
        <v-btn type="submit" color="primary" :loading="loading" prepend-icon="mdi-check"> 修改密码 </v-btn>
        <v-btn variant="outlined" @click="resetForm"> 重置 </v-btn>
      </div>
    </v-form>

    <!-- 成功/错误消息 -->
    <v-alert
      v-if="successMessage"
      type="success"
      density="compact"
      class="mt-4"
      closable
      @click:close="successMessage = ''"
    >
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
