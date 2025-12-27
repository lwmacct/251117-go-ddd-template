<script setup lang="ts">
/**
 * 系统设置页面 - 动态渲染版本
 * 根据后端 Schema 数据动态生成设置界面
 */
import { ref, computed, onMounted, watch } from "vue";
import { useSettings } from "./composables/useSettings";
import { useSettingsDependency } from "./composables/useSettingsDependency";
import { useJsonLogicValidation } from "./composables/useJsonLogicValidation";
import DynamicSettingField from "./components/DynamicSettingField.vue";

const { loading, saving, schema, settingsMap, errorMessage, successMessage, fetchSchema, batchUpdateSettings, clearMessages } =
  useSettings();

// 当前 Tab
const currentTab = ref("");

// 表单值（按 key 存储）
const formValues = ref<Record<string, unknown>>({});

// 依赖关系处理
const { isDisabled, getFinalHint } = useSettingsDependency(settingsMap, formValues);

// JSON Logic 验证
const { errors: validationErrors, validateAll, getError, clearErrors } = useJsonLogicValidation(schema, formValues);

// 获取字段错误消息（转换为数组格式供 Vuetify 使用）
const getFieldErrors = (key: string): string[] => {
  const error = getError(key);
  return error ? [error] : [];
};

// Tab 列表
const tabs = computed(() =>
  schema.value.map((cat) => ({
    value: cat.category ?? "",
    label: cat.label ?? cat.category ?? "",
    icon: cat.icon ?? "mdi-cog",
  })),
);

// 当前 Tab 的分组
const currentGroups = computed(() => {
  const category = schema.value.find((c) => c.category === currentTab.value);
  return category?.groups ?? [];
});

// 解析值类型
const parseValue = (value: unknown, valueType: string | undefined): unknown => {
  if (value === undefined || value === null) return "";
  // 如果已经是正确的类型，直接返回
  if (valueType === "boolean" && typeof value === "boolean") return value;
  if (valueType === "number" && typeof value === "number") return value;
  if (valueType === "json" && typeof value === "object") return value;
  // 字符串形式的解析（兼容旧数据）
  if (typeof value === "string") {
    switch (valueType) {
      case "boolean":
        return value === "true";
      case "number":
        return Number(value) || 0;
      case "json":
        try {
          return JSON.parse(value);
        } catch {
          return {};
        }
      default:
        return value;
    }
  }
  return value;
};

// 初始化表单值
const initFormValues = () => {
  const values: Record<string, unknown> = {};
  schema.value.forEach((cat) => {
    cat.groups?.forEach((group) => {
      group.settings?.forEach((setting) => {
        if (setting.key) {
          values[setting.key] = parseValue(setting.default_value, setting.value_type);
        }
      });
    });
  });
  formValues.value = values;
};

// 监听 schema 变化，初始化表单值
watch(schema, () => {
  initFormValues();
  // 设置默认 Tab
  if (schema.value.length > 0 && !currentTab.value) {
    currentTab.value = schema.value[0]?.category ?? "general";
  }
});

// 保存所有设置
const saveAllSettings = async () => {
  // 执行前端验证
  const errors = validateAll();
  if (errors.size > 0) {
    // 滚动到第一个错误字段（可选）
    return;
  }

  const updates = Object.entries(formValues.value).map(([key, value]) => ({
    key,
    value: value as object,
  }));
  await batchUpdateSettings(updates);
};

// 重置表单
const resetForm = () => {
  initFormValues();
  clearErrors();
};

onMounted(async () => {
  await fetchSchema();
});
</script>

<template>
  <div class="settings-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">系统设置</h1>
        <p class="text-body-2 text-medium-emphasis mb-6">配置系统参数和偏好设置</p>
      </v-col>
    </v-row>

    <!-- 消息提示 -->
    <v-row v-if="errorMessage || successMessage">
      <v-col cols="12">
        <v-alert v-if="errorMessage" type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
        </v-alert>
        <v-alert v-if="successMessage" type="success" closable @click:close="clearMessages">
          {{ successMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <!-- 加载状态 -->
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />

    <!-- 空状态 -->
    <v-row v-if="!loading && schema.length === 0">
      <v-col cols="12">
        <v-alert type="info" variant="tonal">
          <v-icon start>mdi-information</v-icon>
          暂无配置项，请在数据库中初始化设置数据。
        </v-alert>
      </v-col>
    </v-row>

    <!-- 设置表单 -->
    <v-row v-if="schema.length > 0">
      <v-col cols="12">
        <v-card>
          <!-- 动态 Tab -->
          <v-tabs v-model="currentTab" bg-color="primary">
            <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value" :prepend-icon="tab.icon">
              {{ tab.label }}
            </v-tab>
          </v-tabs>

          <v-card-text class="pa-6">
            <v-tabs-window v-model="currentTab">
              <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
                <!-- 按 Group 渲染 -->
                <template v-for="group in currentGroups" :key="group.group">
                  <!-- 分组标题 -->
                  <div v-if="group.label" class="text-subtitle-1 font-weight-medium mb-3 mt-4">
                    {{ group.label }}
                  </div>

                  <v-row>
                    <v-col
                      v-for="setting in group.settings"
                      :key="setting.key"
                      cols="12"
                      :md="setting.ui_config?.input_type === 'switch' ? 12 : 6"
                    >
                      <DynamicSettingField
                        v-if="setting.key"
                        v-model="formValues[setting.key]"
                        :setting="setting"
                        :disabled="isDisabled(setting)"
                        :hint="getFinalHint(setting)"
                        :error-messages="getFieldErrors(setting.key)"
                      />
                    </v-col>
                  </v-row>
                </template>

                <!-- 空分组提示 -->
                <v-alert v-if="currentGroups.length === 0" type="info" variant="tonal" class="mt-4">
                  <v-icon start>mdi-information</v-icon>
                  该分类下暂无配置项
                </v-alert>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card-text>

          <v-card-actions class="pa-6">
            <v-spacer />
            <v-btn variant="text" :disabled="saving" @click="resetForm">重置</v-btn>
            <v-btn color="primary" :loading="saving" @click="saveAllSettings">保存所有设置</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.settings-page {
  width: 100%;
}
</style>
