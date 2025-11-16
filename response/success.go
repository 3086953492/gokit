package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginatedData 分页数据结构
type PaginatedData struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

// Success 返回成功响应
// c: gin 上下文
// message: 成功消息，如果为空则使用配置中的默认消息
// data: 响应数据
func Success(c *gin.Context, message string, data any) {
	cfg := getConfig()

	if message == "" {
		message = cfg.DefaultSuccessMessage
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Paginated 返回分页响应
// c: gin 上下文
// data: 分页数据项
// total: 总记录数
// page: 当前页码
// pageSize: 每页大小
func Paginated(c *gin.Context, data any, total int64, page, pageSize int) {
	cfg := getConfig()

	// 计算总页数，向上取整
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: cfg.DefaultPaginatedMessage,
		Data: PaginatedData{
			Items:      data,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}
