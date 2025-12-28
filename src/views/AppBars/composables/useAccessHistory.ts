/**
 * 访问历史管理 Composable
 *
 * 设计特点（相对 SaaS 模板的优化）：
 * - 独立存储完整页面信息，不依赖菜单项过滤
 * - 支持多用户隔离（按 userId）
 * - 5 秒内重复访问去重
 * - 登出时支持清理历史
 */

import { ref, computed } from "vue";
import type { AccessRecord } from "../types";

/** 配置常量 */
const MAX_RECORDS = 20;
const STORAGE_KEY_PREFIX = "recent_pages";
const DEDUP_INTERVAL_MS = 5000; // 5 秒去重窗口

/**
 * 获取存储 Key
 */
function getStorageKey(userId: number | string | null): string {
  const userKey = userId || "guest";
  return `${STORAGE_KEY_PREFIX}_${userKey}`;
}

/**
 * 从 localStorage 加载历史记录
 */
function loadHistory(userId: number | string | null): AccessRecord[] {
  try {
    const key = getStorageKey(userId);
    const data = localStorage.getItem(key);
    if (!data) return [];

    const records: AccessRecord[] = JSON.parse(data);
    return records
      .filter((r) => r.path && r.title && r.timestamp)
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, MAX_RECORDS);
  } catch {
    return [];
  }
}

/**
 * 保存历史记录到 localStorage
 */
function saveHistory(userId: number | string | null, records: AccessRecord[]): void {
  try {
    const key = getStorageKey(userId);
    const trimmed = records.slice(0, MAX_RECORDS);
    localStorage.setItem(key, JSON.stringify(trimmed));
  } catch {
    // 静默失败
  }
}

/**
 * 访问历史管理 Composable
 */
export function useAccessHistory(userId: number | string | null = null) {
  const history = ref<AccessRecord[]>(loadHistory(userId));

  /** 最近访问记录（响应式） */
  const recentPages = computed(() => history.value);

  /**
   * 记录新的访问
   */
  function recordAccess(record: Omit<AccessRecord, "timestamp">) {
    const now = Date.now();

    // 去重：同一路径 5 秒内不重复记录
    const lastRecord = history.value[0];
    if (lastRecord && lastRecord.path === record.path && now - lastRecord.timestamp < DEDUP_INTERVAL_MS) {
      return;
    }

    // 移除相同路径的旧记录
    history.value = history.value.filter((item) => item.path !== record.path);

    // 创建新记录
    const newRecord: AccessRecord = {
      ...record,
      timestamp: now,
    };

    // 添加到开头
    history.value = [newRecord, ...history.value];
    saveHistory(userId, history.value);
  }

  /**
   * 获取最近 N 条记录
   */
  function getRecent(limit = 10): AccessRecord[] {
    return history.value.slice(0, limit);
  }

  /**
   * 清空历史记录
   */
  function clearHistory() {
    history.value = [];
    try {
      const key = getStorageKey(userId);
      localStorage.removeItem(key);
    } catch {
      // 静默失败
    }
  }

  /**
   * 删除指定路径的记录
   */
  function removeByPath(path: string) {
    const beforeCount = history.value.length;
    history.value = history.value.filter((r) => r.path !== path);
    if (history.value.length < beforeCount) {
      saveHistory(userId, history.value);
    }
  }

  /**
   * 更新某条记录的收藏状态
   */
  function updateFavoriteStatus(path: string, isFavorite: boolean) {
    const record = history.value.find((item) => item.path === path);
    if (record) {
      record.isFavorite = isFavorite;
      saveHistory(userId, history.value);
    }
  }

  /**
   * 重新加载历史记录
   */
  function reload() {
    history.value = loadHistory(userId);
  }

  return {
    recentPages,
    recordAccess,
    getRecent,
    clearHistory,
    removeByPath,
    updateFavoriteStatus,
    reload,
  };
}
