package json

// AllInStringList 检查所有给定的值是否都存在于列表中
// 参数:
//   - list: 目标字符串列表
//   - values: 需要检查的一个或多个字符串值
//
// 返回:
//   - bool: 当所有 values 都在 list 中时返回 true，否则返回 false
//
// 行为约定:
//   - 匹配规则为大小写敏感的字符串全等比较
//   - 当 values 为空时返回 true（没有需要校验的值视为通过）
//   - 当 list 为空且 values 非空时返回 false
//
// 使用示例:
//
//	grantTypes := []string{"authorization_code", "client_credentials", "refresh_token"}
//	if AllInStringList(grantTypes, "authorization_code", "refresh_token") {
//	    // 所有值都在列表中
//	}
func AllInStringList(list []string, values ...string) bool {
	// 当没有需要检查的值时，视为通过
	if len(values) == 0 {
		return true
	}

	// 当列表为空但有需要检查的值时，返回 false
	if len(list) == 0 {
		return false
	}

	// 将列表转换为 map 以提高查找效率
	listMap := make(map[string]struct{}, len(list))
	for _, item := range list {
		listMap[item] = struct{}{}
	}

	// 检查所有值是否都在 map 中
	for _, value := range values {
		if _, exists := listMap[value]; !exists {
			return false
		}
	}

	return true
}
