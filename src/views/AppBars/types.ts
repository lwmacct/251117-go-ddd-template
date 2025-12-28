/**
 * AppBars 组件类型定义
 */

/**
 * AppBars 组件属性
 */
export interface AppBarProps {
  /** Logo 图片路径（支持 URL 或本地资源） */
  logo?: string;
  /** Logo 备用图标（无 logo 时使用） */
  logoIcon?: string;
  /** 应用名称（方形 logo 时显示在旁边） */
  appName?: string;
  /** 是否显示导航图标 */
  showNavIcon?: boolean;
  /** 导航图标 */
  navIcon?: string;
  /** 是否显示抽屉菜单 */
  showDrawer?: boolean;
  /** 抽屉宽度 */
  drawerWidth?: number;
  /** 阴影高度 */
  elevation?: number;
  /** 背景色 */
  color?: string;
  /** 高度 */
  height?: number;
}

/**
 * 菜单项接口
 */
export interface MenuItem {
  /** 唯一标识 */
  id: string;
  /** 标题 */
  title: string;
  /** 图标 */
  icon: string;
  /** 路由路径 */
  path: string;
  /** 分类 */
  category?: string;
  /** 描述 */
  description?: string;
  /** 是否收藏 */
  isFavorite?: boolean;
  /** 徽章 */
  badge?: string;
}

/**
 * 访问记录接口
 */
export interface AccessRecord {
  /** 路由路径 */
  path: string;
  /** 页面标题 */
  title: string;
  /** 页面图标 */
  icon: string;
  /** 页面分类 */
  category?: string;
  /** 访问时间戳 (毫秒) */
  timestamp: number;
  /** 是否收藏 */
  isFavorite?: boolean;
}

/**
 * 抽屉菜单配置
 */
export interface DrawerConfig {
  /** 宽度 */
  width: number;
  /** 是否临时模式 */
  temporary: boolean;
}

/**
 * 悬停面板类型
 */
export type HoverPanelType = "recent" | "all" | null;
