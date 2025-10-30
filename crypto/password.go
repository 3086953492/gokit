package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 算法加密密码
// 使用默认的加密强度 (cost=10)
// 参数:
//   - password: 需要加密的明文密码
//
// 返回:
//   - string: 加密后的密码哈希值
//   - error: 加密过程中的错误
func HashPassword(password string) (string, error) {
	return HashPasswordWithCost(password, bcrypt.DefaultCost)
}

// HashPasswordWithCost 使用指定的加密强度加密密码
// 参数:
//   - password: 需要加密的明文密码
//   - cost: 加密强度，范围 4-31，值越大越安全但计算时间越长
//
// 返回:
//   - string: 加密后的密码哈希值
//   - error: 加密过程中的错误
//
// 注意:
//   - bcrypt 最大支持 72 字节的密码
//   - 推荐 cost 值: 10-12 (默认为 10)
func HashPasswordWithCost(password string, cost int) (string, error) {
	// 验证 cost 参数
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("加密强度必须在 %d 到 %d 之间", bcrypt.MinCost, bcrypt.MaxCost)
	}

	// 验证密码长度
	if len(password) == 0 {
		return "", fmt.Errorf("密码不能为空")
	}

	if len(password) > 72 {
		return "", fmt.Errorf("密码长度不能超过 72 字节")
	}

	// 生成密码哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码是否匹配
// 参数:
//   - hashedPassword: 存储的密码哈希值
//   - password: 需要验证的明文密码
//
// 返回:
//   - bool: true 表示密码匹配，false 表示密码不匹配
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
