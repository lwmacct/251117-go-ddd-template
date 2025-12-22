/**
 * 数组工具函数
 * 提供常用的数组操作方法
 */

// ============================================================================
// 分组与分块
// ============================================================================

/**
 * 数组分块
 * @example
 * chunk([1, 2, 3, 4, 5], 2) // [[1, 2], [3, 4], [5]]
 */
export function chunk<T>(arr: T[], size: number): T[][] {
  if (size <= 0) return [];

  const result: T[][] = [];
  for (let i = 0; i < arr.length; i += size) {
    result.push(arr.slice(i, i + size));
  }
  return result;
}

/**
 * 数组分组
 * @example
 * groupBy([{ type: 'a', value: 1 }, { type: 'b', value: 2 }], 'type')
 * // { a: [{ type: 'a', value: 1 }], b: [{ type: 'b', value: 2 }] }
 */
export function groupBy<T>(arr: T[], key: keyof T | ((item: T) => string | number)): Record<string, T[]> {
  return arr.reduce(
    (groups, item) => {
      const groupKey = typeof key === "function" ? String(key(item)) : String(item[key]);

      if (!groups[groupKey]) {
        groups[groupKey] = [];
      }
      groups[groupKey].push(item);

      return groups;
    },
    {} as Record<string, T[]>
  );
}

/**
 * 按条件分组为两个数组
 * @example
 * partition([1, 2, 3, 4, 5], n => n % 2 === 0) // [[2, 4], [1, 3, 5]]
 */
export function partition<T>(arr: T[], predicate: (item: T, index: number) => boolean): [T[], T[]] {
  const pass: T[] = [];
  const fail: T[] = [];

  arr.forEach((item, index) => {
    if (predicate(item, index)) {
      pass.push(item);
    } else {
      fail.push(item);
    }
  });

  return [pass, fail];
}

// ============================================================================
// 去重与交集
// ============================================================================

/**
 * 数组去重
 * @example
 * unique([1, 2, 2, 3]) // [1, 2, 3]
 * unique([{ id: 1 }, { id: 1 }], 'id') // [{ id: 1 }]
 */
export function unique<T>(arr: T[], key?: keyof T | ((item: T) => unknown)): T[] {
  if (!key) {
    return [...new Set(arr)];
  }

  const seen = new Set<unknown>();
  return arr.filter((item) => {
    const k = typeof key === "function" ? key(item) : item[key];
    if (seen.has(k)) {
      return false;
    }
    seen.add(k);
    return true;
  });
}

/**
 * 数组交集
 * @example
 * intersection([1, 2, 3], [2, 3, 4]) // [2, 3]
 */
export function intersection<T>(...arrays: T[][]): T[] {
  if (arrays.length === 0) return [];
  if (arrays.length === 1) return [...arrays[0]];

  return arrays.reduce((a, b) => a.filter((item) => b.includes(item)));
}

/**
 * 数组差集
 * @example
 * difference([1, 2, 3], [2, 3, 4]) // [1]
 */
export function difference<T>(arr: T[], ...others: T[][]): T[] {
  const otherSet = new Set(others.flat());
  return arr.filter((item) => !otherSet.has(item));
}

/**
 * 数组并集
 * @example
 * union([1, 2], [2, 3], [3, 4]) // [1, 2, 3, 4]
 */
export function union<T>(...arrays: T[][]): T[] {
  return [...new Set(arrays.flat())];
}

// ============================================================================
// 查找与索引
// ============================================================================

/**
 * 查找第一个匹配项的索引
 * @example
 * findIndex([{ id: 1 }, { id: 2 }], item => item.id === 2) // 1
 */
export function findIndex<T>(arr: T[], predicate: (item: T, index: number) => boolean): number {
  for (let i = 0; i < arr.length; i++) {
    if (predicate(arr[i], i)) {
      return i;
    }
  }
  return -1;
}

/**
 * 查找最后一个匹配项的索引
 * @example
 * findLastIndex([1, 2, 1, 2], n => n === 1) // 2
 */
export function findLastIndex<T>(arr: T[], predicate: (item: T, index: number) => boolean): number {
  for (let i = arr.length - 1; i >= 0; i--) {
    if (predicate(arr[i], i)) {
      return i;
    }
  }
  return -1;
}

/**
 * 根据条件查找所有索引
 * @example
 * findAllIndices([1, 2, 1, 2], n => n === 1) // [0, 2]
 */
