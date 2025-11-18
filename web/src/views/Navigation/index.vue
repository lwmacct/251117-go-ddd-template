<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import type { MenuItem } from './types.ts'

interface Props {
  /**
   * 菜单项列表
   */
  items: MenuItem[]
  /**
   * 抽屉宽度
   */
  width?: number | string
  /**
   * Rail 模式宽度（折叠后宽度）
   */
  railWidth?: number
  /**
   * 是否显示边框
   */
  border?: boolean
  /**
   * 是否固定
   */
  permanent?: boolean
  /**
   * 是否显示图标
   */
  showIcon?: boolean
  /**
   * 颜色主题
   */
  color?: string
  /**
   * 默认是否折叠
   */
  defaultRail?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: 220,
  railWidth: 56, // 与 Navbar 图标宽度对齐：v-app-bar-nav-icon 默认 56px
  border: true,
  permanent: true,
  showIcon: true,
  color: undefined,
  defaultRail: false,
})

const route = useRoute()

// 折叠状态
const rail = ref(props.defaultRail)

/**
 * 切换折叠状态
 */
const toggleRail = () => {
  rail.value = !rail.value
}

/**
 * 判断菜单项是否激活
 */
const isActive = (item: MenuItem): boolean => {
  if (item.exact) {
    return route.path === item.path
  }
  return route.path.startsWith(item.path)
}

/**
 * 获取激活的菜单项
 */
const activeItem = computed(() => {
  return props.items.find((item) => isActive(item))
})
</script>

<template>
  <v-navigation-drawer :model-value="true" :rail="rail" :width="width" :rail-width="railWidth" :border="border"
    :permanent="permanent" :color="color">
    <!-- 菜单列表 -->
    <v-list nav density="compact" bg-color="transparent">
      <v-list-item v-for="item in items" :key="item.path" :to="item.path" :active="isActive(item)"
        :prepend-icon="showIcon ? item.icon : undefined" :title="item.title" :value="item.path" rounded="lg">
        <!-- 折叠状态下的 Tooltip -->
        <v-tooltip v-if="rail" activator="parent" location="end">
          {{ item.title }}
        </v-tooltip>

        <!-- 可选：显示徽章 -->
        <template v-if="item.badge && !rail" #append>
          <v-badge :content="item.badge" :color="item.badgeColor || 'error'" inline />
        </template>
      </v-list-item>
    </v-list>

    <!-- 折叠/展开按钮（底部） -->
    <template #append>
      <div class="pa-2">
        <v-btn :icon="rail ? 'mdi-chevron-right' : 'mdi-chevron-left'" variant="text" block @click="toggleRail" />
      </div>
    </template>
  </v-navigation-drawer>
</template>
