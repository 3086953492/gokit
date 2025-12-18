package response

// --------------------------------------------------------------------
// Writer 接口：框架无关的响应写入抽象
// --------------------------------------------------------------------

// JSONWriter 写入 JSON 响应。
type JSONWriter interface {
	// JSON 写入 JSON 响应，status 为 HTTP 状态码，body 为响应体。
	JSON(status int, body any)
}

// RedirectWriter 执行 HTTP 重定向。
type RedirectWriter interface {
	// Redirect 执行重定向，status 为 301/302/307/308，location 为目标地址。
	Redirect(status int, location string)
}

// ResponseWriter 综合写入能力，供 provider 实现。
type ResponseWriter interface {
	JSONWriter
	RedirectWriter
}

// --------------------------------------------------------------------
// response.Error 接口：自有的可响应错误抽象（脱离 errors 包依赖）
// --------------------------------------------------------------------

// Error 描述一个可用于响应的错误。
// 调用方可实现此接口，也可以使用 NewError 快速构造。
type Error interface {
	error
	// Code 稳定的业务错误码，用于前端/调用方分支判断。
	Code() string
	// Status 建议的 HTTP 状态码；若为 0 表示未指定，Manager 会使用 fallback。
	Status() int
	// Title RFC7807 title（简短描述）；若为空 Manager 会使用 Code。
	Title() string
	// Detail RFC7807 detail（更长的可暴露信息）；可为空。
	Detail() string
	// ProblemType RFC7807 type（URI 引用）；若为空 Manager 自动生成 urn。
	ProblemType() string
}

// --------------------------------------------------------------------
// RFC7807 Problem 结构（仅用于错误响应）
// --------------------------------------------------------------------

// Problem 表示 RFC7807 problem+json 响应体。
// 根据计划，不携带 details 扩展字段。
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// --------------------------------------------------------------------
// 成功响应结构（不走 RFC7807）
// --------------------------------------------------------------------

// SuccessBody 通用成功响应结构。
type SuccessBody struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// PaginatedData 分页数据结构。
type PaginatedData struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

// PaginatedBody 分页成功响应结构。
type PaginatedBody struct {
	Message string        `json:"message,omitempty"`
	Data    PaginatedData `json:"data"`
}

// --------------------------------------------------------------------
// simpleError：Error 接口的简单实现
// --------------------------------------------------------------------

// simpleError 是 Error 接口的简单实现。
type simpleError struct {
	code        string
	status      int
	title       string
	detail      string
	problemType string
}

func (e *simpleError) Error() string {
	if e.detail != "" {
		return e.detail
	}
	if e.title != "" {
		return e.title
	}
	return e.code
}

func (e *simpleError) Code() string        { return e.code }
func (e *simpleError) Status() int         { return e.status }
func (e *simpleError) Title() string       { return e.title }
func (e *simpleError) Detail() string      { return e.detail }
func (e *simpleError) ProblemType() string { return e.problemType }

// NewError 创建一个实现 Error 接口的简单错误。
// code 为必填；其余字段可选（空字符串/0 表示未指定）。
func NewError(code string, status int, title, detail, problemType string) Error {
	return &simpleError{
		code:        code,
		status:      status,
		title:       title,
		detail:      detail,
		problemType: problemType,
	}
}