export function findAllIndices<T>(arr: T[], predicate: (item: T, index: number) => boolean): number[] {
  const indices: number[] = [];
  arr.forEach((item, index) => {
    if (predicate(item, index)) {
      indices.push(index);
    }
  });
  return indices;
}

// ============================================================================
// 排序
// ============================================================================

/**
 * 按字段排序
 * @example
 * sortBy([{ name: 'b' }, { name: 'a' }], 'name') // [{ name: 'a' }, { name: 'b' }]
 * sortBy(users, 'age', 'desc')
 */
export function sortBy<T>(arr: T[], key: keyof T | ((item: T) => number | string), order: "asc" | "desc" = "asc"): T[] {
  const multiplier = order === "asc" ? 1 : -1;

  return [...arr].sort((a, b) => {
    const valueA = typeof key === "function" ? key(a) : a[key];
    const valueB = typeof key === "function" ? key(b) : b[key];

    if (valueA < valueB) return -1 * multiplier;
    if (valueA > valueB) return 1 * multiplier;
    return 0;
  });
}

/**
 * 多字段排序
 * @example
 * sortByMultiple(users, [
 *   { key: 'role', order: 'asc' },
 *   { key: 'name', order: 'asc' }
 * ])
 */
export function sortByMultiple<T>(
  arr: T[],
  sorts: Array<{ key: keyof T | ((item: T) => number | string); order?: "asc" | "desc" }>
): T[] {
  return [...arr].sort((a, b) => {
    for (const { key, order = "asc" } of sorts) {
      const valueA = typeof key === "function" ? key(a) : a[key];
      const valueB = typeof key === "function" ? key(b) : b[key];
      const multiplier = order === "asc" ? 1 : -1;

      if (valueA < valueB) return -1 * multiplier;
      if (valueA > valueB) return 1 * multiplier;
    }
    return 0;
  });
}

// ============================================================================
// 变换
// ============================================================================

/**
 * 打平数组
 * @example
 * flatten([[1, 2], [3, [4, 5]]]) // [1, 2, 3, [4, 5]]
 * flatten([[1, [2, [3]]]], 2) // [1, 2, 3]
 */
export function flatten<T>(arr: unknown[], depth = 1): T[] {
  if (depth < 1) return arr as T[];

  return arr.reduce<T[]>((result, item) => {
    if (Array.isArray(item) && depth > 0) {
      result.push(...flatten<T>(item, depth - 1));
    } else {
      result.push(item as T);
    }
    return result;
  }, []);
}

/**
 * 数组洗牌（随机排序）
 * @example
 * shuffle([1, 2, 3, 4, 5]) // [3, 1, 5, 2, 4] (随机)
 */
export function shuffle<T>(arr: T[]): T[] {
  const result = [...arr];
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [result[i], result[j]] = [result[j], result[i]];
  }
  return result;
}

/**
 * 随机取样
 * @example
 * sample([1, 2, 3, 4, 5]) // 3 (随机一个)
 * sample([1, 2, 3, 4, 5], 2) // [1, 4] (随机两个)
 */
export function sample<T>(arr: T[], count?: number): T | T[] {
  if (count === undefined) {
    return arr[Math.floor(Math.random() * arr.length)];
  }

  const shuffled = shuffle(arr);
  return shuffled.slice(0, Math.min(count, arr.length));
}

// ============================================================================
// 移动与交换
// ============================================================================

/**
 * 移动元素
 * @example
 * move([1, 2, 3, 4], 0, 2) // [2, 3, 1, 4]
 */
export function move<T>(arr: T[], fromIndex: number, toIndex: number): T[] {
  const result = [...arr];
  const [item] = result.splice(fromIndex, 1);
  result.splice(toIndex, 0, item);
  return result;
}

/**
 * 交换元素
 * @example
 * swap([1, 2, 3, 4], 0, 2) // [3, 2, 1, 4]
 */
export function swap<T>(arr: T[], indexA: number, indexB: number): T[] {
  const result = [...arr];
  [result[indexA], result[indexB]] = [result[indexB], result[indexA]];
  return result;
}

// ============================================================================
// 聚合
// ============================================================================

/**
 * 求和
 * @example
 * sum([1, 2, 3]) // 6
 * sum([{ value: 1 }, { value: 2 }], item => item.value) // 3
 */
export function sum<T>(arr: T[], getter?: (item: T) => number): number {
  if (!getter) {
    return (arr as unknown as number[]).reduce((acc, val) => acc + val, 0);
  }
  return arr.reduce((acc, item) => acc + getter(item), 0);
}

