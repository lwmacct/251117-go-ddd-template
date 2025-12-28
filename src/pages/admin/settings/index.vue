<script setup lang="ts">
/**
 * 系统设置页面 - 动态渲染版本
 * 根据后端 Schema 数据动态生成设置界面
 *
 * 特性：
 * - 支持按 Tab 懒加载（减少首屏数据量）
 * - 支持字段验证（JSON Logic 规则）
 * - 支持字段依赖（条件可见/禁用）
 * - 响应式布局（宽屏垂直 Tabs，窄屏水平 Tabs）
 * - URL 参数同步（Tab 状态持久化）
 */
import { ref, computed, onMounted, watch } from "vue";
import { useResponsiveTabs, type TabItem } from "@/composables";
import ResponsiveTabs from "@/components/ResponsiveTabs.vue";
import { useSettings } from "./composables/useSettings";
import { useSettingsDependency } from "./composables/useSettingsDependency";
import { useJsonLogicValidation } from "./composables/useJsonLogicValidation";
import DynamicSettingField from "./components/DynamicSettingField.vue";

const {
  loading,
  categories,
  schema,
  settingsMap,
  errorMessage,
  successMessage,
  fetchCategories,
  fetchSchemaByCategory,
  isCategoryLoaded,
  updateSettingQuietly,
  clearMessages,
} = useSettings();

// 表单值（按 key 存储）
const formValues = ref<Record<string, unknown>>({});

// 正在保存的设置 keys
const savingKeys = ref<Set<string>>(new Set());

// 依赖关系处理
const { isDisabled, getFinalHint } = useSettingsDependency(settingsMap, formValues);

// JSON Logic 验证
const { validate, getError } = useJsonLogicValidation(schema, formValues);

// 获取字段错误消息（转换为数组格式供 Vuetify 使用）
const getFieldErrors = (key: string): string[] => {
  const error = getError(key);
  return error ? [error] : [];
};

// Tab 列表（从 categories 元信息获取，而非 schema）
const tabs = computed<TabItem[]>(() =>
  categories.value.map((cat) => ({
    value: cat.key ?? "",
    label: cat.label ?? cat.key ?? "",
    icon: cat.icon ?? "mdi-cog",
  })),
);

// 响应式 Tabs（含懒加载回调）
const { currentTab, isVertical, handleTabChange } = useResponsiveTabs({
  defaultTab: "general",
  onTabChange: async (tab) => {
    // 懒加载：仅加载未加载的 category
    if (!isCategoryLoaded(tab)) {
      await fetchSchemaByCategory(tab);
    }
  },
});

// 获取指定 Tab 的分组
const getGroupsByCategory = (category: string) => {
  const cat = schema.value.find((c) => c.category === category);
  return cat?.groups ?? [];
};

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
});

// 处理即时变更（switch/select）
const handleFieldChange = async (key: string, value: unknown) => {
  // 先验证
  const error = validate(key);
  if (error) return; // 验证失败，不保存

  savingKeys.value.add(key);
  await updateSettingQuietly(key, value as object);
  savingKeys.value.delete(key);
};

// 处理失焦保存（text 类控件）
const handleFieldBlur = async (key: string) => {
  const value = formValues.value[key];
  if (value === undefined) return;

  // 先验证
  const error = validate(key);
  if (error) return; // 验证失败，不保存

  savingKeys.value.add(key);
  await updateSettingQuietly(key, value as object);
  savingKeys.value.delete(key);
};

onMounted(async () => {
  // 1. 先获取分类列表（用于渲染 tabs）
  await fetchCategories();

  // 2. 加载当前 tab 的数据（URL 参数或默认第一个）
  const targetTab = currentTab.value || categories.value[0]?.key || "general";
  if (!isCategoryLoaded(targetTab)) {
    await fetchSchemaByCategory(targetTab);
  }
});
</script>

<template>
  <div class="settings-page">
    <v-row>
      <v-col cols="12">
        <div class="d-flex align-center mb-2">
          <h1 class="text-h4">系统设置</h1>
          <v-spacer />
          <v-btn variant="outlined" color="primary" to="/admin/setting-categories">
            <v-icon start>mdi-folder-cog</v-icon>
            管理分类
          </v-btn>
        </div>
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
          <ResponsiveTabs :model-value="currentTab" :tabs="tabs" :vertical="isVertical" @update:model-value="handleTabChange">
            <template v-for="tab in tabs" :key="tab.value" #[tab.value]>
              <!-- 按 Group 渲染 -->
              <template v-for="group in getGroupsByCategory(tab.value)" :key="group.group">
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
                      :disabled="isDisabled(setting) || savingKeys.has(setting.key)"
                      :hint="getFinalHint(setting)"
                      :error-messages="getFieldErrors(setting.key)"
                      @change="handleFieldChange(setting.key!, $event)"
                      @blur="handleFieldBlur(setting.key!)"
                    />
                  </v-col>
                </v-row>
              </template>

              <!-- 空分组提示（仅在分类加载完成后显示） -->
              <v-alert
                v-if="isCategoryLoaded(tab.value) && getGroupsByCategory(tab.value).length === 0"
                type="info"
                variant="tonal"
                class="mt-4"
              >
                <v-icon start>mdi-information</v-icon>
                该分类下暂无配置项
              </v-alert>
            </template>
          </ResponsiveTabs>
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
