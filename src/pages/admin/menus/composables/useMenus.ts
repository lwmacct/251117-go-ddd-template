/**
 * Admin 菜单管理 Composable
 */
import { ref } from "vue";
import { adminMenuApi, extractData, type Menu } from "@/api";
import { type HandlerCreateMenuRequest, type HandlerUpdateMenuRequest, type HandlerReorderMenusRequest } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";

// 扁平化菜单结构（用于导出）
interface FlatMenu extends Omit<Menu, "children"> {
  level: number;
  parent_title?: string;
}

export function useMenus() {
  const menus = ref<Menu[]>([]);
  const loading = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  const fetchMenus = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminMenuApi.apiAdminMenusGet();
      menus.value = (response.data.data ?? []) as Menu[];
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取菜单列表失败";
      console.error("Failed to fetch menus:", error);
    } finally {
      loading.value = false;
    }
  };

  const fetchMenu = async (id: number): Promise<Menu | null> => {
    try {
      const response = await adminMenuApi.apiAdminMenusIdGet(id);
      return extractData<Menu>(response.data) ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取菜单详情失败";
      return null;
    }
  };

  const createMenu = async (data: HandlerCreateMenuRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminMenuApi.apiAdminMenusPost(data);
      successMessage.value = "菜单创建成功";
      await fetchMenus();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建菜单失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const updateMenu = async (id: number, data: HandlerUpdateMenuRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminMenuApi.apiAdminMenusIdPut(id, data);
      successMessage.value = "菜单更新成功";
      await fetchMenus();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新菜单失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const deleteMenu = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminMenuApi.apiAdminMenusIdDelete(id);
      successMessage.value = "菜单删除成功";
      await fetchMenus();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除菜单失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const reorderMenus = async (data: HandlerReorderMenusRequest): Promise<boolean> => {
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminMenuApi.apiAdminMenusReorderPost(data);
      successMessage.value = "菜单排序已更新";
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新排序失败";
      return false;
    }
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  /**
   * 扁平化菜单树结构
   */
  const flattenMenus = (menuList: Menu[], level = 0, parentTitle?: string): FlatMenu[] => {
    const result: FlatMenu[] = [];
    for (const menu of menuList) {
      const { children, ...rest } = menu;
      result.push({
        ...rest,
        level,
        parent_title: parentTitle,
      });
      if (children && children.length > 0) {
        result.push(...flattenMenus(children, level + 1, menu.title));
      }
    }
    return result;
  };

  /**
   * 导出菜单列表为 CSV
   */
  const exportMenus = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      // 获取所有菜单
      const response = await adminMenuApi.apiAdminMenusGet();
      const menuList = (response.data.data ?? []) as Menu[];

      if (menuList.length === 0) {
        errorMessage.value = "没有可导出的数据";
        return;
      }

      // 扁平化菜单树
      const flatMenus = flattenMenus(menuList);

      // 定义 CSV 列
      const columns: CSVColumn<FlatMenu>[] = [
        { header: "ID", key: "id" },
        { header: "层级", key: (item) => "─".repeat(item.level) + (item.level > 0 ? " " : "") },
        { header: "标题", key: "title" },
        { header: "路径", key: "path" },
        { header: "图标", key: (item) => item.icon || "-" },
        { header: "父级菜单", key: (item) => item.parent_title || "-" },
        { header: "排序", key: "order" },
        { header: "可见", key: (item) => (item.visible ? "是" : "否") },
        {
          header: "创建时间",
          key: (item) =>
            item.created_at
              ? new Date(item.created_at).toLocaleString("zh-CN", {
                  year: "numeric",
                  month: "2-digit",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                })
              : "-",
        },
      ];

      // 生成文件名
      const timestamp = new Date().toISOString().slice(0, 10);
      const filename = `菜单列表_${timestamp}.csv`;

      // 导出
      exportToCSV(flatMenus, columns, { filename, withBOM: true });
      successMessage.value = `成功导出 ${flatMenus.length} 个菜单项`;
    } catch (error) {
      errorMessage.value = (error as Error).message || "导出失败";
      console.error("Failed to export menus:", error);
    } finally {
      loading.value = false;
    }
  };

  return {
    menus,
    loading,
    errorMessage,
    successMessage,
    fetchMenus,
    fetchMenu,
    createMenu,
    updateMenu,
    deleteMenu,
    reorderMenus,
    clearMessages,
    exportMenus,
  };
}
