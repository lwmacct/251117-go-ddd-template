<script setup lang="ts">
import { ref } from "vue";

interface Props {
  modelValue: boolean;
  token: string;
  tokenName: string;
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const copied = ref(false);

const copyToken = async () => {
  try {
    await navigator.clipboard.writeText(props.token);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (error) {
    console.error("Failed to copy token:", error);
  }
};

const closeDialog = () => {
  emit("update:modelValue", false);
};
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="700"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="bg-primary">
        <v-icon start>mdi-key-variant</v-icon>
        Token 创建成功
      </v-card-title>

      <v-card-text class="pt-6">
        <v-alert type="warning" prominent class="mb-4">
          <div class="text-h6">请立即复制保存此 Token！</div>
          <div class="mt-2">出于安全考虑，此 Token 将只显示这一次。关闭此窗口后将无法再次查看。</div>
        </v-alert>

        <div class="mb-4">
          <div class="text-subtitle-2 mb-2">Token 名称</div>
          <v-chip color="primary">{{ tokenName }}</v-chip>
        </div>

        <div>
          <div class="text-subtitle-2 mb-2">Token 值</div>
          <v-card variant="outlined" class="pa-4">
            <code style="word-break: break-all; font-size: 14px; line-height: 1.8">{{ token }}</code>
          </v-card>
        </div>

        <div class="mt-4 d-flex justify-center">
          <v-btn :color="copied ? 'success' : 'primary'" size="large" @click="copyToken">
            <v-icon start>{{ copied ? "mdi-check" : "mdi-content-copy" }}</v-icon>
            {{ copied ? "已复制" : "复制 Token" }}
          </v-btn>
        </div>

        <v-alert type="info" class="mt-4" density="compact">
          <strong>使用方式：</strong>在 API 请求中添加 Header: <code>Authorization: Bearer YOUR_TOKEN</code>
        </v-alert>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="elevated" @click="closeDialog">我已保存，关闭窗口</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
