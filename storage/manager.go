package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"
)

// Manager 对象存储管理器，作为 storage 包对外的统一入口。
// Manager 是线程安全的，内部持有一个 Store 实现。
type Manager struct {
	mu    sync.RWMutex
	store Store
	opts  *Options
}

// NewManager 创建存储管理器。
// 必须通过 WithStore 指定具体的存储后端。
func NewManager(opts ...Option) (*Manager, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.store == nil {
		return nil, fmt.Errorf("%w: store is required, use WithStore option", ErrInvalidConfig)
	}

	return &Manager{
		store: o.store,
		opts:  o,
	}, nil
}

// Upload 上传对象到存储后端。
// key 为对象唯一标识，r 为对象内容流。
func (m *Manager) Upload(ctx context.Context, key string, r io.Reader, opts ...WriteOption) (*ObjectMeta, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	writeOpts := &WriteOptions{}
	for _, opt := range opts {
		opt(writeOpts)
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	meta, err := store.Upload(ctx, key, r, writeOpts)
	if err != nil {
		return nil, fmt.Errorf("storage upload %s: %w", key, err)
	}
	return meta, nil
}

// Download 从存储后端下载对象。
// 返回对象内容流（调用方负责 Close）和元信息。
func (m *Manager) Download(ctx context.Context, key string, opts ...ReadOption) (io.ReadCloser, *ObjectMeta, error) {
	if err := validateKey(key); err != nil {
		return nil, nil, err
	}

	readOpts := &ReadOptions{}
	for _, opt := range opts {
		opt(readOpts)
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	rc, meta, err := store.Download(ctx, key, readOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("storage download %s: %w", key, err)
	}
	return rc, meta, nil
}

// Delete 删除存储后端中的对象。
func (m *Manager) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	if err := validateKey(key); err != nil {
		return err
	}

	deleteOpts := &DeleteOptions{}
	for _, opt := range opts {
		opt(deleteOpts)
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	if err := store.Delete(ctx, key, deleteOpts); err != nil {
		return fmt.Errorf("storage delete %s: %w", key, err)
	}
	return nil
}

// List 列举指定前缀下的对象。
// 返回的切片非 nil（即使为空）。
func (m *Manager) List(ctx context.Context, prefix string, opts ...ListOption) (*ListResult, error) {
	listOpts := &ListOptions{
		MaxKeys: 1000, // 默认值
	}
	for _, opt := range opts {
		opt(listOpts)
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	result, err := store.List(ctx, prefix, listOpts)
	if err != nil {
		return nil, fmt.Errorf("storage list prefix=%s: %w", prefix, err)
	}

	// 确保返回非 nil 切片
	if result.Objects == nil {
		result.Objects = []*ObjectMeta{}
	}
	if result.CommonPrefixes == nil {
		result.CommonPrefixes = []string{}
	}

	return result, nil
}

// Exists 检查对象是否存在。
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	if err := validateKey(key); err != nil {
		return false, err
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	exists, err := store.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("storage exists %s: %w", key, err)
	}
	return exists, nil
}

// Head 获取对象元信息，不下载内容。
func (m *Manager) Head(ctx context.Context, key string) (*ObjectMeta, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	meta, err := store.Head(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("storage head %s: %w", key, err)
	}
	return meta, nil
}

// Store 返回内部的 Store 实现。
// 一般情况下不需要直接访问，仅供特殊场景使用。
func (m *Manager) Store() Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.store
}

// validateKey 校验对象 key 的合法性。
func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("%w: key cannot be empty", ErrInvalidKey)
	}
	return nil
}

// DeleteByURL 根据对象 URL 删除存储后端中的对象。
// 仅支持本库生成的公开直链（即 ObjectMeta.URL 格式）；
// 若当前 Store 未实现 URLKeyResolver，返回 ErrURLDeleteUnsupported。
func (m *Manager) DeleteByURL(ctx context.Context, rawURL string, opts ...DeleteOption) error {
	key, err := m.parseKeyFromURL(rawURL)
	if err != nil {
		return err
	}
	return m.Delete(ctx, key, opts...)
}

// parseKeyFromURL 解析 rawURL 并提取对象 key。
// 返回的 key 已做 URL 解码；若 URL 不合法或域名不允许，返回对应错误。
func (m *Manager) parseKeyFromURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("%w: url cannot be empty", ErrInvalidURL)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	// 必须是 http(s) scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("%w: unsupported scheme %q", ErrInvalidURL, u.Scheme)
	}

	m.mu.RLock()
	store := m.store
	m.mu.RUnlock()

	resolver, ok := store.(URLKeyResolver)
	if !ok {
		return "", ErrURLDeleteUnsupported
	}

	if !isHostAllowed(u.Host, resolver.AllowedHosts()) {
		return "", fmt.Errorf("%w: host %q", ErrDomainNotAllowed, u.Host)
	}

	key, err := resolver.KeyFromURL(u)
	if err != nil {
		return "", err
	}

	return key, nil
}

// isHostAllowed 检查 host 是否在允许列表中（大小写不敏感）。
func isHostAllowed(host string, allowed []string) bool {
	for _, h := range allowed {
		if equalFoldHost(host, h) {
			return true
		}
	}
	return false
}

// equalFoldHost 大小写不敏感比较两个 host。
func equalFoldHost(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
