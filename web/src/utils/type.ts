/**
 * 类型检查和转换工具函数
 * 提供运行时类型检查和类型安全的转换
 */

// ============================================================================
// 基础类型检查
// ============================================================================

/**
 * 检查是否为字符串
 * @example
 * isString('hello') // true
 * isString(123) // false
 */
export function isString(value: unknown): value is string {
  return typeof value === "string";
}

/**
 * 检查是否为数字
 * @example
 * isNumber(123) // true
 * isNumber(NaN) // false
 * isNumber('123') // false
 */
export function isNumber(value: unknown): value is number {
  return typeof value === "number" && !Number.isNaN(value);
}

/**
 * 检查是否为有限数字
 * @example
 * isFiniteNumber(123) // true
 * isFiniteNumber(Infinity) // false
 */
export function isFiniteNumber(value: unknown): value is number {
  return isNumber(value) && Number.isFinite(value);
}

/**
 * 检查是否为整数
 * @example
 * isInteger(123) // true
 * isInteger(123.5) // false
 */
export function isInteger(value: unknown): value is number {
  return Number.isInteger(value);
}

/**
 * 检查是否为布尔值
 * @example
 * isBoolean(true) // true
 * isBoolean(1) // false
 */
export function isBoolean(value: unknown): value is boolean {
  return typeof value === "boolean";
}

/**
 * 检查是否为 null
 * @example
 * isNull(null) // true
 * isNull(undefined) // false
 */
export function isNull(value: unknown): value is null {
  return value === null;
}

/**
 * 检查是否为 undefined
 * @example
 * isUndefined(undefined) // true
 * isUndefined(null) // false
 */
export function isUndefined(value: unknown): value is undefined {
  return value === undefined;
}

/**
 * 检查是否为 null 或 undefined
 * @example
 * isNullish(null) // true
 * isNullish(undefined) // true
 * isNullish(0) // false
 */
export function isNullish(value: unknown): value is null | undefined {
  return value === null || value === undefined;
}

/**
 * 检查是否已定义（非 undefined）
 * @example
 * isDefined(null) // true
 * isDefined(undefined) // false
 */
export function isDefined<T>(value: T): value is Exclude<T, undefined> {
  return value !== undefined;
}

/**
 * 检查是否有值（非 null 且非 undefined）
 * @example
 * hasValue(0) // true
 * hasValue('') // true
 * hasValue(null) // false
 */
export function hasValue<T>(value: T): value is NonNullable<T> {
  return value !== null && value !== undefined;
}

// ============================================================================
// 复杂类型检查
// ============================================================================

/**
 * 检查是否为对象
 * @example
 * isObject({}) // true
 * isObject([]) // false
 * isObject(null) // false
 */
export function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/**
 * 检查是否为普通对象（Plain Object）
 * @example
 * isPlainObject({}) // true
 * isPlainObject(new Date()) // false
 */
export function isPlainObject(value: unknown): value is Record<string, unknown> {
  if (!isObject(value)) return false;
  const proto = Object.getPrototypeOf(value);
  return proto === null || proto === Object.prototype;
}

/**
 * 检查是否为数组
 * @example
 * isArray([1, 2, 3]) // true
 * isArray('123') // false
 */
export function isArray(value: unknown): value is unknown[] {
  return Array.isArray(value);
}

/**
 * 检查是否为指定类型的数组
 * @example
 * isArrayOf([1, 2, 3], isNumber) // true
 * isArrayOf([1, '2', 3], isNumber) // false
 */
export function isArrayOf<T>(value: unknown, guard: (item: unknown) => item is T): value is T[] {
  return isArray(value) && value.every(guard);
}

/**
 * 检查是否为非空数组
 * @example
 * isNonEmptyArray([1]) // true
 * isNonEmptyArray([]) // false
 */
export function isNonEmptyArray<T>(value: T[] | unknown): value is [T, ...T[]] {
  return isArray(value) && value.length > 0;
}

