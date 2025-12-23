package subject

// Options 定义 Manager 的配置选项
type Options struct {
	// Secret HMAC 密钥，必须提供且建议至少 32 字节
	Secret []byte

	// Length 输出 sub 的长度（不含前缀），0 表示不截断（完整 43 字符）
	Length int

	// Prefix sub 的前缀，例如 "u_"，便于区分来源
	Prefix string

	// AllowShortSecret 是否允许短密钥（小于 32 字节），默认 false
	AllowShortSecret bool
}

// Option 是配置选项函数类型
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		Secret:           nil,
		Length:           0, // 不截断
		Prefix:           "",
		AllowShortSecret: false,
	}
}

// WithSecret 设置 HMAC 密钥（字节切片）
func WithSecret(secret []byte) Option {
	return func(o *Options) {
		o.Secret = secret
	}
}

// WithSecretString 设置 HMAC 密钥（字符串形式）
func WithSecretString(secret string) Option {
	return func(o *Options) {
		o.Secret = []byte(secret)
	}
}

// WithLength 设置输出 sub 的长度（不含前缀）
// 0 表示不截断，使用完整的 43 字符
// 有效范围：1-43
func WithLength(n int) Option {
	return func(o *Options) {
		o.Length = n
	}
}

// WithPrefix 设置 sub 的前缀
// 例如 "u_" 会使输出变为 "u_xxxxxxx..."
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// WithAllowShortSecret 允许使用短密钥（不推荐，仅用于测试）
func WithAllowShortSecret(allow bool) Option {
	return func(o *Options) {
		o.AllowShortSecret = allow
	}
}

