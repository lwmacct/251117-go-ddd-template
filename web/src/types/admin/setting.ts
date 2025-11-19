/**
 * 系统设置相关类型定义
 */

/** 设置项类型 */
export type SettingType = 'string' | 'int' | 'bool' | 'json';

/** 设置项 */
export interface Setting {
  key: string;
  value: string;
  type: SettingType;
  description?: string;
  group?: string; // general, security, notification, backup
  created_at?: string;
  updated_at?: string;
}

/** 设置分组 */
export interface SettingGroup {
  group: string;
  label: string;
  settings: Setting[];
}

/** 更新设置请求 */
export interface UpdateSettingsRequest {
  settings: Array<{
    key: string;
    value: string;
  }>;
}
