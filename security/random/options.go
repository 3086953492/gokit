package random

import (
	"crypto/rand"
	"io"
)

// Options 定义 Generator 的配置选项
type Options struct {
	// Alphabet 字符集，默认为 URL-safe Base64
	Alphabet string

	// Reader 随机源，默认为 crypto/rand.Reader
	Reader io.Reader
}

// Option 是配置选项函数类型
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		Alphabet: AlphabetURLSafe,
		Reader:   rand.Reader,
	}
}

// WithAlphabet 设置字符集
func WithAlphabet(alphabet string) Option {
	return func(o *Options) {
		o.Alphabet = alphabet
	}
}

// WithReader 设置随机源
// 可用于测试时注入可预测的随机源
func WithReader(r io.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

