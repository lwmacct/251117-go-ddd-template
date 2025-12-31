<script setup lang="ts">
import { ref, watch } from "vue";
import { adminUserApi, extractList } from "@/api";
import type { UserUserWithRolesDTO } from "@models";

interface Props {
  modelValue: number | null;
  label?: string;
  disabled?: boolean;
  hideDetails?: boolean;
}

defineProps<Props>();
const emit = defineEmits<{
  (e: "update:modelValue", value: number | null): void;
}>();

const search = ref("");
const loading = ref(false);
const items = ref<UserUserWithRolesDTO[]>([]);

/**
 * 搜索用户
 */
const searchUsers = async () => {
  if (!search.value || search.value.length < 2) {
    items.value = [];
    return;
  }

  loading.value = true;
  try {
    const response = await adminUserApi.apiSystemUsersGet(20, 1, search.value);
    const result = extractList<UserUserWithRolesDTO>(response.data);
    items.value = result.data;
  } catch (err) {
    console.error("Failed to search users:", err);
    items.value = [];
  } finally {
    loading.value = false;
  }
};

// 防抖搜索
let debounceTimer: ReturnType<typeof setTimeout> | null = null;
watch(search, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }
  debounceTimer = setTimeout(() => {
    searchUsers();
  }, 300);
});

/**
 * 格式化用户显示
 */
const displayName = (user: UserUserWithRolesDTO) => {
  if (user.real_name) {
    return `${user.real_name} (@${user.username})`;
  }
  return `@${user.username}`;
};
</script>

<template>
  <v-autocomplete
    :model-value="modelValue"
    :items="items"
    :item-title="displayName"
    item-value="id"
    :label="label"
    :loading="loading"
    :disabled="disabled"
    :hide-details="hideDetails"
    :search="search"
    variant="outlined"
    clearable
    return-object
    @update:model-value="emit('update:modelValue', ($event as UserUserWithRolesDTO | null)?.id || null)"
    @update:search="search = $event"
  >
    <template #chip="{ props: chipProps, item }">
      <v-chip v-bind="chipProps" :text="item.raw.username" />
    </template>
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps" :title="item.raw.username" :subtitle="item.raw.email || item.raw.real_name">
        <template #prepend>
          <v-avatar color="primary" size="32">
            <span>{{ item.raw.username?.[0]?.toUpperCase() || "?" }}</span>
          </v-avatar>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>
