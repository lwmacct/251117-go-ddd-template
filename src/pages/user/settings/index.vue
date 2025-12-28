<script setup lang="ts">
/**
 * 用户个人设置页面
 * 允许用户自定义配置偏好，可恢复到系统默认值
 *
 * 特性：
 * - 支持按 Tab 懒加载（减少首屏数据量）
 * - 支持字段验证（JSON Logic 规则）
 * - 支持字段依赖（条件可见/禁用）
 * - 响应式布局（宽屏垂直 Tabs，窄屏水平 Tabs）
 * - URL 参数同步（Tab 状态持久化）
 */
import { ref, computed, onMounted, watch } from "vue";
import jsonLogic from "json-logic-js";
import { useResponsiveTabs, type TabItem } from "@/composables";
import ResponsiveTabs from "@/components/ResponsiveTabs.vue";
import { useUserSettings } from "./composables/useUserSettings";
import UserSettingField from "./components/UserSettingField.vue";

const {
  loading,
  categories,
  schema,
  settingsMap,
  errorMessage,
  fetchCategories,
  fetchSchemaByCategory,
  isCategoryLoaded,
  updateSettingQuietly,
  resetSetting,
  clearMessages,
} = useUserSettings();

// Snackbar 状态
const snackbar = ref({
  show: false,
  message: "",
  color: "success",
});

const showSnackbar = (message: string, color = "success") => {
  snackbar.value = { show: true, message, color };
};

// 表单值（按 key 存储）
const formValues = ref<Record<string, unknown>>({});

// 正在重置的设置 key
const resettingKey = ref<string | null>(null);

// 正在保存的设置 keys
const savingKeys = ref<Set<string>>(new Set());

// 验证错误（按 key 存储）
const validationErrors = ref<Record<string, string[]>>({});

// Tab 列表（从 categories 元信息获取，而非 schema）
const tabs = computed<TabItem[]>(() =>
  categories.value.map((cat) => ({
    value: cat.category ?? "",
    label: cat.label ?? cat.category ?? "",
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

// 初始化表单值（使用 value 字段，即用户实际生效值）
const initFormValues = () => {
  const values: Record<string, unknown> = {};
  schema.value.forEach((cat) => {
    cat.groups?.forEach((group) => {
      group.settings?.forEach((setting) => {
        if (setting.key) {
          // 使用 value 字段（实际生效值），而非 default_value
          values[setting.key] = parseValue(setting.value, setting.value_type);
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

// =============== 验证逻辑 ===============

/**
 * 验证单个字段
 * 支持 JSON Logic 规则和简单规则两种格式
 */
const validateField = (key: string, value: unknown): boolean => {
  const setting = settingsMap.value.get(key);
  const validation = setting?.ui_config?.validation;
  if (!validation) return true;

  try {
    // 构建数据上下文（支持跨字段验证）
    const data = {
      value,
      key,
      settings: formValues.value,
    };

    // 执行 JSON Logic 验证
    const result = jsonLogic.apply(validation, data);
    if (!result) {
      const msg = (validation as { message?: string }).message || `${setting?.label || key}验证失败`;
      validationErrors.value[key] = [msg];
      return false;
    }
    delete validationErrors.value[key];
    return true;
  } catch (e) {
    console.error(`Validation error for ${key}:`, e);
    validationErrors.value[key] = [`${setting?.label || key}验证规则执行失败`];
    return false;
  }
};

// 获取字段验证错误
const getFieldErrors = (key: string): string[] => {
  return validationErrors.value[key] ?? [];
};

// =============== 依赖逻辑 ===============

interface DependsOnConfig {
  key: string;
  value?: unknown;
  operator?: "eq" | "ne" | "gt" | "lt";
}

/**
 * 检查设置项是否因依赖关系被禁用
 */
const isDisabled = (key: string): boolean => {
  const setting = settingsMap.value.get(key);
  const dependsOn = setting?.ui_config?.depends_on as DependsOnConfig | undefined;
  if (!dependsOn) return false;

  const depValue = formValues.value[dependsOn.key];
  const expectedValue = dependsOn.value;
  const operator = dependsOn.operator || "eq";

  switch (operator) {
    case "eq":
      return depValue !== expectedValue;
    case "ne":
      return depValue === expectedValue;
    case "gt":
      return !(Number(depValue) > Number(expectedValue));
    case "lt":
      return !(Number(depValue) < Number(expectedValue));
    default:
      return false;
  }
};

/**
 * 获取字段提示（优先显示禁用原因）
 */
const getFieldHint = (key: string): string | undefined => {
  const setting = settingsMap.value.get(key);
  const dependsOn = setting?.ui_config?.depends_on as DependsOnConfig | undefined;

  if (dependsOn && isDisabled(key)) {
    const depSetting = settingsMap.value.get(dependsOn.key);
    return `需要先启用「${depSetting?.label || dependsOn.key}」`;
  }

  return setting?.ui_config?.hint;
};

// 处理即时变更（switch/select）
const handleFieldChange = async (key: string, value: unknown) => {
  // 验证
  if (!validateField(key, value)) return;

  savingKeys.value.add(key);
  await updateSettingQuietly(key, value as object);
  savingKeys.value.delete(key);
};

// 处理失焦保存（text 类控件）
const handleFieldBlur = async (key: string) => {
  const value = formValues.value[key];
  if (value === undefined) return;

  // 验证
  if (!validateField(key, value)) return;

  savingKeys.value.add(key);
  await updateSettingQuietly(key, value as object);
  savingKeys.value.delete(key);
};

// 重置单个设置
const handleReset = async (key: string) => {
  resettingKey.value = key;
  const success = await resetSetting(key);
  if (success) {
    // 更新表单值为默认值
    const setting = settingsMap.value.get(key);
    if (setting) {
      formValues.value[key] = parseValue(setting.default_value, setting.value_type);
    }
    // 清除验证错误
    delete validationErrors.value[key];
    // 显示临时提示
    showSnackbar("已恢复默认值");
  }
  resettingKey.value = null;
};

onMounted(async () => {
  // 1. 先获取分类列表（用于渲染 tabs）
  await fetchCategories();

  // 2. 加载当前 tab 的数据（URL 参数或默认第一个）
  const targetTab = currentTab.value || categories.value[0]?.category || "general";
  if (!isCategoryLoaded(targetTab)) {
    await fetchSchemaByCategory(targetTab);
  }
});
</script>

<template>
  <div class="user-settings-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">个人设置</h1>
        <p class="text-body-2 text-medium-emphasis mb-6">
          自定义您的偏好设置。已修改的设置项可点击
          <v-icon size="x-small" color="info" class="mx-1">mdi-restore</v-icon>
          按钮恢复默认值。
        </p>
      </v-col>
    </v-row>

    <!-- 错误提示 -->
    <v-row v-if="errorMessage">
      <v-col cols="12">
        <v-alert type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
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
          暂无可配置的设置项。
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
                    <UserSettingField
                      v-if="setting.key"
                      v-model="formValues[setting.key]"
                      :setting="setting"
                      :disabled="savingKeys.has(setting.key) || isDisabled(setting.key)"
                      :resetting="resettingKey === setting.key"
                      :error-messages="getFieldErrors(setting.key)"
                      :hint="getFieldHint(setting.key)"
                      @change="handleFieldChange(setting.key!, $event)"
                      @blur="handleFieldBlur(setting.key!)"
                      @reset="handleReset(setting.key!)"
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

    <!-- 临时提示 -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="3000" location="bottom">
      {{ snackbar.message }}
    </v-snackbar>
  </div>
</template>

<style scoped>
.user-settings-page {
  width: 100%;
}
</style>
