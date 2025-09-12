package validator

import "github.com/go-playground/validator/v10"

// ValidatorPackage 验证器包接口
type ValidatorPackage interface {
	// GetValidators 返回包中所有验证器方法
	// 约定：方法名即为验证器名称，会自动转换为snake_case
	GetValidators() map[string]ValidatorFunc
}

// ValidatorFunc 验证器函数类型
type ValidatorFunc = validator.Func

// AutoRegistry 自动注册器接口
type AutoRegistry interface {
	// RegisterPackage 注册验证器包
	RegisterPackage(name string, pkg ValidatorPackage, options ...RegisterOption) AutoRegistry

	// RegisterFunc 注册单个验证器函数
	RegisterFunc(tag string, fn ValidatorFunc) AutoRegistry

	// Apply 应用所有注册的验证器
	Apply() error
}

// RegisterOption 注册选项
type RegisterOption func(*PackageConfig)

// PackageConfig 包配置
type PackageConfig struct {
	Prefix    string            // 标签前缀
	Tags      map[string]string // 自定义标签映射
	Skip      []string          // 跳过的方法名
	Transform TagTransformer    // 自定义标签转换器
}

// TagTransformer 标签转换器
type TagTransformer func(methodName string) string

type FieldLevel = validator.FieldLevel
