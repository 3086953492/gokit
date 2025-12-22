package password

// Hasher 定义密码哈希与校验的接口
type Hasher interface {
	// Hash 对明文密码进行哈希处理
	// 返回哈希后的字符串或错误
	Hash(password string) (string, error)

	// Compare 比较明文密码与已存储的哈希值
	// 返回 nil 表示匹配，ErrMismatch 表示不匹配，其他错误表示校验过程失败
	Compare(hashed, password string) error
}

