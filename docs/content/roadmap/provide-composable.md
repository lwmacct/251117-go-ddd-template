# Provide Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:36+4`
- [已实现功能](#已实现功能) `:40+22`
  - [上下文创建](#上下文创建) `:42+7`
  - [特殊上下文](#特殊上下文) `:49+8`
  - [注入工具](#注入工具) `:57+5`
- [使用方式](#使用方式) `:62+254`
  - [基本上下文](#基本上下文) `:64+25`
  - [状态上下文](#状态上下文) `:89+18`
  - [响应式上下文](#响应式上下文) `:107+26`
  - [只读上下文](#只读上下文) `:133+26`
  - [事件总线上下文](#事件总线上下文) `:159+33`
  - [主题上下文](#主题上下文) `:192+22`
  - [国际化上下文](#国际化上下文) `:214+33`
  - [持久化上下文](#持久化上下文) `:247+32`
  - [工厂上下文](#工厂上下文) `:279+21`
  - [注入工具](#注入工具-1) `:300+16`
- [API](#api) `:316+41`
  - [createContext](#createcontext) `:318+13`
  - [createEventBusContext](#createeventbuscontext) `:331+8`
  - [createThemeContext](#createthemecontext) `:339+9`
  - [createI18nContext](#createi18ncontext) `:348+9`
- [代码位置](#代码位置) `:357+7`

<!--TOC-->

## 需求背景

前端需要增强的依赖注入工具函数，支持类型安全、状态管理、事件总线等高级功能。

## 已实现功能

### 上下文创建

- `createContext` - 类型安全的上下文
- `createStateContext` - 状态上下文
- `createReactiveContext` - 响应式上下文
- `createReadonlyContext` - 只读上下文

### 特殊上下文

- `createEventBusContext` - 事件总线上下文
- `createThemeContext` - 主题上下文
- `createI18nContext` - 国际化上下文
- `createStorageContext` - 持久化上下文
- `createFactoryContext` - 工厂上下文

### 注入工具

- `useOptionalInject` - 可选注入
- `useRequiredInject` - 必需注入

## 使用方式

### 基本上下文

```typescript
import { createContext } from "@/composables/useProvide";

// 定义类型
interface User {
  id: number;
  name: string;
}

// 创建上下文
const UserContext = createContext<User>("User");

// 父组件 - 提供
UserContext.provide({ id: 1, name: "John" });

// 子组件 - 注入
const user = UserContext.inject();
console.log(user.name); // 'John'

// 带默认值
const user = UserContext.inject({ id: 0, name: "Guest" });
```

### 状态上下文

```typescript
import { createStateContext } from "@/composables/useProvide";

// 创建可变状态上下文
const CountContext = createStateContext<number>("Count");

// 父组件
const { state: count, setState: setCount } = CountContext.provide(0);

// 子组件
const { state: count, setState: setCount } = CountContext.inject();

// 更新状态
setCount(count.value + 1);
```

### 响应式上下文

```typescript
import { createReactiveContext } from "@/composables/useProvide";

interface AppState {
  user: User | null;
  settings: Settings;
  notifications: Notification[];
}

const AppContext = createReactiveContext<AppState>("App");

// 父组件
const state = AppContext.provide({
  user: null,
  settings: { theme: "light" },
  notifications: [],
});

// 子组件 - 直接修改响应式对象
const state = AppContext.inject();
state.user = { id: 1, name: "John" };
state.notifications.push({ message: "Hello" });
```

### 只读上下文

```typescript
import { createReadonlyContext } from "@/composables/useProvide";

interface Config {
  apiUrl: string;
  timeout: number;
  features: string[];
}

const ConfigContext = createReadonlyContext<Config>("Config");

// 父组件
ConfigContext.provide({
  apiUrl: "/api",
  timeout: 5000,
  features: ["feature-a", "feature-b"],
});

// 子组件 - 只能读取，无法修改
const config = ConfigContext.inject();
console.log(config.apiUrl); // '/api'
// config.apiUrl = '/new-api' // TypeScript 错误
```

### 事件总线上下文

```typescript
import { createEventBusContext } from "@/composables/useProvide";

// 定义事件类型
interface AppEvents {
  "user:login": { id: number; name: string };
  "user:logout": void;
  "notification:show": { message: string; type: "success" | "error" };
}

const EventBus = createEventBusContext<AppEvents>("EventBus");

// 父组件
EventBus.provide();

// 子组件 A - 发送事件
const bus = EventBus.inject();
bus.emit("user:login", { id: 1, name: "John" });

// 子组件 B - 监听事件
const bus = EventBus.inject();
const unsubscribe = bus.on("user:login", (user) => {
  console.log("User logged in:", user.name);
});

// 取消订阅
unsubscribe();
// 或
bus.off("user:login", handler);
```

### 主题上下文

```typescript
import { createThemeContext } from "@/composables/useProvide";

const ThemeProvider = createThemeContext("Theme");

// 父组件
const theme = ThemeProvider.provide("light");

// 子组件
const { theme, isDark, setTheme, toggleDark } = ThemeProvider.inject();

// 使用
console.log(isDark.value); // false
setTheme("dark");
console.log(isDark.value); // true

// 切换
toggleDark(); // 切换暗色/亮色模式
```

### 国际化上下文

```typescript
import { createI18nContext } from "@/composables/useProvide";

const messages = {
  en: {
    greeting: "Hello, {name}!",
    welcome: "Welcome to our app",
  },
  zh: {
    greeting: "你好, {name}!",
    welcome: "欢迎使用我们的应用",
  },
};

const I18n = createI18nContext("I18n", messages);

// 父组件
const i18n = I18n.provide("zh");

// 子组件
const { t, locale, setLocale, availableLocales } = I18n.inject();

console.log(t("greeting", { name: "World" }));
// 输出: 你好, World!

// 切换语言
setLocale("en");
console.log(t("greeting", { name: "World" }));
// 输出: Hello, World!
```

### 持久化上下文

```typescript
import { createStorageContext } from "@/composables/useProvide";

interface UserSettings {
  theme: string;
  fontSize: number;
  notifications: boolean;
}

const SettingsContext = createStorageContext<UserSettings>("Settings", {
  storageKey: "app-settings",
  storage: localStorage,
  defaultValue: { theme: "light", fontSize: 14, notifications: true },
});

// 父组件
const settings = SettingsContext.provide({
  theme: "light",
  fontSize: 14,
  notifications: true,
});

// 子组件
const settings = SettingsContext.inject();

// 修改会自动保存到 localStorage
settings.theme = "dark";
settings.fontSize = 16;
```

### 工厂上下文

```typescript
import { createFactoryContext } from "@/composables/useProvide";

// 创建 Logger 工厂
const LoggerContext = createFactoryContext("Logger", (name: string) => ({
  log: (msg: string) => console.log(`[${name}] ${msg}`),
  warn: (msg: string) => console.warn(`[${name}] ${msg}`),
  error: (msg: string) => console.error(`[${name}] ${msg}`),
}));

// 父组件
LoggerContext.provide();

// 子组件 - 每次调用创建新实例
const createLogger = LoggerContext.inject();
const logger = createLogger("UserService");
logger.log("User created");
```

### 注入工具

```typescript
import { useOptionalInject, useRequiredInject } from "@/composables/useProvide";

// 可选注入 - 可能返回 undefined
const user = useOptionalInject(UserKey);
if (user) {
  console.log(user.name);
}

// 必需注入 - 不存在则报错
const config = useRequiredInject(ConfigKey, "Config");
// 如果 Config 未提供，会抛出错误：[Config] 必须在提供者组件内使用
```

## API

### createContext

| 参数         | 类型   | 说明           |
| ------------ | ------ | -------------- |
| name         | string | 上下文名称     |
| defaultValue | T      | 默认值（可选） |

| 返回值  | 类型                 | 说明     |
| ------- | -------------------- | -------- |
| provide | `(value: T) => void` | 提供函数 |
| inject  | `(default?: T) => T` | 注入函数 |
| key     | InjectionKey         | 注入键   |

### createEventBusContext

| 返回值方法 | 类型                              | 说明     |
| ---------- | --------------------------------- | -------- |
| emit       | `(event, payload) => void`        | 发送事件 |
| on         | `(event, handler) => unsubscribe` | 订阅事件 |
| off        | `(event, handler) => void`        | 取消订阅 |

### createThemeContext

| 返回值     | 类型                   | 说明         |
| ---------- | ---------------------- | ------------ |
| theme      | Ref\<string\>          | 当前主题     |
| isDark     | ComputedRef\<boolean\> | 是否暗色     |
| setTheme   | `(theme) => void`      | 设置主题     |
| toggleDark | `() => void`           | 切换暗色模式 |

### createI18nContext

| 返回值           | 类型                       | 说明     |
| ---------------- | -------------------------- | -------- |
| locale           | Ref\<string\>              | 当前语言 |
| availableLocales | string[]                   | 可用语言 |
| setLocale        | `(locale) => void`         | 设置语言 |
| t                | `(key, params?) => string` | 翻译函数 |

## 代码位置

```
web/src/
└── composables/
    └── useProvide.ts    # Provide Composable
```
