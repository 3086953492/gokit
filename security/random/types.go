package random

// Alphabet 预定义的常用字符集
const (
	// AlphabetURLSafe URL 安全字符集（Base64URL），共 64 个字符
	AlphabetURLSafe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	// AlphabetAlphanumeric 字母数字字符集，共 62 个字符
	AlphabetAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// AlphabetHex 十六进制字符集，共 16 个字符
	AlphabetHex = "0123456789abcdef"

	// AlphabetNumeric 纯数字字符集，共 10 个字符
	AlphabetNumeric = "0123456789"

	// AlphabetLowercase 小写字母字符集，共 26 个字符
	AlphabetLowercase = "abcdefghijklmnopqrstuvwxyz"
)

