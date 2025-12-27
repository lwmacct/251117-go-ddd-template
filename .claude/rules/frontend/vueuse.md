---
paths:
  - "src/**/*.vue"
  - "src/**/*.ts"
---

# VueUse 使用规范

优先使用 `@vueuse/core` 处理浏览器 API，确保自动生命周期管理和响应式。

## 必须使用 VueUse

| 场景     | ❌ 禁止                    | ✅ 推荐                        |
| -------- | -------------------------- | ------------------------------ |
| 定时器   | `setInterval`/`setTimeout` | `useIntervalFn`/`useTimeoutFn` |
| 本地存储 | `localStorage`             | `useLocalStorage`              |
| 剪贴板   | `navigator.clipboard`      | `useClipboard`                 |
| 事件监听 | `addEventListener`         | `useEventListener`             |
| 拖拽上传 | 手动 drag 事件             | `useDropZone`                  |
| 文件选择 | `<input type="file">`      | `useFileDialog`                |
| 防抖节流 | 手动实现                   | `refDebounced`/`useThrottleFn` |
| 系统主题 | `matchMedia`               | `usePreferredDark`             |
