package random

import (
	"fmt"
	"io"
)

// generateString 使用 rejection sampling 生成随机字符串
// 确保字符分布均匀，避免 math/big 的性能开销
func generateString(r io.Reader, length int, alphabet string) (string, error) {
	alphabetLen := len(alphabet)
	if alphabetLen == 0 {
		return "", ErrEmptyAlphabet
	}
	if alphabetLen > 256 {
		return "", ErrAlphabetTooLarge
	}

	// 计算 rejection sampling 的阈值
	// 为了保证均匀分布，我们需要拒绝 >= maxValid 的值
	// 例如：alphabet 长度 64，则 maxValid = 256 - (256 % 64) = 256
	// 例如：alphabet 长度 62，则 maxValid = 256 - (256 % 62) = 248
	maxValid := 256 - (256 % alphabetLen)

	result := make([]byte, length)
	// 预分配缓冲区，减少读取次数
	// 考虑到 rejection 可能需要多次读取，预分配 1.5 倍空间
	bufSize := length + length/2
	if bufSize < 16 {
		bufSize = 16
	}
	buf := make([]byte, bufSize)

	generated := 0
	for generated < length {
		// 批量读取随机字节
		n, err := r.Read(buf)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrReadFailed, err)
		}

		// 处理读取的字节
		for i := 0; i < n && generated < length; i++ {
			// rejection sampling：拒绝超出阈值的值
			if int(buf[i]) < maxValid {
				result[generated] = alphabet[int(buf[i])%alphabetLen]
				generated++
			}
		}
	}

	return string(result), nil
}

// generateBytes 生成指定长度的随机字节
func generateBytes(r io.Reader, length int) ([]byte, error) {
	result := make([]byte, length)
	_, err := io.ReadFull(r, result)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFailed, err)
	}
	return result, nil
}

