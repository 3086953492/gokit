package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Manager 提供密码哈希与校验的统一入口
// 默认使用 bcrypt 算法，线程安全
type Manager struct {
	opts *Options
}

// NewManager 创建新的密码管理器实例
// 可通过 Option 函数自定义配置
func NewManager(opts ...Option) (*Manager, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 验证 cost 参数
	if options.Cost < bcrypt.MinCost || options.Cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("%w: 必须在 %d 到 %d 之间", ErrInvalidCost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	// 验证最大密码长度
	if options.MaxPasswordBytes <= 0 {
		return nil, fmt.Errorf("%w: 最大密码长度必须大于 0", ErrInvalidCost)
	}
	if options.MaxPasswordBytes > 72 {
		options.MaxPasswordBytes = 72 // bcrypt 硬限制
	}

	return &Manager{opts: options}, nil
}

// Hash 对明文密码进行哈希处理
// 返回 bcrypt 哈希字符串或错误
func (m *Manager) Hash(password string) (string, error) {
	if err := m.validatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), m.opts.Cost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashFailed, err)
	}

	return string(hashedBytes), nil
}

// Compare 比较明文密码与已存储的哈希值
// 返回 nil 表示匹配，ErrMismatch 表示不匹配，其他错误表示校验过程失败
func (m *Manager) Compare(hashed, password string) error {
	if hashed == "" {
		return ErrHashInvalid
	}

	// 对于 Compare，我们不强制校验密码长度和空值
	// 因为可能是验证旧数据，只校验是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatch
	}

	// 其他错误（如哈希格式无效）
	return fmt.Errorf("%w: %v", ErrHashInvalid, err)
}

// validatePassword 校验密码输入
func (m *Manager) validatePassword(password string) error {
	if !m.opts.AllowEmpty && len(password) == 0 {
		return ErrEmptyPassword
	}

	if len(password) > m.opts.MaxPasswordBytes {
		return fmt.Errorf("%w: 最大 %d 字节", ErrPasswordTooLong, m.opts.MaxPasswordBytes)
	}

	return nil
}

