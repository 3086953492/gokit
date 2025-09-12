package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
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

// RegisterValidation 注册自定义验证规则
func RegisterValidation(tag string, fn validator.Func) error {
	return GetValidator().RegisterValidation(tag, fn)
}

// ResetValidator 重置验证器实例（主要用于测试）
func ResetValidator() {
	instance = nil
}

// RegisterPackageValidators 注册包级验证器（便捷方法）
func RegisterPackageValidators(pkg ValidatorPackage, options ...RegisterOption) error {
	return NewAutoRegistry().
		RegisterPackage("default", pkg, options...).
		Apply()
}

// RegisterMultiplePackages 注册多个包
func RegisterMultiplePackages(packages map[string]ValidatorPackage, globalOptions ...RegisterOption) error {
	registry := NewAutoRegistry()

	for name, pkg := range packages {
		registry.RegisterPackage(name, pkg, globalOptions...)
	}

	return registry.Apply()
}
