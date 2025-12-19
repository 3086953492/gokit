package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Parse 从 *gin.Context 解析分页参数并归一化。
//
// 默认 page=1, page_size=10, 最大 page_size=100。可通过 Option 覆盖。
//
//	p := pagination.Parse(c)
//	p := pagination.Parse(c, pagination.WithMaxPageSize(50))
func Parse(c *gin.Context, opts ...Option) Page {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	page := parseIntQuery(c, "page", o.DefaultPage)
	if page < 1 {
		page = o.DefaultPage
	}

	pageSize := parseIntQuery(c, "page_size", o.DefaultPageSize)
	if pageSize < 1 {
		pageSize = o.DefaultPageSize
	}
	if pageSize > o.MaxPageSize {
		pageSize = o.MaxPageSize
	}

	offset := (page - 1) * pageSize

	return Page{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// parseIntQuery 从 query 读取整数，失败或不存在返回 defaultVal
func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	s := c.Query(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

