# 键盘快捷键

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

## 需求背景

为提升效率用户体验，需要支持键盘快捷键操作，如快速保存、关闭对话框等。

## 已实现功能

### useKeyboard Composable

- 支持组合键（Ctrl、Shift、Alt、Meta）
- 自动跳过输入框中的快捷键
- 支持条件触发
- 提供常用快捷键预设

## 使用方式

### 基础用法

```typescript
import { useKeyboard } from "@/composables/useKeyboard";

useKeyboard([
  {
    key: "ctrl+s",
    handler: () => handleSave(),
    description: "保存",
  },
  {
    key: "escape",
    handler: () => closeDialog(),
    description: "关闭",
  },
]);
```

### 单个快捷键

```typescript
import { useShortcut } from "@/composables/useKeyboard";

useShortcut("ctrl+k", () => openSearch());
```

### 条件触发

```typescript
useKeyboard([
  {
    key: "delete",
    handler: () => deleteSelected(),
    when: () => hasSelection.value,
  },
]);
```

## 常用快捷键预设

```typescript
import { commonShortcuts } from "@/composables/useKeyboard";

// 可用预设
commonShortcuts.save     // "ctrl+s"
commonShortcuts.new      // "ctrl+n"
commonShortcuts.search   // "ctrl+k"
commonShortcuts.escape   // "escape"
commonShortcuts.delete   // "delete"
commonShortcuts.refresh  // "ctrl+r"
commonShortcuts.undo     // "ctrl+z"
commonShortcuts.redo     // "ctrl+shift+z"
```

## 代码位置

```
web/src/
└── composables/
    └── useKeyboard.ts    # 键盘快捷键 Composable
```

## 注意事项

- 输入框中只有 Escape 键会触发
- macOS 上 Ctrl 和 Cmd 被视为相同
- 建议为每个快捷键提供 description 用于帮助面板
