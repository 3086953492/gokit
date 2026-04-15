# cache

两级缓存组件：统一封装本地内存缓存与 Redis 缓存，提供类型化访问器、批量失效能力，以及并发防击穿（singleflight）。

## 安装

```bash
go get github.com/3086953492/gokit/cache
```

## 核心能力

- `Manager`：缓存主入口，线程安全；负责本地缓存与 Redis 的读写协同。
- `Keyed[T]`：类型化缓存访问器，预绑定前缀和 TTL，业务侧只传可变 key 片段。
- `Group`：缓存失效组，聚合多个前缀后一键失效。
- `BuildKey`：统一 key 拼接函数，生成形如 `prefix|part1|part2` 的 key。
- 防击穿：`Manager.GetOrSet` / `Keyed.GetOrSet` 内部使用 `singleflight` 合并并发回源。

## Redis 依赖接口

`Manager` 依赖 `RedisBackend` 接口（`redis.Manager` 可直接满足）：

- `GetBytes(ctx, key)`
- `SetBytes(ctx, key, value, ttl)`
- `Del(ctx, keys...)`
- `ScanKeys(ctx, pattern, count)`
- `Exists(ctx, key)`

## 快速开始

```go
package main

import (
	"context"
	"time"

	"github.com/3086953492/gokit/cache"
)

type ProductList struct {
	Items []string `json:"items"`
	Total int      `json:"total"`
}

func main() {
	// backend 需实现 cache.RedisBackend（例如 redis.Manager）
	var backend cache.RedisBackend

	mgr, err := cache.NewManager(
		backend,
		cache.WithDefaultTTL(5*time.Minute),
		cache.WithLocalCache(true),
		cache.WithLocalCacheTTL(1*time.Minute),
		cache.WithLocalCacheMaxSize(1000),
		cache.WithScanCount(100),
	)
	if err != nil {
		panic(err)
	}
	defer mgr.Close()

	ctx := context.Background()

	listCache := cache.NewKeyed[ProductList](
		mgr,
		"product:list",
		cache.WithKeyedTTL(10*time.Minute),
	)

	// 命中则直接返回，未命中则执行回调并写入缓存
	result, err := listCache.GetOrSet(ctx, func() (*ProductList, error) {
		return &ProductList{
			Items: []string{"A", "B"},
			Total: 2,
		}, nil
	}, 1, 20, "zh-CN")
	if err != nil {
		panic(err)
	}

	_ = result
}
```

## 批量失效示例

```go
g := cache.NewGroup(
	mgr,
	"product:list",
	"product:detail",
)

if err := g.InvalidateAll(ctx); err != nil {
	// handle error
}
```

## Manager 配置项

默认配置如下：

- `DefaultTTL`: `5m`（Redis 默认过期时间）
- `LocalCacheEnabled`: `true`（默认启用本地缓存）
- `LocalCacheTTL`: `1m`（本地缓存 TTL）
- `LocalCacheMaxSize`: `1000`（本地缓存最大条目数，`0` 表示不限制）
- `ScanCount`: `100`（按前缀删除时，扫描建议数量）

可用选项：

- `WithDefaultTTL(ttl)`
- `WithLocalCache(enabled)`
- `WithLocalCacheTTL(ttl)`
- `WithLocalCacheMaxSize(size)`
- `WithScanCount(count)`

## 关键行为说明

- 本地缓存采用惰性过期策略，无后台清理 goroutine。
- 本地缓存达到上限时，先清理过期条目；若仍超限，会删除一个已有条目腾挪空间。
- `Manager.Get` 未命中时返回 `ErrCacheMiss`。
- `DeleteByPrefix` 会清理本地缓存，并通过 Redis `ScanKeys + Del` 删除远端前缀 key。
- `Close` 会标记管理器关闭并清空本地缓存；关闭后调用方法返回 `ErrManagerClosed`。

## 错误约定

包内预定义错误（可通过 `errors.Is` 判断）：

- `ErrNilRedisBackend`
- `ErrCacheMiss`
- `ErrManagerClosed`
- `ErrInvalidKey`（当前实现中暂未主动返回，保留为公共错误约定）

## API 概览

`Manager`：

- `Get(ctx, key, dest)`
- `Set(ctx, key, value, ttl...)`
- `GetOrSet(ctx, key, dest, fn, ttl...)`
- `Delete(ctx, key)`
- `DeleteByPrefix(ctx, prefix)`
- `DeleteByPrefixes(ctx, prefixes)`
- `Exists(ctx, key)`
- `Close()`

`Keyed[T]`：

- `NewKeyed[T](mgr, prefix, opts...)`
- `Prefix()`
- `Get(ctx, parts...)`
- `Set(ctx, value, parts...)`
- `GetOrSet(ctx, fn, parts...)`
- `Delete(ctx, parts...)`
- `InvalidateAll(ctx)`
- `Exists(ctx, parts...)`
