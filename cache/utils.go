package cache

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// normalizeValue 将任意值标准化为字符串，保证相同值的不同类型表示能得到相同结果
func normalizeValue(v any) string {
	if v == nil {
		return "nil"
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.String:
		return val.String()
	case reflect.Slice, reflect.Array:
		// 处理切片：递归标准化每个元素，用逗号连接
		parts := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			parts[i] = normalizeValue(val.Index(i).Interface())
		}
		return "[" + strings.Join(parts, ",") + "]"
	case reflect.Map:
		return fmt.Sprint(v)
	default:
		// 其他类型（结构体等）兜底使用 fmt.Sprint
		return fmt.Sprint(v)
	}
}

// BuildKey 通用的 key 构造函数
// 使用方式：BuildKey("prefix", part1, part2, ...) => "prefix|part1|part2|..."
func BuildKey(prefix string, parts ...any) string {
	if len(parts) == 0 {
		return prefix
	}

	normalizedParts := make([]string, len(parts))
	for i, part := range parts {
		normalizedParts[i] = normalizeValue(part)
	}

	return prefix + "|" + strings.Join(normalizedParts, "|")
}

