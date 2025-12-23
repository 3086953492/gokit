package subject

// MaxSubLength HMAC-SHA256 输出经 base64url 编码后的最大长度（不含前缀）
// SHA256 输出 32 字节，base64url 编码后为 43 字符（无 padding）
const MaxSubLength = 43

// MinSecretLength 建议的最小密钥长度（字节）
const MinSecretLength = 32