/**
 * 检查是否为函数
 * @example
 * isFunction(() => {}) // true
 * isFunction({}) // false
 */
export function isFunction(value: unknown): value is (...args: unknown[]) => unknown {
  return typeof value === "function";
}

/**
 * 检查是否为 async 函数
 * @example
 * isAsyncFunction(async () => {}) // true
 * isAsyncFunction(() => {}) // false
 */
export function isAsyncFunction(value: unknown): value is (...args: unknown[]) => Promise<unknown> {
  return isFunction(value) && value.constructor.name === "AsyncFunction";
}

/**
 * 检查是否为 Symbol
 * @example
 * isSymbol(Symbol('test')) // true
 */
export function isSymbol(value: unknown): value is symbol {
  return typeof value === "symbol";
}

/**
 * 检查是否为 BigInt
 * @example
 * isBigInt(BigInt(123)) // true
 */
export function isBigInt(value: unknown): value is bigint {
  return typeof value === "bigint";
}

// ============================================================================
// 特殊类型检查
// ============================================================================

/**
 * 检查是否为 Date 对象
 * @example
 * isDate(new Date()) // true
 */
export function isDate(value: unknown): value is Date {
  return value instanceof Date && !isNaN(value.getTime());
}

/**
 * 检查是否为正则表达式
 * @example
 * isRegExp(/test/) // true
 */
export function isRegExp(value: unknown): value is RegExp {
  return value instanceof RegExp;
}

/**
 * 检查是否为 Error 对象
 * @example
 * isError(new Error('test')) // true
 */
export function isError(value: unknown): value is Error {
  return value instanceof Error;
}

/**
 * 检查是否为 Promise
 * @example
 * isPromise(Promise.resolve()) // true
 */
export function isPromise<T = unknown>(value: unknown): value is Promise<T> {
  return (
    value instanceof Promise ||
    (isObject(value) &&
      isFunction((value as { then?: unknown }).then) &&
      isFunction((value as { catch?: unknown }).catch))
  );
}

/**
 * 检查是否为 Map
 * @example
 * isMap(new Map()) // true
 */
export function isMap<K = unknown, V = unknown>(value: unknown): value is Map<K, V> {
  return value instanceof Map;
}

/**
 * 检查是否为 Set
 * @example
 * isSet(new Set()) // true
 */
export function isSet<T = unknown>(value: unknown): value is Set<T> {
  return value instanceof Set;
}

/**
 * 检查是否为 WeakMap
 * @example
 * isWeakMap(new WeakMap()) // true
 */
export function isWeakMap<K extends object = object, V = unknown>(value: unknown): value is WeakMap<K, V> {
  return value instanceof WeakMap;
}

/**
 * 检查是否为 WeakSet
 * @example
 * isWeakSet(new WeakSet()) // true
 */
export function isWeakSet<T extends object = object>(value: unknown): value is WeakSet<T> {
  return value instanceof WeakSet;
}

// ============================================================================
// 浏览器特定类型检查
// ============================================================================

/**
 * 检查是否为 DOM 元素
 * @example
 * isElement(document.body) // true
 */
export function isElement(value: unknown): value is Element {
  return value instanceof Element;
}

/**
 * 检查是否为 HTML 元素
 * @example
 * isHTMLElement(document.body) // true
 */
export function isHTMLElement(value: unknown): value is HTMLElement {
  return value instanceof HTMLElement;
}

/**
 * 检查是否为 Node
 * @example
 * isNode(document.body) // true
 */
export function isNode(value: unknown): value is Node {
  return value instanceof Node;
}

/**
 * 检查是否为 Blob
 * @example
 * isBlob(new Blob()) // true
 */
export function isBlob(value: unknown): value is Blob {
  return value instanceof Blob;
}

/**
 * 检查是否为 File
 * @example
 * isFile(file) // true
 */
