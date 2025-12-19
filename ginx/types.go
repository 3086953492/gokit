package ginx

// Response 成功响应的标准结构
type Response struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
	Meta    *PageMeta `json:"meta,omitempty"`
}

// Page 分页参数（由 ParsePage 解析得到）
type Page struct {
	Page     int // 当前页码（从 1 开始）
	PageSize int // 每页条数
	Offset   int // 偏移量（用于 SQL OFFSET）
	Limit    int // 等于 PageSize，方便使用
}

// PageMeta 分页响应元数据
type PageMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPageMeta 根据 Page 和总记录数创建 PageMeta
func NewPageMeta(p Page, total int64) PageMeta {
	totalPages := 0
	if p.PageSize > 0 {
		totalPages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
	return PageMeta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// Problem RFC 7807 Problem Details 结构
type Problem struct {
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Status     int            `json:"status"`
	Detail     string         `json:"detail,omitempty"`
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"-"` // 扩展字段，序列化时展平到顶层
}

// ContentTypeProblem RFC 7807 规定的 Content-Type
const ContentTypeProblem = "application/problem+json"
