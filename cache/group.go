package cache

import "context"

// Group 将多个缓存前缀聚合在一起，用于批量失效。
//
// 典型用法是在 Service 构造时把读缓存的前缀收集到一个 Group，
// 写操作只需一行 InvalidateAll 即可清除所有关联缓存。
//
// 使用示例：
//
//	g := cache.NewGroup(mgr, listCache.Prefix(), detailCache.Prefix())
//	g.InvalidateAll(ctx)
type Group struct {
	mgr      *Manager
	prefixes []string
}

// NewGroup 创建缓存失效组。
func NewGroup(mgr *Manager, prefixes ...string) *Group {
	dst := make([]string, len(prefixes))
	copy(dst, prefixes)
	return &Group{mgr: mgr, prefixes: dst}
}

// InvalidateAll 删除组内所有前缀对应的缓存。
func (g *Group) InvalidateAll(ctx context.Context) error {
	return g.mgr.DeleteByPrefixes(ctx, g.prefixes)
}