export function isFile(value: unknown): value is File {
  return value instanceof File;
}

/**
 * 检查是否为 FormData
 * @example
 * isFormData(new FormData()) // true
 */
export function isFormData(value: unknown): value is FormData {
  return value instanceof FormData;
}

/**
 * 检查是否为 ArrayBuffer
 * @example
 * isArrayBuffer(new ArrayBuffer(8)) // true
 */
export function isArrayBuffer(value: unknown): value is ArrayBuffer {
  return value instanceof ArrayBuffer;
}

// ============================================================================
// 字符串类型检查
// ============================================================================

/**
 * 检查是否为空字符串
 * @example
 * isEmptyString('') // true
 * isEmptyString('  ') // false
 */
export function isEmptyString(value: unknown): value is "" {
  return value === "";
}

/**
 * 检查是否为空白字符串
 * @example
 * isBlankString('  ') // true
 * isBlankString('') // true
 */
export function isBlankString(value: unknown): boolean {
  return isString(value) && value.trim() === "";
}

/**
 * 检查是否为非空字符串
 * @example
 * isNonEmptyString('hello') // true
 * isNonEmptyString('') // false
 */
export function isNonEmptyString(value: unknown): value is string {
  return isString(value) && value.length > 0;
}

/**
 * 检查是否为非空白字符串
 * @example
 * isNonBlankString('hello') // true
 * isNonBlankString('  ') // false
 */
export function isNonBlankString(value: unknown): value is string {
  return isString(value) && value.trim().length > 0;
}

// ============================================================================
// 类型转换
// ============================================================================

/**
 * 安全转换为字符串
 * @example
 * toString(123) // '123'
 * toString(null) // ''
 */
export function toString(value: unknown): string {
  if (isNullish(value)) return "";
  if (isString(value)) return value;
  return String(value);
}

/**
 * 安全转换为数字
 * @example
 * toNumber('123') // 123
 * toNumber('abc') // 0
 * toNumber('abc', -1) // -1
 */
export function toNumber(value: unknown, defaultValue: number = 0): number {
  if (isNumber(value)) return value;
  if (isString(value)) {
    const num = parseFloat(value);
    return isNaN(num) ? defaultValue : num;
  }
  return defaultValue;
}

/**
 * 安全转换为整数
 * @example
 * toInteger('123.5') // 123
 * toInteger('abc') // 0
 */
export function toInteger(value: unknown, defaultValue: number = 0): number {
  const num = toNumber(value, defaultValue);
  return Math.trunc(num);
}

/**
 * 安全转换为布尔值
 * @example
 * toBoolean('true') // true
 * toBoolean('false') // false
 * toBoolean(1) // true
 * toBoolean(0) // false
 */
export function toBoolean(value: unknown): boolean {
  if (isBoolean(value)) return value;
  if (isString(value)) {
    const lower = value.toLowerCase().trim();
    return lower === "true" || lower === "1" || lower === "yes";
  }
  if (isNumber(value)) return value !== 0;
  return Boolean(value);
}

/**
 * 安全转换为数组
 * @example
 * toArray('hello') // ['hello']
 * toArray([1, 2]) // [1, 2]
 * toArray(null) // []
 */
export function toArray<T>(value: T | T[] | null | undefined): T[] {
  if (isNullish(value)) return [];
  if (isArray(value)) return value;
  return [value] as T[];
}

/**
 * 安全转换为日期
 * @example
 * toDate('2024-01-01') // Date
 * toDate(1704067200000) // Date
 * toDate('invalid') // null
 */
export function toDate(value: unknown): Date | null {
  if (value instanceof Date) return value;
  if (isString(value) || isNumber(value)) {
    const date = new Date(value);
    return isNaN(date.getTime()) ? null : date;
  }
  return null;
}

// ============================================================================
// 类型断言
// ============================================================================

/**
 * 断言值不为 null 或 undefined
 * @example
 * assertDefined(value, 'Value is required')
 */
