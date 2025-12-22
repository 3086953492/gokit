package random

import (
	"fmt"
)

// Generator 提供密码学安全的随机字符串/字节生成
// 使用 crypto/rand 作为默认随机源，线程安全
type Generator struct {
	opts *Options
}

// NewGenerator 创建新的随机生成器实例
// 可通过 Option 函数自定义字符集和随机源
func NewGenerator(opts ...Option) (*Generator, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 验证字符集
	if len(options.Alphabet) == 0 {
		return nil, ErrEmptyAlphabet
	}
	if len(options.Alphabet) > 256 {
		return nil, ErrAlphabetTooLarge
	}

	// 验证随机源
	if options.Reader == nil {
		return nil, fmt.Errorf("随机源不能为 nil")
	}

	return &Generator{opts: options}, nil
}

// String 生成指定长度的随机字符串
// 使用 Generator 配置的字符集
func (g *Generator) String(n int) (string, error) {
	if n <= 0 {
		return "", ErrInvalidLength
	}
	return generateString(g.opts.Reader, n, g.opts.Alphabet)
}

// StringWithAlphabet 使用指定字符集生成随机字符串
// 忽略 Generator 配置的默认字符集
func (g *Generator) StringWithAlphabet(n int, alphabet string) (string, error) {
	if n <= 0 {
		return "", ErrInvalidLength
	}
	if len(alphabet) == 0 {
		return "", ErrEmptyAlphabet
	}
	return generateString(g.opts.Reader, n, alphabet)
}

// Bytes 生成指定长度的随机字节
func (g *Generator) Bytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidLength
	}
	return generateBytes(g.opts.Reader, n)
}

// URLSafe 生成 URL 安全的随机字符串（便捷函数）
// 使用默认随机源和 URL-safe 字符集
func URLSafe(n int) (string, error) {
	if n <= 0 {
		return "", ErrInvalidLength
	}
	return generateString(defaultOptions().Reader, n, AlphabetURLSafe)
}

// Hex 生成十六进制随机字符串（便捷函数）
// 使用默认随机源和十六进制字符集
func Hex(n int) (string, error) {
	if n <= 0 {
		return "", ErrInvalidLength
	}
	return generateString(defaultOptions().Reader, n, AlphabetHex)
}
