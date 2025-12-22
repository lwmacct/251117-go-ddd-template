<script setup lang="ts">
import { ref, watch } from "vue";
import type { CreateTokenRequest } from "@/types/user";

interface Props {
  modelValue: boolean;
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: CreateTokenRequest): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const formData = ref<CreateTokenRequest>({
  name: "",
  permissions: [],
  expires_at: undefined,
  ip_whitelist: [],
});

const valid = ref(false);
const form = ref();
const expiresEnabled = ref(false);
const ipWhitelistEnabled = ref(false);
const ipWhitelistText = ref("");

const rules = {
  name: [(v: string) => !!v || "Token 名称不能为空", (v: string) => (v && v.length >= 3) || "Token 名称至少3个字符"],
};

watch(
  () => props.modelValue,
  (newVal) => {
    if (!newVal) {
      resetForm();
    }
  }
);

const resetForm = () => {
  formData.value = {
    name: "",
    permissions: [],
    expires_at: undefined,
    ip_whitelist: [],
  };
  expiresEnabled.value = false;
  ipWhitelistEnabled.value = false;
  ipWhitelistText.value = "";
  form.value?.resetValidation();
};

const closeDialog = () => {
  emit("update:modelValue", false);
};

const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  const data: CreateTokenRequest = {
    name: formData.value.name,
    permissions: formData.value.permissions?.length ? formData.value.permissions : undefined,
    expires_at: expiresEnabled.value ? formData.value.expires_at : undefined,
    ip_whitelist:
      ipWhitelistEnabled.value && ipWhitelistText.value
        ? ipWhitelistText.value
            .split("\n")
            .map((ip) => ip.trim())
            .filter((ip) => ip)
        : undefined,
  };

  emit("save", data);
  closeDialog();
};

const expirationOptions = [
  { title: "7 天", value: 7 },
  { title: "30 天", value: 30 },
  { title: "90 天", value: 90 },
  { title: "1 年", value: 365 },
];

const setExpiration = (days: number) => {
  const date = new Date();
  date.setDate(date.getDate() + days);
  formData.value.expires_at = date.toISOString();
};
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="600"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title>
        <span class="text-h5">创建 Personal Access Token</span>
      </v-card-title>

      <v-card-text>
        <v-form ref="form" v-model="valid">
          <v-text-field
            v-model="formData.name"
            label="Token 名称"
            :rules="rules.name"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            hint="给这个 Token 一个描述性的名称"
            persistent-hint
          ></v-text-field>

          <div class="mt-4">
            <v-checkbox v-model="expiresEnabled" label="设置过期时间" density="compact" hide-details></v-checkbox>

            <div v-if="expiresEnabled" class="ml-8 mt-2">
              <v-chip-group>
                <v-chip
                  v-for="option in expirationOptions"
                  :key="option.value"
                  size="small"
                  @click="setExpiration(option.value)"
                >
                  {{ option.title }}
                </v-chip>
              </v-chip-group>

              <v-text-field
                v-model="formData.expires_at"
                label="过期时间"
                type="datetime-local"
                variant="outlined"
                density="compact"
                class="mt-2"
                hint="留空则永不过期"
              ></v-text-field>
            </div>
          </div>

          <div class="mt-4">
            <v-checkbox v-model="ipWhitelistEnabled" label="IP 白名单" density="compact" hide-details></v-checkbox>

            <div v-if="ipWhitelistEnabled" class="ml-8 mt-2">
              <v-textarea
                v-model="ipWhitelistText"
                label="IP 地址列表（每行一个）"
                variant="outlined"
                rows="3"
                hint="例如: 192.168.1.1"
                persistent-hint
              ></v-textarea>
            </div>
          </div>
        </v-form>

        <v-alert type="info" class="mt-4" density="compact">
          <strong>注意：</strong>Token 只会在创建时显示一次，请妥善保存。
        </v-alert>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" :disabled="!valid" @click="handleSave">创建</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
