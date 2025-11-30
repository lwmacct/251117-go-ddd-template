<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter, type RouteLocationMatched } from "vue-router";

/**
 * 应用面包屑导航组件
 * 自动根据当前路由生成面包屑路径
 */

interface Props {
  /** 自定义首页图标 */
  homeIcon?: string;
  /** 自定义首页文本 */
  homeText?: string;
  /** 是否显示首页链接 */
  showHome?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  homeIcon: "mdi-home",
  homeText: "首页",
  showHome: true,
});

const route = useRoute();
const router = useRouter();

interface BreadcrumbItem {
  title: string;
  to?: string;
  disabled?: boolean;
  icon?: string;
}

/**
 * 根据路由匹配信息生成面包屑项
 */
const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const items: BreadcrumbItem[] = [];

  // 添加首页
  if (props.showHome) {
    items.push({
      title: props.homeText,
      to: "/",
      icon: props.homeIcon,
    });
  }

  // 过滤有效的路由匹配项
  const matchedRoutes = route.matched.filter((record: RouteLocationMatched) => {
    // 跳过没有 title 的路由
    if (!record.meta?.title) return false;
    // 跳过首页（已单独添加）
    if (record.path === "/" || record.path === "") return false;
    return true;
  });

  // 添加匹配的路由
  matchedRoutes.forEach((record: RouteLocationMatched, index: number) => {
    const isLast = index === matchedRoutes.length - 1;
    const title = record.meta?.title as string;
    const icon = record.meta?.icon as string | undefined;

    items.push({
      title,
      to: isLast ? undefined : record.path,
      disabled: isLast,
      icon,
    });
  });

  return items;
});

/**
 * 导航到指定路由
 */
const navigateTo = (item: BreadcrumbItem) => {
  if (item.to) {
    router.push(item.to);
  }
};
</script>

<template>
  <nav class="app-breadcrumb" aria-label="面包屑导航">
    <ol class="d-flex align-center flex-wrap pa-0 ma-0" style="list-style: none">
      <li v-for="(item, index) in breadcrumbs" :key="index" class="d-flex align-center">
        <!-- 分隔符 -->
        <v-icon v-if="index > 0" icon="mdi-chevron-right" size="small" class="mx-1 text-medium-emphasis" />

        <!-- 面包屑项 -->
        <span
          :class="['d-flex align-center', item.disabled ? 'text-medium-emphasis' : 'breadcrumb-link cursor-pointer']"
          @click="!item.disabled && navigateTo(item)"
        >
          <v-icon v-if="item.icon && index === 0" :icon="item.icon" size="small" class="mr-1" />
          <span class="text-body-2">{{ item.title }}</span>
        </span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
.app-breadcrumb {
  min-height: 32px;
}

.cursor-pointer {
  cursor: pointer;
}

.breadcrumb-link {
  cursor: pointer;
}

.breadcrumb-link:hover {
  text-decoration: underline;
}
</style>
