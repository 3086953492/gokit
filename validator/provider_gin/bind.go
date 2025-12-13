// Package provider_gin 提供 Gin 框架的验证器绑定支持
package provider_gin

import (
	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/validator"
)

// BindAndValidate JSON 数据绑定并验证
// 返回验证结果和绑定错误（如果有）
func BindAndValidate(m *validator.Manager, c *gin.Context, req any) (*validator.Result, error) {
	if err := c.ShouldBindJSON(req); err != nil {
		return nil, err
	}
	return m.Validate(req), nil
}

// BindQueryAndValidate Query 参数绑定并验证
func BindQueryAndValidate(m *validator.Manager, c *gin.Context, req any) (*validator.Result, error) {
	if err := c.ShouldBindQuery(req); err != nil {
		return nil, err
	}
	return m.Validate(req), nil
}

// BindURIAndValidate URI 参数绑定并验证
func BindURIAndValidate(m *validator.Manager, c *gin.Context, req any) (*validator.Result, error) {
	if err := c.ShouldBindUri(req); err != nil {
		return nil, err
	}
	return m.Validate(req), nil
}

// BindFormAndValidate Form 表单数据绑定并验证
func BindFormAndValidate(m *validator.Manager, c *gin.Context, req any) (*validator.Result, error) {
	if err := c.ShouldBind(req); err != nil {
		return nil, err
	}
	return m.Validate(req), nil
}
