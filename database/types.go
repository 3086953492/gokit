package database

// MySQLConfig 包含构建 MySQL DSN 所需的连接参数
type MySQLConfig struct {
	Host      string // 数据库主机地址
	Port      int    // 端口号
	User      string // 用户名
	Password  string // 密码
	DBName    string // 数据库名
	Charset   string // 字符集，如 utf8mb4
	ParseTime bool   // 是否解析时间类型
	Loc       string // 时区，如 Local
}
