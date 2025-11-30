/**
 * 对象工具函数
 * 提供常用的对象操作方法
 */

// ============================================================================
// 类型定义
// ============================================================================

type DeepPartial<T> = T extends object ? { [P in keyof T]?: DeepPartial<T[P]> } : T;

type Path = string | (string | number)[];

// ============================================================================
// 深度操作
// ============================================================================

/**
 * 深拷贝
 * @example
 * const copy = deepClone({ a: { b: 1 } })
 */
export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== "object") {
    return obj;
  }

  if (obj instanceof Date) {
    return new Date(obj.getTime()) as T;
  }

  if (obj instanceof RegExp) {
    return new RegExp(obj.source, obj.flags) as T;
  }

  if (Array.isArray(obj)) {
    return obj.map((item) => deepClone(item)) as T;
  }

  if (obj instanceof Map) {
    const map = new Map();
    obj.forEach((value, key) => {
      map.set(deepClone(key), deepClone(value));
    });
    return map as T;
  }

  if (obj instanceof Set) {
    const set = new Set();
    obj.forEach((value) => {
      set.add(deepClone(value));
    });
    return set as T;
  }

  const result = {} as T;
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      result[key] = deepClone(obj[key]);
    }
  }

  return result;
}

/**
 * 深度合并
 * @example
 * deepMerge({ a: { b: 1 } }, { a: { c: 2 } }) // { a: { b: 1, c: 2 } }
 */
export function deepMerge<T extends object>(target: T, ...sources: DeepPartial<T>[]): T {
  if (!sources.length) return target;

  const source = sources.shift();
  if (source === undefined) return target;

  for (const key in source) {
    if (Object.prototype.hasOwnProperty.call(source, key)) {
      const sourceValue = source[key];
      const targetValue = (target as Record<string, unknown>)[key];

      if (isObject(sourceValue) && isObject(targetValue)) {
        (target as Record<string, unknown>)[key] = deepMerge(targetValue as object, sourceValue as object);
      } else {
        (target as Record<string, unknown>)[key] = deepClone(sourceValue);
      }
    }
  }

  return deepMerge(target, ...sources);
}

/**
 * 深度比较
 * @example
 * deepEqual({ a: 1 }, { a: 1 }) // true
 */
export function deepEqual(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  if (a === null || b === null) return false;
  if (typeof a !== typeof b) return false;

  if (a instanceof Date && b instanceof Date) {
    return a.getTime() === b.getTime();
  }

  if (a instanceof RegExp && b instanceof RegExp) {
    return a.toString() === b.toString();
  }

  if (Array.isArray(a) && Array.isArray(b)) {
    if (a.length !== b.length) return false;
    return a.every((item, index) => deepEqual(item, b[index]));
  }

  if (typeof a === "object" && typeof b === "object") {
    const keysA = Object.keys(a as object);
    const keysB = Object.keys(b as object);

    if (keysA.length !== keysB.length) return false;

    return keysA.every((key) => deepEqual((a as Record<string, unknown>)[key], (b as Record<string, unknown>)[key]));
  }

  return false;
}

// ============================================================================
// 路径操作
// ============================================================================

/**
 * 解析路径字符串
 * @example
 * parsePath('a.b.c') // ['a', 'b', 'c']
 * parsePath('a[0].b') // ['a', 0, 'b']
 */
export function parsePath(path: Path): (string | number)[] {
  if (Array.isArray(path)) {
    return path;
  }

  const result: (string | number)[] = [];
  const regex = /([^.[\]]+)|\[(\d+)\]/g;
  let match;

  while ((match = regex.exec(path)) !== null) {
    if (match[1]) {
      result.push(match[1]);
    } else if (match[2]) {
      result.push(parseInt(match[2], 10));
    }
  }

  return result;
}

/**
 * 获取嵌套属性值
 * @example
 * get({ a: { b: 1 } }, 'a.b') // 1
 * get({ a: [1, 2] }, 'a[1]') // 2
 * get({ a: 1 }, 'b.c', 'default') // 'default'
 */
export function get<T = unknown>(obj: unknown, path: Path, defaultValue?: T): T | undefined {
  const keys = parsePath(path);
  let result: unknown = obj;

  for (const key of keys) {
    if (result === null || result === undefined) {
      return defaultValue;
    }
    result = (result as Record<string | number, unknown>)[key];
  }

  return (result === undefined ? defaultValue : result) as T;
}

