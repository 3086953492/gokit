package subject

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Manager 提供 OIDC/OAuth subject 标识符生成
// 使用 HMAC-SHA256 + Base64URL 生成稳定、不可逆的 sub
// 线程安全
type Manager struct {
	opts *Options
}

// NewManager 创建新的 subject 管理器实例
// 必须通过 WithSecret 或 WithSecretString 提供密钥
func NewManager(opts ...Option) (*Manager, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 验证密钥
	if len(options.Secret) == 0 {
		return nil, ErrEmptySecret
	}
	if !options.AllowShortSecret && len(options.Secret) < MinSecretLength {
		return nil, fmt.Errorf("%w: 当前 %d 字节", ErrSecretTooShort, len(options.Secret))
	}

	// 验证截断长度
	if options.Length < 0 || options.Length > MaxSubLength {
		return nil, fmt.Errorf("%w: 必须在 0-%d 之间", ErrInvalidLength, MaxSubLength)
	}

	return &Manager{opts: options}, nil
}

// Sub 根据用户 ID 生成稳定的 subject 标识符
// 输出格式：[prefix]base64url(HMAC-SHA256(secret, userID))[:length]
// 同一 userID + 同一 secret 始终生成相同的 sub
func (m *Manager) Sub(userID string) (string, error) {
	if userID == "" {
		return "", ErrEmptyUserID
	}

	// 计算 HMAC-SHA256
	h := hmac.New(sha256.New, m.opts.Secret)
	h.Write([]byte(userID))
	hash := h.Sum(nil)

	// Base64URL 编码（无 padding）
	encoded := base64.RawURLEncoding.EncodeToString(hash)

	// 可选截断
	if m.opts.Length > 0 && m.opts.Length < len(encoded) {
		encoded = encoded[:m.opts.Length]
	}

	// 添加前缀
	if m.opts.Prefix != "" {
		return m.opts.Prefix + encoded, nil
	}

	return encoded, nil
}

// SubWithSector 根据用户 ID 和 sector 生成 pairwise subject 标识符
// 适用于需要跨客户端隔离的场景（Pairwise Subject Identifier）
// 输出格式：[prefix]base64url(HMAC-SHA256(secret, userID + ":" + sector))[:length]
func (m *Manager) SubWithSector(userID, sector string) (string, error) {
	if userID == "" {
		return "", ErrEmptyUserID
	}
	// sector 可以为空，表示退化为 Public 模式
	if sector == "" {
		return m.Sub(userID)
	}

	// 计算 HMAC-SHA256，将 sector 纳入输入
	h := hmac.New(sha256.New, m.opts.Secret)
	h.Write([]byte(userID))
	h.Write([]byte(":"))
	h.Write([]byte(sector))
	hash := h.Sum(nil)

	// Base64URL 编码（无 padding）
	encoded := base64.RawURLEncoding.EncodeToString(hash)

	// 可选截断
	if m.opts.Length > 0 && m.opts.Length < len(encoded) {
		encoded = encoded[:m.opts.Length]
	}

	// 添加前缀
	if m.opts.Prefix != "" {
		return m.opts.Prefix + encoded, nil
	}

	return encoded, nil
}

