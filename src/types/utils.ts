/**
 * TypeScript 工具类型
 * 提供常用的类型转换和工具函数
 */

/**
 * 将类型 T 的所有属性变为可空
 * @example Nullable<string> => string | null
 */
export type Nullable<T> = T | null;

/**
 * 将类型 T 的所有属性变为可选
 * @example Optional<string> => string | undefined
 */
export type Optional<T> = T | undefined;

/**
 * 将类型 T 的所有属性（包括嵌套属性）变为可选
 * @example DeepPartial<{ a: { b: string } }> => { a?: { b?: string } }
 */
export type DeepPartial<T> = T extends object
  ? {
      [P in keyof T]?: DeepPartial<T[P]>;
    }
  : T;

/**
 * 将类型 T 的所有属性（包括嵌套属性）变为只读
 * @example DeepReadonly<{ a: { b: string } }> => { readonly a: { readonly b: string } }
 */
export type DeepReadonly<T> = T extends object
  ? {
      readonly [P in keyof T]: DeepReadonly<T[P]>;
    }
  : T;

/**
 * 提取 Promise 类型的返回值
 * @example Awaited<Promise<string>> => string
 */
export type Awaited<T> = T extends Promise<infer U> ? U : T;

/**
 * 提取异步函数的返回值类型
 * @example AsyncReturnType<() => Promise<string>> => string
 */
export type AsyncReturnType<T extends (...args: unknown[]) => Promise<unknown>> = Awaited<ReturnType<T>>;

/**
 * 提取数组元素类型
 * @example ArrayElement<string[]> => string
 */
export type ArrayElement<T> = T extends readonly (infer E)[] ? E : never;

/**
 * 从联合类型中排除 null 和 undefined
 * @example NonNullable<string | null | undefined> => string
 */
export type NonNullableDeep<T> = T extends object
  ? {
      [P in keyof T]-?: NonNullableDeep<NonNullable<T[P]>>;
    }
  : NonNullable<T>;

/**
 * 选择对象中指定的键
 * @example Pick<{ a: 1, b: 2 }, 'a'> => { a: 1 }
 */
export type PickByValue<T, V> = {
  [K in keyof T as T[K] extends V ? K : never]: T[K];
};

/**
 * 排除对象中指定值类型的键
 * @example OmitByValue<{ a: string, b: number }, number> => { a: string }
 */
export type OmitByValue<T, V> = {
  [K in keyof T as T[K] extends V ? never : K]: T[K];
};

/**
 * 将对象的所有属性变为必需
 * @example RequiredKeys<{ a?: string, b: number }> => 'b'
 */
export type RequiredKeys<T> = {
  [K in keyof T]-?: undefined extends T[K] ? never : K;
}[keyof T];

/**
 * 提取对象的可选键
 * @example OptionalKeys<{ a?: string, b: number }> => 'a'
 */
export type OptionalKeys<T> = {
  [K in keyof T]-?: undefined extends T[K] ? K : never;
}[keyof T];

/**
 * 函数参数类型
 * @example FunctionArgs<(a: string, b: number) => void> => [string, number]
 */
export type FunctionArgs<T extends (...args: unknown[]) => unknown> = T extends (...args: infer A) => unknown ? A : never;

/**
 * 字符串字面量联合类型
 * @example StringLiteral<'a' | 'b'> => 'a' | 'b'
 */
export type StringLiteral<T> = T extends string ? (string extends T ? never : T) : never;

/**
 * 将两个类型合并，后者覆盖前者
 * @example Merge<{ a: string }, { a: number, b: string }> => { a: number, b: string }
 */
export type Merge<T, U> = Omit<T, keyof U> & U;

/**
 * 可变数组类型（移除 readonly）
 * @example Mutable<readonly string[]> => string[]
 */
export type Mutable<T> = T extends readonly (infer U)[] ? U[] : { -readonly [K in keyof T]: T[K] };
