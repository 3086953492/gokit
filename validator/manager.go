package validator

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	instance *validator.Validate
	once     sync.Once
)

// GetValidator 获取验证器实例（线程安全的懒加载）
func GetValidator() *validator.Validate {
	once.Do(func() {
		instance = validator.New()
	})
	return instance
}

// Register 注册单个自定义验证规则
func Register(tag string, fn ValidatorFunc) error {
	return GetValidator().RegisterValidation(tag, fn)
}

// RegisterBatch 批量注册自定义验证规则
func RegisterBatch(validators map[string]ValidatorFunc) error {
	v := GetValidator()
	for tag, fn := range validators {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("注册验证器 %s 失败: %w", tag, err)
		}
	}
	return nil
}

// ResetValidator 重置验证器实例（主要用于测试）
func ResetValidator() {
	instance = nil
	once = sync.Once{}
}
