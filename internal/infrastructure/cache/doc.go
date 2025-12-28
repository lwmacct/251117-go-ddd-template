// Package cache 提供缓存预热和管理服务。
//
// 本包负责系统启动时的缓存预热，确保热点数据在服务启动后立即可用。
//
// 主要组件：
//   - [SettingCacheWarmer]: 系统设置缓存预热服务
//
// 预热机制：
//   - 使用分布式锁（[cache.SettingCacheService.TryAcquireWarmupLock]）防止多实例重复预热
//   - 双重检查避免不必要的数据库查询
//   - 预热失败时降级为惰性加载
package cache
