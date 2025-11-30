/**
 * 键盘快捷键 Composable
 * 用于注册和管理全局或组件级别的键盘快捷键
 */
import { onMounted, onUnmounted, ref } from "vue";

export interface KeyboardShortcut {
  /** 按键组合，如 "ctrl+s", "escape", "ctrl+shift+p" */
  key: string;
  /** 回调函数 */
  handler: (event: KeyboardEvent) => void;
  /** 是否阻止默认行为 */
  preventDefault?: boolean;
  /** 是否阻止事件冒泡 */
  stopPropagation?: boolean;
  /** 描述（用于帮助提示） */
  description?: string;
  /** 是否仅在特定条件下触发 */
  when?: () => boolean;
}

/**
 * 解析按键字符串为标准格式
 */
function parseKeyString(keyString: string): {
  key: string;
  ctrl: boolean;
  shift: boolean;
  alt: boolean;
  meta: boolean;
} {
  const parts = keyString
    .toLowerCase()
    .split("+")
    .map((p) => p.trim());
  const key = parts[parts.length - 1] || "";

  return {
    key,
    ctrl: parts.includes("ctrl") || parts.includes("control"),
    shift: parts.includes("shift"),
    alt: parts.includes("alt"),
    meta: parts.includes("meta") || parts.includes("cmd") || parts.includes("command"),
  };
}

/**
 * 检查事件是否匹配快捷键
 */
function matchesShortcut(event: KeyboardEvent, parsed: ReturnType<typeof parseKeyString>): boolean {
  const eventKey = event.key.toLowerCase();

  // 检查修饰键
  if (parsed.ctrl !== (event.ctrlKey || event.metaKey)) return false;
  if (parsed.shift !== event.shiftKey) return false;
  if (parsed.alt !== event.altKey) return false;

  // 检查主键
  return eventKey === parsed.key || event.code.toLowerCase() === `key${parsed.key}`;
}

/**
 * 创建键盘快捷键管理器
 */
export function useKeyboard(shortcuts: KeyboardShortcut[]) {
  const isEnabled = ref(true);

  // 解析所有快捷键
  const parsedShortcuts = shortcuts.map((shortcut) => ({
    ...shortcut,
    parsed: parseKeyString(shortcut.key),
  }));

  const handleKeydown = (event: KeyboardEvent) => {
    if (!isEnabled.value) return;

    // 跳过输入框中的快捷键（除非是 Escape）
    const target = event.target as HTMLElement;
    const isInput = target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.isContentEditable;

    for (const shortcut of parsedShortcuts) {
      // 检查是否匹配
      if (!matchesShortcut(event, shortcut.parsed)) continue;

      // 检查条件
      if (shortcut.when && !shortcut.when()) continue;

      // 输入框中只允许 Escape 键
      if (isInput && shortcut.parsed.key !== "escape") continue;

      // 执行处理函数
      if (shortcut.preventDefault !== false) {
        event.preventDefault();
      }
      if (shortcut.stopPropagation) {
        event.stopPropagation();
      }

      shortcut.handler(event);
      break;
    }
  };

  onMounted(() => {
    window.addEventListener("keydown", handleKeydown);
  });

  onUnmounted(() => {
    window.removeEventListener("keydown", handleKeydown);
  });

  /**
   * 启用/禁用快捷键
   */
  const setEnabled = (enabled: boolean) => {
    isEnabled.value = enabled;
  };

  /**
   * 获取所有快捷键描述（用于帮助面板）
   */
  const getShortcutDescriptions = () => {
    return parsedShortcuts
      .filter((s) => s.description)
      .map((s) => ({
        key: s.key,
        description: s.description,
      }));
  };

  return {
    isEnabled,
    setEnabled,
    getShortcutDescriptions,
  };
}

/**
 * 简单的单个快捷键注册
 */
export function useShortcut(
  key: string,
  handler: (event: KeyboardEvent) => void,
  options: Omit<KeyboardShortcut, "key" | "handler"> = {}
) {
  return useKeyboard([{ key, handler, ...options }]);
}

/**
 * 常用快捷键预设
 */
export const commonShortcuts = {
  /** 保存 */
  save: "ctrl+s",
  /** 新建 */
  new: "ctrl+n",
  /** 搜索 */
  search: "ctrl+k",
  /** 关闭/取消 */
  escape: "escape",
  /** 删除 */
  delete: "delete",
  /** 刷新 */
  refresh: "ctrl+r",
  /** 撤销 */
  undo: "ctrl+z",
  /** 重做 */
  redo: "ctrl+shift+z",
  /** 全选 */
  selectAll: "ctrl+a",
} as const;
