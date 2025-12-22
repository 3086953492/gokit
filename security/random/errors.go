package random

import "errors"

// 随机字符串生成相关错误变量
var (
	// ErrInvalidLength 表示长度参数无效
	ErrInvalidLength = errors.New("长度必须大于 0")

	// ErrEmptyAlphabet 表示字符集为空
	ErrEmptyAlphabet = errors.New("字符集不能为空")

	// ErrAlphabetTooLarge 表示字符集长度超过限制
	ErrAlphabetTooLarge = errors.New("字符集长度超过 256")

	// ErrReadFailed 表示从随机源读取失败
	ErrReadFailed = errors.New("随机源读取失败")
)

