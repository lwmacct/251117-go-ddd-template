<script setup lang="ts">
import { computed } from "vue";
import { checkPasswordStrength } from "@/utils/auth/validation";

/**
 * 密码强度指示器组件
 * 显示密码强度的进度条和文字提示
 */

interface Props {
  password: string;
  showHints?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showHints: true,
});

// 计算密码强度
const strength = computed(() => {
  if (!props.password) return null;
  return checkPasswordStrength(props.password);
});

// 强度配置
const strengthConfig = computed(() => {
  switch (strength.value) {
    case "weak":
      return {
        color: "error",
        text: "弱",
        value: 33,
        icon: "mdi-shield-alert",
      };
    case "medium":
      return {
        color: "warning",
        text: "中",
        value: 66,
        icon: "mdi-shield-half-full",
      };
    case "strong":
      return {
        color: "success",
        text: "强",
        value: 100,
        icon: "mdi-shield-check",
      };
    default:
      return null;
  }
});

// 密码要求检查
const requirements = computed(() => {
  const pwd = props.password || "";
  return [
    { text: "至少 8 个字符", met: pwd.length >= 8 },
    { text: "包含小写字母", met: /[a-z]/.test(pwd) },
    { text: "包含大写字母", met: /[A-Z]/.test(pwd) },
    { text: "包含数字", met: /\d/.test(pwd) },
    { text: "包含特殊字符", met: /[^a-zA-Z0-9]/.test(pwd) },
  ];
});
</script>

<template>
  <div v-if="password" class="password-strength">
    <!-- 强度条 -->
    <div class="d-flex align-center mb-1">
      <v-progress-linear
        :model-value="strengthConfig?.value || 0"
        :color="strengthConfig?.color || 'grey'"
        height="6"
        rounded
        class="flex-grow-1"
      />
      <v-chip v-if="strengthConfig" :color="strengthConfig.color" size="x-small" variant="tonal" class="ml-2">
        <v-icon size="14" class="mr-1">{{ strengthConfig.icon }}</v-icon>
        {{ strengthConfig.text }}
      </v-chip>
    </div>

    <!-- 密码要求提示 -->
    <div v-if="showHints" class="requirements mt-2">
      <div
        v-for="(req, index) in requirements"
        :key="index"
        class="d-flex align-center text-body-2"
        :class="req.met ? 'text-success' : 'text-grey'"
      >
        <v-icon size="14" class="mr-1">
          {{ req.met ? "mdi-check-circle" : "mdi-circle-outline" }}
        </v-icon>
        {{ req.text }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.password-strength {
  margin-top: 4px;
}

.requirements {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4px;
}

@media (max-width: 600px) {
  .requirements {
    grid-template-columns: 1fr;
  }
}
</style>
