package validator

import (
	"github.com/gin-gonic/gin"
)

// 默认的全局验证器实例，用于便捷方法
var defaultValidator = New()

// BindAndValidate Gin 框架 JSON 数据绑定并验证
// 返回验证结果和绑定错误（如果有）
func BindAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	return defaultValidator.BindAndValidate(c, req)
}

// BindQueryAndValidate Gin 框架 Query 参数绑定并验证
func BindQueryAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	return defaultValidator.BindQueryAndValidate(c, req)
}

// BindURIAndValidate Gin 框架 URI 参数绑定并验证
func BindURIAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	return defaultValidator.BindURIAndValidate(c, req)
}

// BindAndValidate JSON 数据绑定并验证（实例方法）
func (v *Validator) BindAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(req); err != nil {
		return newValidationResult(), err
	}

	// 验证数据
	result := v.Validate(req)
	return result, nil
}

// BindQueryAndValidate Query 参数绑定并验证（实例方法）
func (v *Validator) BindQueryAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	// 绑定 Query 参数
	if err := c.ShouldBindQuery(req); err != nil {
		return newValidationResult(), err
	}

	// 验证数据
	result := v.Validate(req)
	return result, nil
}

// BindURIAndValidate URI 参数绑定并验证（实例方法）
func (v *Validator) BindURIAndValidate(c *gin.Context, req interface{}) (*ValidationResult, error) {
	// 绑定 URI 参数
	if err := c.ShouldBindUri(req); err != nil {
		return newValidationResult(), err
	}

	// 验证数据
	result := v.Validate(req)
	return result, nil
}

// SetDefaultValidator 设置默认验证器（用于全局函数）
// 可用于注册全局的自定义验证规则
func SetDefaultValidator(v *Validator) {
	defaultValidator = v
}

// GetDefaultValidator 获取默认验证器
func GetDefaultValidator() *Validator {
	return defaultValidator
}

