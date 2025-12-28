/**
 * Settings 页面专用的分类懒加载 Composable
 *
 * 基于通用 useLazyLoader 封装，提供 Settings Schema 场景的便捷 API：
 * - 按 category key 加载
 * - 自动按 category 字段合并（更新已有分类或追加新分类）
 *
 * @example
 * ```typescript
 * const { data: schema, load, isLoaded, loading } = useLazyCategorySchema(
 *   async (categoryKey) => {
 *     const response = await settingsApi.get(categoryKey);
 *     return response.data.data ?? [];
 *   }
 * );
 *
 * // 加载指定分类
 * await load('general');
 * await load('notification');
 *
 * // schema 中包含所有已加载分类的数据
 * ```
 */
import { useLazyLoader, type UseLazyLoaderReturn } from "./useLazyLoader";

/**
 * Settings Schema 数据项的基础类型约束
 * 要求数据项必须有 category 字段
 */
export interface CategorySchemaItem {
  category?: string;
}

/**
 * Settings 分类懒加载返回类型
 * 扩展自 useLazyLoader，添加别名以匹配 Settings 业务语义
 */
export interface UseLazyCategorySchemaReturn<T extends CategorySchemaItem> extends UseLazyLoaderReturn<string, T> {
  /** 已加载的 Schema 数据（data 的别名，语义更清晰） */
  schema: UseLazyLoaderReturn<string, T>["data"];

  /** 按分类加载 Schema（load 的别名） */
  fetchSchemaByCategory: UseLazyLoaderReturn<string, T>["load"];

  /** 检查分类是否已加载（isLoaded 的别名） */
  isCategoryLoaded: UseLazyLoaderReturn<string, T>["isLoaded"];

  /** 已加载的分类 key 集合（loaded 的别名） */
  loadedCategories: UseLazyLoaderReturn<string, T>["loaded"];
}

/**
 * Settings 分类懒加载 Composable
 *
 * @param fetcher 按 category key 获取数据的函数
 * @returns 懒加载状态和方法（包含 Settings 语义的别名）
 */
export function useLazyCategorySchema<T extends CategorySchemaItem>(
  fetcher: (categoryKey: string) => Promise<T[]>,
): UseLazyCategorySchemaReturn<T> {
  const loader = useLazyLoader<string, T>({
    fetcher,
    merge: (existing, incoming) => {
      // 按 category 合并：更新已有分类或追加新分类
      const result = [...existing];

      incoming.forEach((item) => {
        const existingIndex = result.findIndex((e) => e.category === item.category);
        if (existingIndex !== -1) {
          // 更新已有分类
          result[existingIndex] = item;
        } else {
          // 追加新分类
          result.push(item);
        }
      });

      return result;
    },
    getKey: (item) => item.category ?? "",
  });

  return {
    // 原始返回值
    ...loader,

    // Settings 语义别名
    schema: loader.data,
    fetchSchemaByCategory: loader.load,
    isCategoryLoaded: loader.isLoaded,
    loadedCategories: loader.loaded,
  };
}
