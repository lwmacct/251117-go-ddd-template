<script setup lang="ts">
import { ref } from "vue";
import draggable from "vuedraggable";
import type { Menu } from "@/api";

interface Props {
  menus: Menu[];
  level?: number;
}

interface Emits {
  (e: "edit", menu: Menu): void;
  (e: "delete", menu: Menu): void;
  (e: "update:menus", menus: Menu[]): void;
}

const props = withDefaults(defineProps<Props>(), {
  level: 0,
});

const emit = defineEmits<Emits>();

const drag = ref(false);
const localMenus = ref(props.menus);

const handleDragEnd = () => {
  emit("update:menus", localMenus.value);
};
</script>

<template>
  <draggable
    v-model="localMenus"
    item-key="id"
    handle=".drag-handle"
    :group="{ name: 'menus' }"
    class="menu-tree"
    @start="drag = true"
    @end="
      drag = false;
      handleDragEnd();
    "
  >
    <template #item="{ element }">
      <div class="menu-item" :style="{ marginLeft: `${level * 24}px` }">
        <v-card variant="outlined" class="mb-2">
          <v-list-item>
            <template #prepend>
              <v-icon class="drag-handle" style="cursor: move">mdi-drag</v-icon>
            </template>

            <v-list-item-title class="d-flex align-center">
              <v-icon v-if="element.icon" size="small" class="mr-2">{{ element.icon }}</v-icon>
              {{ element.title }}
              <v-chip v-if="!element.visible" size="x-small" class="ml-2" color="warning">隐藏</v-chip>
            </v-list-item-title>

            <v-list-item-subtitle>
              {{ element.path }} <span class="text-caption ml-2">(排序: {{ element.order }})</span>
            </v-list-item-subtitle>

            <template #append>
              <v-btn icon="mdi-pencil" size="small" variant="text" @click="emit('edit', element)"></v-btn>
              <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="emit('delete', element)"></v-btn>
            </template>
          </v-list-item>

          <!-- 递归渲染子菜单 -->
          <div v-if="element.children && element.children.length > 0" class="ml-4">
            <MenuTree
              :menus="element.children"
              :level="level + 1"
              @edit="emit('edit', $event)"
              @delete="emit('delete', $event)"
              @update:menus="emit('update:menus', $event)"
            />
          </div>
        </v-card>
      </div>
    </template>
  </draggable>
</template>

<style scoped>
.menu-tree {
  width: 100%;
}

.menu-item {
  position: relative;
}

.drag-handle {
  cursor: move !important;
}
</style>
