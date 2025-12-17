package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/3086953492/gokit/config/internal/loader"
	"github.com/3086953492/gokit/config/internal/resolve"
	"github.com/3086953492/gokit/config/internal/watcher"
)

// Manager 配置管理器，提供配置加载、访问和可选的热更新功能
type Manager struct {
	opts       *Options
	configPath string

	mu     sync.RWMutex
	config Config

	watching bool
	watchMu  sync.Mutex
}

// NewManager 创建配置管理器
func NewManager(opts ...Option) (*Manager, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 解析配置文件路径
	configPath, err := resolve.ResolvePath(resolve.ResolveOptions{
		ConfigFile:    options.ConfigFile,
		ConfigDir:     options.ConfigDir,
		EnvConfigKey:  options.EnvConfigKey,
		Mode:          options.Mode,
		ModeConfigMap: options.ModeConfigMap,
		Formats:       options.Formats,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigNotFound, err)
	}

	options.Logger.Info("resolved config path", "path", configPath)

	return &Manager{
		opts:       options,
		configPath: configPath,
		config:     DefaultConfig(),
	}, nil
}

// Load 加载配置文件，返回配置的值拷贝
func (m *Manager) Load(ctx context.Context) (Config, error) {
	var cfg Config
	if err := loader.Load(m.configPath, &cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	m.mu.Lock()
	m.config = cfg
	m.mu.Unlock()

	m.opts.Logger.Info("config loaded", "path", m.configPath)
	return cfg, nil
}

// Config 返回当前配置的值拷贝（线程安全）
func (m *Manager) Config() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// ConfigPath 返回当前使用的配置文件路径
func (m *Manager) ConfigPath() string {
	return m.configPath
}

// Watch 开始监听配置文件变化，阻塞直到 ctx 取消
// onReload 会在配置重新加载后被调用（仅在加载成功时）
// 若配置加载失败，会记录日志但不覆盖旧配置
func (m *Manager) Watch(ctx context.Context, onReload func(Config)) error {
	m.watchMu.Lock()
	if m.watching {
		m.watchMu.Unlock()
		return ErrAlreadyWatching
	}
	m.watching = true
	m.watchMu.Unlock()

	w, err := watcher.New(m.configPath, watcherLogger{m.opts.Logger})
	if err != nil {
		m.watchMu.Lock()
		m.watching = false
		m.watchMu.Unlock()
		return fmt.Errorf("create watcher: %w", err)
	}

	// 在 goroutine 中运行监听
	go func() {
		defer func() {
			w.Close()
			m.watchMu.Lock()
			m.watching = false
			m.watchMu.Unlock()
		}()

		w.Watch(ctx, func() {
			cfg, err := m.Load(ctx)
			if err != nil {
				m.opts.Logger.Error("reload config failed", "error", err)
				return
			}
			if onReload != nil {
				onReload(cfg)
			}
		})
	}()

	return nil
}

// watcherLogger 适配 Logger 接口到 watcher.Logger
type watcherLogger struct {
	logger Logger
}

func (w watcherLogger) Info(msg string, args ...any) {
	w.logger.Info(msg, args...)
}

func (w watcherLogger) Error(msg string, args ...any) {
	w.logger.Error(msg, args...)
}