export function assertDefined<T>(
  value: T,
  message: string = "Value is null or undefined"
): asserts value is NonNullable<T> {
  if (isNullish(value)) {
    throw new Error(message);
  }
}

/**
 * 断言值为 true
 * @example
 * assert(condition, 'Condition must be true')
 */
export function assert(condition: unknown, message: string = "Assertion failed"): asserts condition {
  if (!condition) {
    throw new Error(message);
  }
}

/**
 * 永远不会执行的断言（用于穷尽检查）
 * @example
 * switch (type) {
 *   case 'a': return handleA();
 *   case 'b': return handleB();
 *   default: assertNever(type);
 * }
 */
export function assertNever(value: never, message?: string): never {
  throw new Error(message ?? `Unexpected value: ${value}`);
}

// ============================================================================
// 类型守卫组合
// ============================================================================

/**
 * 创建联合类型守卫
 * @example
 * const isStringOrNumber = unionGuard(isString, isNumber);
 * isStringOrNumber('hello') // true
 * isStringOrNumber(123) // true
 */
export function unionGuard<T extends (value: unknown) => boolean>(...guards: T[]): (value: unknown) => boolean {
  return (value: unknown) => guards.some((guard) => guard(value));
}

/**
 * 创建交叉类型守卫
 * @example
 * const isNonEmptyStringArray = intersectionGuard(isArray, isNonEmptyArray);
 */
export function intersectionGuard<T extends (value: unknown) => boolean>(...guards: T[]): (value: unknown) => boolean {
  return (value: unknown) => guards.every((guard) => guard(value));
}

/**
 * 创建否定类型守卫
 * @example
 * const isNotNull = notGuard(isNull);
 */
export function notGuard<T extends (value: unknown) => boolean>(guard: T): (value: unknown) => boolean {
  return (value: unknown) => !guard(value);
}

// ============================================================================
// 实用类型工具
// ============================================================================

/**
 * 获取值的类型名称
 * @example
 * getType(null) // 'null'
 * getType([]) // 'array'
 * getType({}) // 'object'
 * getType(new Date()) // 'date'
 */
export function getType(value: unknown): string {
  if (value === null) return "null";
  if (value === undefined) return "undefined";
  if (Array.isArray(value)) return "array";
  if (value instanceof Date) return "date";
  if (value instanceof RegExp) return "regexp";
  if (value instanceof Error) return "error";
  if (value instanceof Map) return "map";
  if (value instanceof Set) return "set";
  if (value instanceof Promise) return "promise";
  return typeof value;
}

/**
 * 检查两个值是否类型相同
 * @example
 * isSameType('a', 'b') // true
 * isSameType('a', 1) // false
 */
export function isSameType(a: unknown, b: unknown): boolean {
  return getType(a) === getType(b);
}

/**
 * 检查值是否为原始类型
 * @example
 * isPrimitive('hello') // true
 * isPrimitive({}) // false
 */
export function isPrimitive(value: unknown): value is string | number | boolean | null | undefined | symbol | bigint {
  return (
    value === null ||
    value === undefined ||
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean" ||
    typeof value === "symbol" ||
    typeof value === "bigint"
  );
}

/**
 * 检查值是否可迭代
 * @example
 * isIterable([1, 2, 3]) // true
 * isIterable('hello') // true
 * isIterable(123) // false
 */
export function isIterable(value: unknown): value is Iterable<unknown> {
  return (
    value !== null &&
    value !== undefined &&
    typeof (value as { [Symbol.iterator]?: unknown })[Symbol.iterator] === "function"
  );
}

/**
 * 检查值是否可序列化为 JSON
 * @example
 * isJsonSerializable({ a: 1 }) // true
 * isJsonSerializable(() => {}) // false
 */
export function isJsonSerializable(value: unknown): boolean {
  try {
    JSON.stringify(value);
    return true;
  } catch {
    return false;
  }
}
