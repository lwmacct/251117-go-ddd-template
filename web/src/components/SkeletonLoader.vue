<script setup lang="ts">
/**
 * 骨架屏加载组件
 * 提供多种预设骨架样式
 */

interface Props {
  /** 骨架类型 */
  type?: "text" | "avatar" | "button" | "card" | "table" | "list" | "custom";
  /** 是否加载中 */
  loading?: boolean;
  /** 行数（用于 text、list 类型） */
  lines?: number;
  /** 表格行数 */
  tableRows?: number;
  /** 表格列数 */
  tableCols?: number;
}

withDefaults(defineProps<Props>(), {
  type: "text",
  loading: true,
  lines: 3,
  tableRows: 5,
  tableCols: 4,
});
</script>

<template>
  <div v-if="loading" class="skeleton-loader">
    <!-- 文本骨架 -->
    <template v-if="type === 'text'">
      <v-skeleton-loader
        v-for="i in lines"
        :key="i"
        type="text"
        :style="{ width: i === lines ? '60%' : '100%' }"
        class="mb-2"
      />
    </template>

    <!-- 头像骨架 -->
    <template v-else-if="type === 'avatar'">
      <div class="d-flex align-center">
        <v-skeleton-loader type="avatar" class="mr-4" />
        <div class="flex-grow-1">
          <v-skeleton-loader type="text" style="width: 40%" class="mb-2" />
          <v-skeleton-loader type="text" style="width: 60%" />
        </div>
      </div>
    </template>

    <!-- 按钮骨架 -->
    <template v-else-if="type === 'button'">
      <v-skeleton-loader type="button" />
    </template>

    <!-- 卡片骨架 -->
    <template v-else-if="type === 'card'">
      <v-skeleton-loader type="card" />
    </template>

    <!-- 表格骨架 -->
    <template v-else-if="type === 'table'">
      <div class="skeleton-table">
        <!-- 表头 -->
        <div class="d-flex mb-4 pa-2 bg-grey-lighten-4 rounded">
          <v-skeleton-loader
            v-for="col in tableCols"
            :key="`header-${col}`"
            type="text"
            class="flex-grow-1 mx-2"
            style="height: 20px"
          />
        </div>
        <!-- 表格行 -->
        <div v-for="row in tableRows" :key="`row-${row}`" class="d-flex mb-3 pa-2">
          <v-skeleton-loader
            v-for="col in tableCols"
            :key="`cell-${row}-${col}`"
            type="text"
            class="flex-grow-1 mx-2"
          />
        </div>
      </div>
    </template>

    <!-- 列表骨架 -->
    <template v-else-if="type === 'list'">
      <div v-for="i in lines" :key="i" class="d-flex align-center mb-4">
        <v-skeleton-loader type="avatar" class="mr-4" />
        <div class="flex-grow-1">
          <v-skeleton-loader type="text" style="width: 30%" class="mb-2" />
          <v-skeleton-loader type="text" style="width: 50%" />
        </div>
        <v-skeleton-loader type="button" style="width: 80px" />
      </div>
    </template>

    <!-- 自定义骨架 -->
    <template v-else>
      <slot name="skeleton" />
    </template>
  </div>

  <!-- 实际内容 -->
  <slot v-else />
</template>

<style scoped>
.skeleton-loader {
  width: 100%;
}

.skeleton-table {
  width: 100%;
}
</style>
