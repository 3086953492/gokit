package validator

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

// autoRegistry 自动注册器实现
type autoRegistry struct {
	packages []packageInfo
	funcs    map[string]ValidatorFunc
	debug    bool
}

// packageInfo 包信息
type packageInfo struct {
	name   string
	pkg    ValidatorPackage
	config *PackageConfig
}

// NewAutoRegistry 创建自动注册器
func NewAutoRegistry() AutoRegistry {
	return &autoRegistry{
		packages: make([]packageInfo, 0),
		funcs:    make(map[string]ValidatorFunc),
		debug:    false,
	}
}

// RegisterPackage 注册验证器包
func (r *autoRegistry) RegisterPackage(name string, pkg ValidatorPackage, options ...RegisterOption) AutoRegistry {
	config := &PackageConfig{}

	// 应用选项
	for _, option := range options {
		option(config)
	}

	r.packages = append(r.packages, packageInfo{
		name:   name,
		pkg:    pkg,
		config: config,
	})

	return r
}

// RegisterFunc 注册单个验证器函数
func (r *autoRegistry) RegisterFunc(tag string, fn ValidatorFunc) AutoRegistry {
	r.funcs[tag] = fn
	return r
}

// Apply 应用所有注册的验证器
func (r *autoRegistry) Apply() error {
	validator := GetValidator()

	// 注册包级验证器
	for _, pkgInfo := range r.packages {
		if err := r.registerPackage(validator, pkgInfo); err != nil {
			return fmt.Errorf("注册包 %s 失败: %w", pkgInfo.name, err)
		}
	}

	// 注册单个函数
	for tag, fn := range r.funcs {
		if err := validator.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("注册验证器 %s 失败: %w", tag, err)
		}

		if r.debug {
			log.Printf("✓ 注册验证器: %s", tag)
		}
	}

	return nil
}

// registerPackage 注册单个包
func (r *autoRegistry) registerPackage(v *validator.Validate, pkgInfo packageInfo) error {
	// 获取包中的验证器
	validators := pkgInfo.pkg.GetValidators()

	for methodName, validatorFunc := range validators {
		// 检查是否跳过
		if r.shouldSkip(methodName, pkgInfo.config) {
			continue
		}

		// 构建标签名
		tagName := BuildTagName(methodName, pkgInfo.config)

		// 注册验证器
		if err := v.RegisterValidation(tagName, validatorFunc); err != nil {
			return fmt.Errorf("注册方法 %s 失败: %w", methodName, err)
		}

		if r.debug {
			log.Printf("✓ 注册验证器: %s.%s -> %s", pkgInfo.name, methodName, tagName)
		}
	}

	return nil
}

// shouldSkip 检查是否应该跳过某个方法
func (r *autoRegistry) shouldSkip(methodName string, config *PackageConfig) bool {
	if config.Skip == nil {
		return false
	}

	for _, skip := range config.Skip {
		if methodName == skip {
			return true
		}
	}

	return false
}

// EnableDebug 启用调试模式
func (r *autoRegistry) EnableDebug() AutoRegistry {
	r.debug = true
	return r
}
