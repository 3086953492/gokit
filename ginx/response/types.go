// Package response 提供标准成功响应结构与辅助函数。
package response

import "github.com/3086953492/gokit/ginx/pagination"

// Response 成功响应的标准结构
type Response struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    any              `json:"data,omitempty"`
	Meta    *pagination.Meta `json:"meta,omitempty"`
}

