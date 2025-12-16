package cache

import (
	"sync"
	"time"
)

// localCache 简单的本地 TTL 缓存实现
// 使用惰性过期策略，无后台 goroutine
type localCache struct {
	mu      sync.RWMutex
	data    map[string]*localCacheEntry
	maxSize int
}

// newLocalCache 创建一个新的本地缓存
func newLocalCache(maxSize int) *localCache {
	return &localCache{
		data:    make(map[string]*localCacheEntry),
		maxSize: maxSize,
	}
}

// get 获取缓存值
// 如果 key 不存在或已过期，返回 nil
func (c *localCache) get(key string) []byte {
	c.mu.RLock()
	entry, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil
	}

	if entry.isExpired() {
		// 惰性删除过期条目
		c.mu.Lock()
		// 双重检查
		if e, exists := c.data[key]; exists && e.isExpired() {
			delete(c.data, key)
		}
		c.mu.Unlock()
		return nil
	}

	return entry.value
}

// set 设置缓存值
func (c *localCache) set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果达到最大容量，执行简单的清理（删除过期条目）
	if c.maxSize > 0 && len(c.data) >= c.maxSize {
		c.evictExpired()
		// 如果仍然满，删除一个随机条目（简单 LRU 替代）
		if len(c.data) >= c.maxSize {
			for k := range c.data {
				delete(c.data, k)
				break
			}
		}
	}

	c.data[key] = &localCacheEntry{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

// delete 删除缓存值
func (c *localCache) delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

// deleteByPrefix 删除指定前缀的所有缓存
func (c *localCache) deleteByPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(c.data, k)
		}
	}
}

// clear 清空所有缓存
func (c *localCache) clear() {
	c.mu.Lock()
	c.data = make(map[string]*localCacheEntry)
	c.mu.Unlock()
}

// evictExpired 清理过期条目（需要持有写锁）
func (c *localCache) evictExpired() {
	now := time.Now()
	for k, v := range c.data {
		if now.After(v.expireAt) {
			delete(c.data, k)
		}
	}
}

