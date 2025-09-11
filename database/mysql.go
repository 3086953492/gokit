package database

import (
	"fmt"

	"github.com/3086953492/YaBase/config/types"
)

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
