package json

import (
	"encoding/json"
	"fmt"
)

// IsSubset 判断 subJSON 是否为 supJSON 的子集
// subJSON 和 supJSON 应该是 JSON 数组字符串，例如: `["authorization_code", "client_credentials"]`
// 返回判断结果和可能的错误
func IsSubset(subJSON string, supJSON string) error {
	var sub []string
	var sup []string

	// 解析 subJSON
	if err := json.Unmarshal([]byte(subJSON), &sub); err != nil {
		return fmt.Errorf("解析 sub JSON 失败: %w", err)
	}

	// 解析 supJSON
	if err := json.Unmarshal([]byte(supJSON), &sup); err != nil {
		return fmt.Errorf("解析 sup JSON 失败: %w", err)
	}

	// 将 sup 转换为 map 以提高查找效率
	supSet := make(map[string]bool)
	for _, item := range sup {
		supSet[item] = true
	}

	// 检查 sub 中的每个元素是否都在 supSet 中
	for _, item := range sub {
		if !supSet[item] {
			return fmt.Errorf("sub 中的元素 %s 不在 sup 中", item)
		}
	}

	return nil
}
