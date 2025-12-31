<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { ProductProductDTO, ProductCreateProductDTO, ProductUpdateProductDTO } from "@models";

interface Props {
  modelValue: boolean;
  product?: ProductProductDTO | null;
  mode: "create" | "edit";
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: ProductCreateProductDTO | ProductUpdateProductDTO): void;
}>();

// 表单数据
const formData = ref<{
  name: string;
  description: string;
  price: number;
  status: string;
}>({
  name: "",
  description: "",
  price: 0,
  status: "active",
});

// 表单引用
const formRef = ref<HTMLFormElement | null>(null);

// 表单验证规则
const rules = {
  name: [(v: string) => !!v || "产品名称不能为空", (v: string) => (v && v.length >= 2) || "产品名称至少2个字符"],
  price: [(v: number) => v >= 0 || "价格不能为负数"],
  status: [(v: string) => !!v || "请选择状态"],
};

// 监听 product 变化，更新表单数据
watch(
  () => props.product,
  (newProduct) => {
    if (newProduct) {
      formData.value = {
        name: newProduct.name || "",
        description: newProduct.description || "",
        price: newProduct.price || 0,
        status: newProduct.status || "active",
      };
    } else {
      formData.value = {
        name: "",
        description: "",
        price: 0,
        status: "active",
      };
    }
  },
  { immediate: true },
);

// 状态选项
const statusOptions = [
  { title: "启用", value: "active" },
  { title: "禁用", value: "inactive" },
];

// 对话框标题
const dialogTitle = computed(() => {
  return props.mode === "create" ? "新建产品" : "编辑产品";
});

/**
 * 保存处理
 */
const handleSave = async () => {
  // 验证表单
  const valid = await formRef.value?.validate?.();
  if (valid?.length > 0) {
    return; // 验证失败
  }

  if (props.mode === "create") {
    emit("save", {
      name: formData.value.name,
      description: formData.value.description || undefined,
      price: formData.value.price,
      status: formData.value.status as "active" | "inactive",
    } as ProductCreateProductDTO);
  } else {
    emit("save", {
      name: formData.value.name || undefined,
      description: formData.value.description || undefined,
      price: formData.value.price || undefined,
      status: formData.value.status || undefined,
    } as ProductUpdateProductDTO);
  }
  emit("update:modelValue", false);
};

/**
 * 关闭对话框
 */
const handleClose = () => {
  emit("update:modelValue", false);
};
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="600" @click:outside="handleClose">
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form ref="formRef">
          <v-text-field v-model="formData.name" label="产品名称" :rules="rules.name" variant="outlined" required class="mb-2" />

          <v-textarea v-model="formData.description" label="产品描述" variant="outlined" rows="3" auto-grow class="mb-2" />

          <v-text-field
            v-model.number="formData.price"
            label="价格"
            :rules="rules.price"
            type="number"
            variant="outlined"
            required
            prefix="¥"
            class="mb-2"
          />

          <v-select
            v-model="formData.status"
            label="状态"
            :items="statusOptions"
            item-title="title"
            item-value="value"
            :rules="rules.status"
            variant="outlined"
            required
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="handleClose">取消</v-btn>
        <v-btn color="primary" variant="elevated" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
