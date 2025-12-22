# 开发计划

本专栏记录项目的开发计划和进度跟踪。

<!--TOC-->

## Table of Contents

- [当前迭代](#当前迭代) `:20+10`
  - [v1.1.0 - 用户体验增强](#v110-用户体验增强) `:22+8`
- [已完成功能](#已完成功能) `:30+93`
  - [v1.2.0 - 功能扩展（前端）](#v120-功能扩展前端) `:32+85`
  - [计划中](#计划中) `:117+6`
- [状态说明](#状态说明) `:123+10`
- [已完成](#已完成) `:133+12`
  - [v1.0.0 - 核心功能](#v100-核心功能) `:135+10`

<!--TOC-->

## 当前迭代

### v1.1.0 - 用户体验增强

| 功能                               | 状态        | 负责人   | 完成日期   |
| ---------------------------------- | ----------- | -------- | ---------- |
| [账户删除功能](./account-deletion) | ✅ 代码完成 | AI Agent | 2024-11-30 |
| [审计日志导出](./auditlog-export)  | ✅ 代码完成 | AI Agent | 2024-11-30 |
| [用户头像上传](./avatar-upload)    | ✅ 代码完成 | AI Agent | 2024-11-30 |

## 已完成功能

### v1.2.0 - 功能扩展（前端）

#### 业务功能

| 功能                           | 状态      | 说明           |
| ------------------------------ | --------- | -------------- |
| [批量用户导入](./batch-import) | ✅ 已完成 | 前后端均已完成 |
| [用户列表导出](./user-export)  | ✅ 已完成 | 导出用户为 CSV |
| [角色列表导出](./role-export)  | ✅ 已完成 | 导出角色为 CSV |
| [菜单列表导出](./menu-export)  | ✅ 已完成 | 导出菜单为 CSV |
| [数据仪表盘](./dashboard)      | ✅ 已实现 | 可视化统计图表 |

#### UI 组件

| 功能                                  | 状态      | 说明               |
| ------------------------------------- | --------- | ------------------ |
| [密码强度指示器](./password-strength) | ✅ 已完成 | 可视化密码强度     |
| [复制到剪贴板](./clipboard-copy)      | ✅ 已完成 | 快速复制 ID 等信息 |
| [确认对话框组件](./confirm-dialog)    | ✅ 已完成 | 可复用确认对话框   |
| [面包屑导航](./breadcrumb)            | ✅ 已完成 | 显示当前页面位置   |
| [空状态组件](./empty-state)           | ✅ 已完成 | 统一空状态样式     |
| [加载遮罩组件](./loading-overlay)     | ✅ 已完成 | 区域加载遮罩       |
| [页面过渡动画](./page-transition)     | ✅ 已完成 | 页面切换动画       |
| [回到顶部按钮](./back-to-top)         | ✅ 已完成 | 快速返回顶部       |
| [骨架屏加载](./skeleton-loader)       | ✅ 已完成 | 加载占位骨架屏     |

#### 工具函数

| 功能                                                          | 状态      | 说明                                |
| ------------------------------------------------------------- | --------- | ----------------------------------- |
| [防抖搜索](./debounce-search)                                 | ✅ 已完成 | 输入自动触发搜索                    |
| [键盘快捷键](./keyboard-shortcuts)                            | ✅ 已完成 | 键盘快捷键支持                      |
| [日期时间格式化](./date-utils)                                | ✅ 已完成 | 统一日期格式化                      |
| [数字格式化](./number-utils)                                  | ✅ 已完成 | 货币/百分比/文件大小格式化          |
| [存储工具](./storage-utils)                                   | ✅ 已完成 | 响应式 localStorage/sessionStorage  |
| [表单验证工具](./form-validation)                             | ✅ 已完成 | 统一验证规则                        |
| [字符串工具函数](./string-utils)                              | ✅ 已完成 | 字符串处理与格式化                  |
| [窗口尺寸检测](./window-size)                                 | ✅ 已完成 | 响应式窗口/断点/媒体查询            |
| [交叉观察器](./intersection-observer)                         | ✅ 已完成 | 懒加载/无限滚动/可见性检测          |
| [点击外部检测](./click-outside)                               | ✅ 已完成 | 下拉菜单/弹出框外部点击             |
| [异步状态管理](./async-state)                                 | ✅ 已完成 | loading/error/data/重试/轮询        |
| [网络状态检测](./network-status)                              | ✅ 已完成 | 在线/离线/弱网检测                  |
| [全屏 API](./fullscreen)                                      | ✅ 已完成 | 元素/文档全屏控制                   |
| [滚动锁定](./scroll-lock)                                     | ✅ 已完成 | 模态框/抽屉滚动锁定                 |
| [焦点陷阱](./focus-trap)                                      | ✅ 已完成 | 模态框焦点管理/无障碍               |
| [拖拽排序](./draggable)                                       | ✅ 已完成 | 列表拖拽排序/文件拖放               |
| [表单状态管理](./form-state)                                  | ✅ 已完成 | 脏状态/验证/提交/重置               |
| [倒计时](./countdown)                                         | ✅ 已完成 | 验证码/活动倒计时/秒表              |
| [事件总线](./event-bus)                                       | ✅ 已完成 | 组件间通信/事件驱动                 |
| [历史记录](./history)                                         | ✅ 已完成 | 撤销/重做/状态快照                  |
| [平台检测](./platform)                                        | ✅ 已完成 | 浏览器/操作系统/设备检测            |
| [限流工具](./throttle)                                        | ✅ 已完成 | 节流/防抖/重试/超时                 |
| [数组工具](./array-utils)                                     | ✅ 已完成 | 分组/去重/排序/树形操作             |
| [对象工具](./object-utils)                                    | ✅ 已完成 | 深拷贝/路径操作/差异比较            |
| [ID 生成工具](./id-utils)                                     | ✅ 已完成 | UUID/NanoID/雪花ID/ULID             |
| [颜色工具](./color-utils)                                     | ✅ 已完成 | 颜色转换/调整/生成调色板            |
| [URL 工具](./url-utils)                                       | ✅ 已完成 | URL解析/查询参数/路径操作           |
| [DOM 工具](./dom-utils)                                       | ✅ 已完成 | 元素操作/事件/滚动/可见性           |
| [Cookie 工具](./cookie-utils)                                 | ✅ 已完成 | Cookie读写/JSON/管理器              |
| [类型工具](./type-utils)                                      | ✅ 已完成 | 类型检查/转换/断言                  |
| [Timer Composable](./timer-composable)                        | ✅ 已完成 | 定时器/间隔器/任务调度              |
| [Permission Composable](./permission-composable)              | ✅ 已完成 | 浏览器权限管理                      |
| [Mouse Composable](./mouse-composable)                        | ✅ 已完成 | 鼠标位置/悬停/拖放                  |
| [Geolocation Composable](./geolocation-composable)            | ✅ 已完成 | 地理位置/距离计算                   |
| [Preferences Composable](./preferences-composable)            | ✅ 已完成 | 用户偏好设置检测与控制              |
| [Title Composable](./title-composable)                        | ✅ 已完成 | 文档标题/Favicon/脚本加载           |
| [WebSocket Composable](./websocket-composable)                | ✅ 已完成 | WebSocket 连接管理/自动重连/心跳    |
| [EventSource Composable](./eventsource-composable)            | ✅ 已完成 | Server-Sent Events 连接管理         |
| [BroadcastChannel Composable](./broadcast-channel-composable) | ✅ 已完成 | 跨标签页通信/状态同步/Leader 选举   |
| [Image Composable](./image-composable)                        | ✅ 已完成 | 图片加载/预加载/懒加载/验证/压缩    |
| [Clone Composable](./clone-composable)                        | ✅ 已完成 | 对象克隆/脏状态检测/状态快照        |
| [Fetch Composable](./fetch-composable)                        | ✅ 已完成 | HTTP 请求/缓存/分页/无限加载        |
| [VModel Composable](./vmodel-composable)                      | ✅ 已完成 | v-model 双向绑定/防抖/节流/代理     |
| [SWR Composable](./swr-composable)                            | ✅ 已完成 | Stale-While-Revalidate 数据获取策略 |
| [Watch Composable](./watch-composable)                        | ✅ 已完成 | 增强的 watch 工具函数               |
| [Computed Composable](./computed-composable)                  | ✅ 已完成 | 增强的 computed 工具函数            |
| [Ref Composable](./ref-composable)                            | ✅ 已完成 | 增强的 ref 工具函数                 |
| [Reactive Composable](./reactive-composable)                  | ✅ 已完成 | 增强的 reactive 工具函数            |
| [Lifecycle Composable](./lifecycle-composable)                | ✅ 已完成 | 增强的生命周期钩子工具函数          |
| [Provide Composable](./provide-composable)                    | ✅ 已完成 | 增强的依赖注入工具函数              |
| [Component Composable](./component-composable)                | ✅ 已完成 | 组件相关工具函数                    |
| [Transition Composable](./transition-composable)              | ✅ 已完成 | 过渡和动画工具函数                  |
| [Queue Composable](./queue-composable)                        | ✅ 已完成 | 队列相关工具函数                    |
| [Media Composable](./media-composable)                        | ✅ 已完成 | 媒体播放控制工具函数                |

### 计划中

| 功能         | 状态      | 说明                           |
| ------------ | --------- | ------------------------------ |
| 邮件通知系统 | 📋 计划中 | 系统事件邮件通知（需后端支持） |

## 状态说明

| 图标 | 状态   |
| ---- | ------ |
| ✅   | 已完成 |
| 🚧   | 进行中 |
| 📋   | 计划中 |
| ❌   | 已取消 |
| ⏸️   | 暂停   |

## 已完成

### v1.0.0 - 核心功能

- ✅ 用户认证（登录/注册）
- ✅ 双因素认证 (2FA)
- ✅ 用户管理 (CRUD)
- ✅ 角色权限管理 (RBAC)
- ✅ 菜单管理
- ✅ 系统设置
- ✅ 审计日志
- ✅ Personal Access Token
