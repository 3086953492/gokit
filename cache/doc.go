// Package cache 提供两级缓存（本地内存 + Redis）管理。
//
// 核心类型：
//   - Manager：底层缓存管理器，负责本地缓存与 Redis 的读写。
//   - Keyed[T]：预定义的类型化缓存访问器，绑定前缀、TTL 和值类型，
//     调用方只需提供可变的键部分即可完成缓存操作。
//   - Group：缓存失效组，将多个前缀聚合在一起用于批量失效。
//
// 典型用法：
//
//	// 创建 Manager
//	mgr, _ := cache.NewManager(redisBackend, cache.WithDefaultTTL(5*time.Minute))
//
//	// 预定义缓存访问器
//	listCache := cache.NewKeyed[ProductList](mgr, "product:list", cache.WithKeyedTTL(10*time.Minute))
//	detailCache := cache.NewKeyed[ProductDetail](mgr, "product:detail")
//
//	// 读写
//	result, err := listCache.GetOrSet(ctx, loadFn, page, pageSize, lang)
//
//	// 批量失效
//	g := cache.NewGroup(mgr, listCache.Prefix(), detailCache.Prefix())
//	g.InvalidateAll(ctx)
package cache
