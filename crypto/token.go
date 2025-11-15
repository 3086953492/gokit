package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// DefaultTokenAlphabet 默认的令牌字符集
// 包含大小写字母、数字、连字符和下划线，共 64 个字符
// 符合 URL-safe 标准 (RFC 4648 Base64URL)
const DefaultTokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// GenerateRandomString 生成指定长度的安全随机字符串
// 使用 crypto/rand 作为随机源，保证密码学安全
// 参数:
//   - length: 生成字符串的长度，必须大于 0
//   - alphabet: 字符集，不能为空
//
// 返回:
//   - string: 生成的随机字符串
//   - error: 生成过程中的错误
//
// 注意:
//   - 使用拒绝采样确保字符分布均匀
//   - 适用于生成令牌、授权码、会话 ID 等场景
func GenerateRandomString(length int, alphabet string) (string, error) {
	// 验证长度参数
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于 0")
	}

	// 验证字符集参数
	if len(alphabet) == 0 {
		return "", fmt.Errorf("字符集不能为空")
	}

	// 预分配结果字符串的字节切片
	result := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(alphabet)))

	// 生成每个字符
	for i := 0; i < length; i++ {
		// 使用 crypto/rand 生成随机数
		randomIndex, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}
		result[i] = alphabet[randomIndex.Int64()]
	}

	return string(result), nil
}

// GenerateSecureToken 生成指定长度的安全令牌
// 使用默认的 URL-safe 字符集 (A-Z a-z 0-9 - _)
// 参数:
//   - length: 生成令牌的长度，必须大于 0
//
// 返回:
//   - string: 生成的安全令牌
//   - error: 生成过程中的错误
//
// 推荐长度:
//   - 短期令牌: 32 字符
//   - 长期令牌: 64 字符
//   - 授权码: 32-43 字符
func GenerateSecureToken(length int) (string, error) {
	return GenerateRandomString(length, DefaultTokenAlphabet)
}

// GenerateAuthorizationCode 生成 OAuth 授权码
// 默认生成 32 字符长度的授权码，符合 OAuth 2.0 规范建议
// 参数:
//   - length: 授权码长度，必须大于 0，推荐 32-43 字符
//
// 返回:
//   - string: 生成的授权码
//   - error: 生成过程中的错误
//
// 注意:
//   - OAuth 2.0 RFC 6749 建议授权码应具有足够的熵
//   - 授权码应该是一次性使用且有效期较短（通常 10 分钟）
//   - 生成的授权码是 URL-safe 的，可直接用于重定向 URL
func GenerateAuthorizationCode(length int) (string, error) {
	// 验证长度参数
	if length <= 0 {
		return "", fmt.Errorf("授权码长度必须大于 0")
	}

	// 推荐长度校验（警告性质，不阻止执行）
	if length < 32 {
		// 注意：这里不返回错误，只是生成较短的授权码
		// 用户可能有特殊需求，但应该意识到安全性降低
	}

	return GenerateRandomString(length, DefaultTokenAlphabet)
}
