package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// BindAndValidate Gin框架数据绑定并验证（推荐使用，替代ValidateStruct）
func BindAndValidate(c *gin.Context, req any) error {
	// 绑定JSON数据
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("数据绑定失败: %w", err)
	}

	// 验证数据
	if err := GetValidator().Struct(req); err != nil {
		return fmt.Errorf("数据验证失败: %w", err)
	}

	return nil
}
