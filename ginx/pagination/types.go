// Package pagination 提供分页参数解析与元数据构建。
package pagination

// Page 分页参数（由 Parse 解析得到）
type Page struct {
	Page     int // 当前页码（从 1 开始）
	PageSize int // 每页条数
	Offset   int // 偏移量（用于 SQL OFFSET）
	Limit    int // 等于 PageSize，方便使用
}

// Meta 分页响应元数据
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta 根据 Page 和总记录数创建 Meta
func NewMeta(p Page, total int64) Meta {
	totalPages := 0
	if p.PageSize > 0 {
		totalPages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
	return Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

