/**
 * 系统设置相关类型定义
 */

/** 设置项类型 */
export type SettingType = "string" | "number" | "boolean" | "json";

/** 设置项 */
export interface Setting {
  key: string;
  value: string;
  category: string;
  value_type: SettingType;
  label?: string;
  created_at?: string;
  updated_at?: string;
}

/** 创建设置请求 */
export interface CreateSettingRequest {
  key: string;
  value: string;
  category: string;
  value_type?: SettingType;
  label?: string;
}

/** 更新单个设置请求 */
export interface UpdateSettingRequest {
  value: string;
  value_type?: SettingType;
  label?: string;
}

/** 批量更新设置请求 */
export interface UpdateSettingsRequest {
  settings: Array<{
    key: string;
    value: string;
  }>;
}