/**
 * 设置嵌套属性值
 * @example
 * set({}, 'a.b.c', 1) // { a: { b: { c: 1 } } }
 * set({}, 'a[0].b', 1) // { a: [{ b: 1 }] }
 */
export function set<T extends object>(obj: T, path: Path, value: unknown): T {
  const keys = parsePath(path);
  const result = deepClone(obj);
  let current: Record<string | number, unknown> = result as Record<string | number, unknown>;

  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i];
    const nextKey = keys[i + 1];

    if (current[key] === undefined || current[key] === null) {
      current[key] = typeof nextKey === "number" ? [] : {};
    }

    current = current[key] as Record<string | number, unknown>;
  }

  current[keys[keys.length - 1]] = value;
  return result;
}

/**
 * 检查是否存在嵌套属性
 * @example
 * has({ a: { b: 1 } }, 'a.b') // true
 * has({ a: { b: 1 } }, 'a.c') // false
 */
export function has(obj: unknown, path: Path): boolean {
  const keys = parsePath(path);
  let current: unknown = obj;

  for (const key of keys) {
    if (current === null || current === undefined) {
      return false;
    }
    if (!Object.prototype.hasOwnProperty.call(current, key)) {
      return false;
    }
    current = (current as Record<string | number, unknown>)[key];
  }

  return true;
}

/**
 * 删除嵌套属性
 * @example
 * unset({ a: { b: 1, c: 2 } }, 'a.b') // { a: { c: 2 } }
 */
export function unset<T extends object>(obj: T, path: Path): T {
  const keys = parsePath(path);
  const result = deepClone(obj);

  if (keys.length === 0) return result;

  let current: Record<string | number, unknown> = result as Record<string | number, unknown>;

  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i];
    if (current[key] === undefined || current[key] === null) {
      return result;
    }
    current = current[key] as Record<string | number, unknown>;
  }

  delete current[keys[keys.length - 1]];
  return result;
}

// ============================================================================
// 筛选与转换
// ============================================================================

/**
 * 选取指定键
 * @example
 * pick({ a: 1, b: 2, c: 3 }, ['a', 'b']) // { a: 1, b: 2 }
 */
export function pick<T extends object, K extends keyof T>(obj: T, keys: K[]): Pick<T, K> {
  const result = {} as Pick<T, K>;
  for (const key of keys) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      result[key] = obj[key];
    }
  }
  return result;
}

/**
 * 排除指定键
 * @example
 * omit({ a: 1, b: 2, c: 3 }, ['c']) // { a: 1, b: 2 }
 */
export function omit<T extends object, K extends keyof T>(obj: T, keys: K[]): Omit<T, K> {
  const result = { ...obj };
  for (const key of keys) {
    delete result[key];
  }
  return result as Omit<T, K>;
}

/**
 * 按条件筛选键值对
 * @example
 * filterObject({ a: 1, b: null, c: 3 }, v => v !== null)
 * // { a: 1, c: 3 }
 */
export function filterObject<T extends object>(
  obj: T,
  predicate: (value: T[keyof T], key: keyof T) => boolean
): Partial<T> {
  const result: Partial<T> = {};
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key) && predicate(obj[key], key)) {
      result[key] = obj[key];
    }
  }
  return result;
}

/**
 * 转换键值
 * @example
 * mapObject({ a: 1, b: 2 }, v => v * 2) // { a: 2, b: 4 }
 */
export function mapObject<T extends object, U>(
  obj: T,
  fn: (value: T[keyof T], key: keyof T) => U
): { [K in keyof T]: U } {
  const result = {} as { [K in keyof T]: U };
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      result[key] = fn(obj[key], key);
    }
  }
  return result;
}

/**
 * 键值反转
 * @example
 * invert({ a: 'x', b: 'y' }) // { x: 'a', y: 'b' }
 */
export function invert<T extends Record<string, string | number>>(obj: T): Record<string, keyof T> {
  const result: Record<string, keyof T> = {};
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      result[String(obj[key])] = key;
    }
  }
  return result;
}

// ============================================================================
// 类型检查
// ============================================================================

/**
 * 检查是否为对象
 */
