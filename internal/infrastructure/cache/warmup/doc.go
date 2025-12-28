// Package warmup 提供缓存预热服务。
//
// # 概述
//
// 本包负责应用启动时的缓存预热，确保热点数据在服务启动后立即可用。
//
// # 主要组件
//
//   - [NewSettingCacheWarmer]: 系统设置缓存预热
//   - [NewSettingCategoryCacheWarmer]: 设置分类缓存预热
//
// # 预热机制
//
// 多实例安全的预热流程：
//  1. 使用分布式锁（SETNX）防止多实例重复预热
//  2. 双重检查避免不必要的数据库查询
//  3. 预热失败不阻塞启动，降级为惰性加载
//
// # 使用方式
//
//	warmer := warmup.NewSettingCacheWarmer(repo, cacheService)
//	if err := warmer.WarmUpWithTimeout(ctx); err != nil {
//	    slog.Warn("warmup failed, using lazy loading", "err", err)
//	}
//
// # 依赖注入原则
//
// 预热服务依赖原始仓储（非缓存装饰器），避免循环依赖。
package warmup
