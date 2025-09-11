package global

import (
	"fmt"
	"sync"
	"time"

	"github.com/3086953492/YaBase/config/types"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局数据库管理器
var (
	globalDB *gorm.DB
	dbMutex  sync.RWMutex
)

// GetGlobalDB 获取全局数据库实例
func GetGlobalDB() *gorm.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return globalDB
}

// InitDBWithDialector 使用dialector初始化数据库连接
// dialector: 数据库方言，如mysql.Open(dsn)、postgres.Open(dsn)等
func InitDBWithDialector(dialector gorm.Dialector) error {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 设置全局数据库实例
	dbMutex.Lock()
	globalDB = db
	dbMutex.Unlock()

	return nil
}

// AutoMigrateModels 自动迁移数据表结构
// models: 需要迁移的模型列表
func AutoMigrateModels(models ...any) error {
	db := GetGlobalDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	if len(models) == 0 {
		return nil
	}

	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if globalDB != nil {
		sqlDB, err := globalDB.DB()
		if err != nil {
			return fmt.Errorf("获取底层数据库连接失败: %w", err)
		}

		err = sqlDB.Close()
		globalDB = nil

		if err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}

	}
	return nil
}

// BuildMySQLDSN 构建MySQL DSN
func BuildMySQLDSN(cfg types.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)
}
