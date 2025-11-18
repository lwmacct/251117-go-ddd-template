/**
 * 菜单项类型定义
 */
export interface MenuItem {
  /**
   * 菜单标题
   */
  title: string;
  /**
   * 路由路径
   */
  path: string;
  /**
   * 图标（MDI 图标名称）
   */
  icon?: string;
  /**
   * 是否精确匹配路由（默认 false，使用 startsWith）
   */
  exact?: boolean;
  /**
   * 徽章内容（可选）
   */
  badge?: string | number;
  /**
   * 徽章颜色（可选）
   */
  badgeColor?: string;
  /**
   * 是否禁用
   */
  disabled?: boolean;
}
