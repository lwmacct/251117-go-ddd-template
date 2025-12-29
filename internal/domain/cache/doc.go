// Package cache 定义缓存领域服务接口。
//
// 本包提供两类缓存接口：
//
// # 业务缓存服务（推荐）
//
// 针对特定业务域的缓存服务，封装 Key 命名、TTL 策略和失效逻辑：
//   - [SettingCategoryCacheService]: Setting 分类元数据缓存
//   - [UserSettingCacheService]: 用户有效设置值缓存
//   - [UserSettingQueryCacheService]: 用户自定义配置缓存
//   - [PermissionCacheService]: 用户权限缓存
//   - [UserWithRolesCacheService]: 用户角色关联缓存
//
// # 通用缓存仓储（已废弃）
//
// 低级别的 KV 缓存接口，不推荐新代码使用：
//   - [CommandRepository]: 写操作（已废弃）
//   - [QueryRepository]: 读操作（已废弃）
//
// 实现位于 infrastructure/cache 包。
package cache
