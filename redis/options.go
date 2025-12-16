package redis

import "time"

// Options 包含 Manager 的配置选项
type Options struct {
	// Address Redis 服务器地址，格式为 host:port
	Address string

	// Password Redis 认证密码
	Password string

	// DB 数据库编号
	DB int

	// DialTimeout 连接超时时间
	DialTimeout time.Duration

	// ReadTimeout 读取超时时间
	ReadTimeout time.Duration

	// WriteTimeout 写入超时时间
	WriteTimeout time.Duration

	// PoolSize 连接池大小
	PoolSize int

	// MinIdleConns 最小空闲连接数
	MinIdleConns int
}

// Option 是配置 Manager 的函数类型
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		Address:      "localhost:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	}
}

// WithAddress 设置 Redis 服务器地址
func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// WithPassword 设置 Redis 认证密码
func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// WithDB 设置数据库编号
func WithDB(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}

// WithDialTimeout 设置连接超时时间
func WithDialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = d
	}
}

// WithReadTimeout 设置读取超时时间
func WithReadTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = d
	}
}

// WithWriteTimeout 设置写入超时时间
func WithWriteTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = d
	}
}

// WithPoolSize 设置连接池大小
func WithPoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}

// WithMinIdleConns 设置最小空闲连接数
func WithMinIdleConns(n int) Option {
	return func(o *Options) {
		o.MinIdleConns = n
	}
}