export function isObject(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

/**
 * 检查是否为空对象
 * @example
 * isEmpty({}) // true
 * isEmpty({ a: 1 }) // false
 */
export function isEmpty(obj: object): boolean {
  return Object.keys(obj).length === 0;
}

/**
 * 检查是否为普通对象
 */
export function isPlainObject(value: unknown): value is Record<string, unknown> {
  if (!isObject(value)) return false;
  const proto = Object.getPrototypeOf(value);
  return proto === null || proto === Object.prototype;
}

// ============================================================================
// 遍历
// ============================================================================

/**
 * 遍历对象
 * @example
 * forEachObject({ a: 1, b: 2 }, (value, key) => {
 *   console.log(key, value)
 * })
 */
export function forEachObject<T extends object>(obj: T, fn: (value: T[keyof T], key: keyof T) => void): void {
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      fn(obj[key], key);
    }
  }
}

/**
 * 深度遍历对象
 * @example
 * deepForEach({ a: { b: 1 } }, (value, path) => {
 *   console.log(path.join('.'), value)
 * })
 */
export function deepForEach(obj: unknown, fn: (value: unknown, path: string[]) => void, path: string[] = []): void {
  if (isObject(obj)) {
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        const newPath = [...path, key];
        const value = obj[key];
        fn(value, newPath);

        if (isObject(value)) {
          deepForEach(value, fn, newPath);
        }
      }
    }
  }
}

// ============================================================================
// 差异比较
// ============================================================================

export interface ObjectDiff {
  added: Record<string, unknown>;
  removed: Record<string, unknown>;
  changed: Record<string, { from: unknown; to: unknown }>;
}

/**
 * 比较两个对象的差异
 * @example
 * diff({ a: 1, b: 2 }, { a: 1, c: 3 })
 * // { added: { c: 3 }, removed: { b: 2 }, changed: {} }
 */
export function diff<T extends object>(prev: T, next: T): ObjectDiff {
  const result: ObjectDiff = {
    added: {},
    removed: {},
    changed: {},
  };

  const prevKeys = new Set(Object.keys(prev));
  const nextKeys = new Set(Object.keys(next));

  // 查找新增
  for (const key of nextKeys) {
    if (!prevKeys.has(key)) {
      result.added[key] = (next as Record<string, unknown>)[key];
    }
  }

  // 查找删除
  for (const key of prevKeys) {
    if (!nextKeys.has(key)) {
      result.removed[key] = (prev as Record<string, unknown>)[key];
    }
  }

  // 查找变化
  for (const key of prevKeys) {
    if (nextKeys.has(key)) {
      const prevValue = (prev as Record<string, unknown>)[key];
      const nextValue = (next as Record<string, unknown>)[key];

      if (!deepEqual(prevValue, nextValue)) {
        result.changed[key] = { from: prevValue, to: nextValue };
      }
    }
  }

  return result;
}

// ============================================================================
// 其他工具
// ============================================================================

/**
 * 移除对象中的 null/undefined 值
 * @example
 * compact({ a: 1, b: null, c: undefined }) // { a: 1 }
 */
export function compact<T extends object>(obj: T): Partial<T> {
  return filterObject(obj, (value) => value != null);
}

/**
 * 创建带默认值的对象
 * @example
 * defaults({ a: 1 }, { a: 0, b: 2 }) // { a: 1, b: 2 }
 */
export function defaults<T extends object>(obj: Partial<T>, ...defaultsObjects: Partial<T>[]): T {
  const result = { ...obj } as T;

  for (const defaults of defaultsObjects) {
    for (const key in defaults) {
      if (Object.prototype.hasOwnProperty.call(defaults, key)) {
        if ((result as Record<string, unknown>)[key] === undefined) {
          (result as Record<string, unknown>)[key] = defaults[key];
        }
      }
    }
  }

  return result;
}

/**
 * 对象转数组
 * @example
 * entries({ a: 1, b: 2 }) // [['a', 1], ['b', 2]]
 */
export function entries<T extends object>(obj: T): [keyof T, T[keyof T]][] {
  return Object.entries(obj) as [keyof T, T[keyof T]][];
}

/**
 * 数组转对象
 * @example
 * fromEntries([['a', 1], ['b', 2]]) // { a: 1, b: 2 }
 */
export function fromEntries<K extends string | number | symbol, V>(entries: [K, V][]): Record<K, V> {
  return Object.fromEntries(entries) as Record<K, V>;
}