/**
 * 平均值
 * @example
 * average([1, 2, 3]) // 2
 */
export function average<T>(arr: T[], getter?: (item: T) => number): number {
  if (arr.length === 0) return 0;
  return sum(arr, getter) / arr.length;
}

/**
 * 最大值项
 * @example
 * maxBy([{ score: 1 }, { score: 3 }], 'score') // { score: 3 }
 */
export function maxBy<T>(arr: T[], key: keyof T | ((item: T) => number)): T | undefined {
  if (arr.length === 0) return undefined;

  return arr.reduce((max, item) => {
    const maxValue = typeof key === "function" ? key(max) : (max[key] as number);
    const itemValue = typeof key === "function" ? key(item) : (item[key] as number);
    return itemValue > maxValue ? item : max;
  });
}

/**
 * 最小值项
 * @example
 * minBy([{ score: 1 }, { score: 3 }], 'score') // { score: 1 }
 */
export function minBy<T>(arr: T[], key: keyof T | ((item: T) => number)): T | undefined {
  if (arr.length === 0) return undefined;

  return arr.reduce((min, item) => {
    const minValue = typeof key === "function" ? key(min) : (min[key] as number);
    const itemValue = typeof key === "function" ? key(item) : (item[key] as number);
    return itemValue < minValue ? item : min;
  });
}

// ============================================================================
// 树形结构
// ============================================================================

export interface TreeNode<T> {
  children?: TreeNode<T>[];
  [key: string]: unknown;
}

/**
 * 数组转树形结构
 * @example
 * const items = [
 *   { id: 1, parentId: null },
 *   { id: 2, parentId: 1 },
 *   { id: 3, parentId: 1 }
 * ]
 * arrayToTree(items, { idKey: 'id', parentKey: 'parentId' })
 */
export function arrayToTree<T extends Record<string, unknown>>(
  items: T[],
  options: {
    idKey?: string;
    parentKey?: string;
    childrenKey?: string;
    rootValue?: unknown;
  } = {}
): (T & { children: T[] })[] {
  const { idKey = "id", parentKey = "parentId", childrenKey = "children", rootValue = null } = options;

  const itemMap = new Map<unknown, T & { children: T[] }>();
  const roots: (T & { children: T[] })[] = [];

  // 创建所有节点的副本
  for (const item of items) {
    itemMap.set(item[idKey], { ...item, [childrenKey]: [] } as T & { children: T[] });
  }

  // 构建树
  for (const item of items) {
    const node = itemMap.get(item[idKey])!;
    const parentId = item[parentKey];

    if (parentId === rootValue) {
      roots.push(node);
    } else {
      const parent = itemMap.get(parentId);
      if (parent) {
        (parent as Record<string, T[]>)[childrenKey].push(node);
      }
    }
  }

  return roots;
}

/**
 * 树形结构转数组
 * @example
 * treeToArray(tree, 'children')
 */
export function treeToArray<T extends TreeNode<T>>(tree: T[], childrenKey: string = "children"): T[] {
  const result: T[] = [];

  const traverse = (nodes: T[]) => {
    for (const node of nodes) {
      result.push(node);
      const children = node[childrenKey] as T[] | undefined;
      if (children && children.length > 0) {
        traverse(children);
      }
    }
  };

  traverse(tree);
  return result;
}

/**
 * 在树中查找节点
 * @example
 * findInTree(tree, node => node.id === 5)
 */
export function findInTree<T extends TreeNode<T>>(
  tree: T[],
  predicate: (node: T) => boolean,
  childrenKey: string = "children"
): T | undefined {
  for (const node of tree) {
    if (predicate(node)) {
      return node;
    }

    const children = node[childrenKey] as T[] | undefined;
    if (children && children.length > 0) {
      const found = findInTree(children, predicate, childrenKey);
      if (found) {
        return found;
      }
    }
  }

  return undefined;
}

/**
 * 过滤树形结构
 * @example
 * filterTree(tree, node => node.visible)
 */
export function filterTree<T extends TreeNode<T>>(
  tree: T[],
  predicate: (node: T) => boolean,
  childrenKey: string = "children"
): T[] {
  return tree.filter(predicate).map((node) => {
    const children = node[childrenKey] as T[] | undefined;
    if (children && children.length > 0) {
      return {
        ...node,
        [childrenKey]: filterTree(children, predicate, childrenKey),
      };
    }
    return { ...node };
  });
}
