// Package problem 提供 RFC 7807 Problem Details 响应。
package problem

// Problem RFC 7807 Problem Details 结构
type Problem struct {
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Status     int            `json:"status"`
	Detail     string         `json:"detail,omitempty"`
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"-"` // 扩展字段，序列化时展平到顶层
}

// ContentType RFC 7807 规定的 Content-Type
const ContentType = "application/problem+json"

