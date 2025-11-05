package validator

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

// Init 初始化并注册验证器函数（自动从函数名生成标签）
// 函数名会自动转换为 snake_case 作为标签名
// 例如：UsernameUnique -> username_unique
func Init(validators ...ValidatorFunc) error {
	v := GetValidator()

	for _, fn := range validators {
		// 获取函数名
		funcName := getFuncName(fn)
		if funcName == "" {
			return fmt.Errorf("无法获取函数名")
		}

		// 转换为 snake_case
		tagName := CamelToSnake(funcName)

		// 注册验证器
		if err := v.RegisterValidation(tagName, fn); err != nil {
			return fmt.Errorf("注册验证器 %s 失败: %w", tagName, err)
		}
	}

	return nil
}

// getFuncName 获取函数名
func getFuncName(fn interface{}) string {
	funcValue := reflect.ValueOf(fn)
	if funcValue.Kind() != reflect.Func {
		return ""
	}

	// 获取完整函数名（包含包路径）
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()

	// 提取最后一个部分作为函数名
	// 例如：main.UsernameUnique -> UsernameUnique
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

// ResetValidator 重置验证器实例（主要用于测试）
func ResetValidator() {
	instance = nil
	once = sync.Once{}
}
