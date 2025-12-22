# 菜单列表导出功能

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:22+4`
- [已实现功能](#已实现功能) `:26+23`
  - [导出特性](#导出特性) `:28+7`
  - [导出字段](#导出字段) `:35+14`
- [技术实现](#技术实现) `:49+28`
  - [树形扁平化](#树形扁平化) `:51+17`
  - [代码位置](#代码位置) `:68+9`
- [导出示例](#导出示例) `:77+8`

<!--TOC-->

## 需求背景

管理员需要导出菜单配置，用于备份、文档生成或跨环境迁移参考。

## 已实现功能

### 导出特性

- 一键导出菜单列表为 CSV 文件
- 自动扁平化树形结构
- 保留层级关系可视化
- UTF-8 编码，支持中文

### 导出字段

| 字段     | 说明                    |
| -------- | ----------------------- |
| ID       | 菜单 ID                 |
| 层级     | 使用 "─" 可视化层级深度 |
| 标题     | 菜单显示名称            |
| 路径     | 菜单路由路径            |
| 图标     | MDI 图标名称            |
| 父级菜单 | 父级菜单标题            |
| 排序     | 同级排序值              |
| 可见     | 是否在导航中显示        |
| 创建时间 | 创建时间                |

## 技术实现

### 树形扁平化

```typescript
// 递归扁平化菜单树
const flattenMenus = (menuList: Menu[], level = 0, parentTitle?: string): FlatMenu[] => {
  const result: FlatMenu[] = [];
  for (const menu of menuList) {
    const { children, ...rest } = menu;
    result.push({ ...rest, level, parent_title: parentTitle });
    if (children?.length > 0) {
      result.push(...flattenMenus(children, level + 1, menu.title));
    }
  }
  return result;
};
```

### 代码位置

```
web/src/pages/admin/menus/
├── index.vue                           # 添加导出按钮
└── composables/
    └── useMenus.ts                     # 添加 exportMenus 方法
```

## 导出示例

```csv
ID,层级,标题,路径,图标,父级菜单,排序,可见,创建时间
1,,系统管理,/admin,mdi-cog,-,1,是,2024/11/30 10:00
2,─,用户管理,/admin/users,mdi-account,系统管理,1,是,2024/11/30 10:00
3,─,角色管理,/admin/roles,mdi-shield,系统管理,2,是,2024/11/30 10:00
```
