package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ManagerConfig 数据库管理器配置
type ManagerConfig struct {
	// 连接池配置
	MaxIdleConns    int           // 最大空闲连接数，默认 10
	MaxOpenConns    int           // 最大打开连接数，默认 100
	ConnMaxLifetime time.Duration // 连接最大存活时间，默认 1 小时

	// Gorm 配置
	LogMode logger.LogLevel // 日志级别，默认 Info
}

// DefaultManagerConfig 返回默认配置
func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogMode:         logger.Info,
	}
}

// Manager 数据库管理器
// 封装 gorm.DB，由使用者自己持有和管理生命周期
type Manager struct {
	db *gorm.DB
}

// NewManager 使用 dialector 创建数据库管理器
// dialector: 数据库方言，如 mysql.Open(dsn)、postgres.Open(dsn) 等
// cfg: 可选配置，不传则使用默认配置
func NewManager(dialector gorm.Dialector, cfg ...ManagerConfig) (*Manager, error) {
	config := DefaultManagerConfig()
	if len(cfg) > 0 {
		config = cfg[0]
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogMode),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &Manager{db: db}, nil
}

// DB 返回底层 *gorm.DB 实例
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// AutoMigrate 自动迁移数据表结构
// models: 需要迁移的模型列表
func (m *Manager) AutoMigrate(models ...any) error {
	if len(models) == 0 {
		return nil
	}
	return m.db.AutoMigrate(models...)
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	return nil
}
