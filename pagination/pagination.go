package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPageAndPageSize(c *gin.Context) (int, int) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}