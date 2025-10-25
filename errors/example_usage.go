package errors

// 本文件提供使用示例，展示如何使用 errors 包的各种功能

import (
	"fmt"
	"io"
)

// ExampleSimpleError 演示创建简单错误
func ExampleSimpleError() error {
	// 使用标准库函数创建简单错误
	return NewSimple("这是一个简单错误")
}

// ExampleAppError 演示创建结构化的 AppError
func ExampleAppError() error {
	// 使用链式 API 创建 AppError
	return NotFound().
		WithMessage("用户不存在").
		WithField("user_id", 123).
		Build()
}

// ExampleErrorWrapping 演示错误包装
func ExampleErrorWrapping() error {
	// 模拟一个底层错误
	originalErr := NewSimple("数据库连接超时")

	// 使用 AppError 包装原始错误
	return Database().
		WithMessage("查询用户失败").
		WithCause(originalErr).
		WithField("user_id", 456).
		WithField("retry_count", 3).
		Build()
}

// ExampleErrorChecking 演示错误检查
func ExampleErrorChecking() {
	err := NotFound().WithMessage("资源不存在").Build()

	// 检查 AppError 类型
	if GetType(err) == TypeNotFound {
		fmt.Println("这是一个 NotFound 错误")
	}

	// 使用 IsAppError 检查
	notFoundErr := NotFound().Build()
	if IsAppError(err, notFoundErr) {
		fmt.Println("错误类型匹配")
	}

	// 使用标准库 Is 检查
	if Is(err, io.EOF) {
		fmt.Println("这是 EOF 错误")
	}
}

// ExampleMultipleErrors 演示合并多个错误
func ExampleMultipleErrors() error {
	err1 := NewSimple("错误 1")
	err2 := NewSimple("错误 2")
	err3 := NotFound().WithMessage("错误 3").Build()

	// 使用标准库 Join 合并错误
	return Join(err1, err2, err3)
}

// ExampleErrorFields 演示使用错误字段
func ExampleErrorFields() {
	err := Internal().
		WithMessage("处理订单失败").
		WithFields(map[string]any{
			"order_id":   "ORDER-12345",
			"user_id":    789,
			"amount":     99.99,
			"step":       "payment",
			"error_code": 5001,
		}).
		Build()

	// 获取所有字段
	fields := GetFields(err)
	fmt.Printf("错误字段: %+v\n", fields)

	// 获取单个字段
	var appErr *AppError
	if As(err, &appErr) {
		if orderID, exists := appErr.GetField("order_id"); exists {
			fmt.Printf("订单号: %v\n", orderID)
		}

		// 检查字段是否存在
		if appErr.HasField("amount") {
			fmt.Println("包含金额信息")
		}
	}
}

// ExampleDatabaseError 演示数据库错误转换
func ExampleDatabaseError() error {
	// 模拟 GORM 错误
	// dbErr := gorm.ErrRecordNotFound

	// 自动转换数据库错误
	// return FromDatabaseError(dbErr)

	// 手动创建数据库错误
	return Database().
		WithMessage("数据库查询失败").
		WithField("table", "users").
		WithField("query", "SELECT * FROM users WHERE id = ?").
		Build()
}

// ExampleValidationError 演示验证错误
func ExampleValidationError() error {
	return Validation().
		WithMessage("用户名格式不正确").
		WithField("field", "username").
		WithField("value", "ab").
		WithField("rule", "minLength:3").
		Build()
}

// ExampleUnauthorizedError 演示未授权错误
func ExampleUnauthorizedError() error {
	return Unauthorized().
		WithMessage("令牌已过期").
		WithField("token_type", "JWT").
		WithField("expired_at", "2024-01-01T00:00:00Z").
		Build()
}

// ExampleComplexScenario 演示复杂场景
func ExampleComplexScenario() error {
	// 模拟多层错误包装
	originalErr := NewSimple("网络连接失败")

	dbErr := Database().
		WithMessage("无法连接到数据库").
		WithCause(originalErr).
		WithField("host", "localhost:3306").
		Build()

	return Internal().
		WithMessage("用户服务不可用").
		WithCause(dbErr).
		WithFields(map[string]any{
			"service":    "UserService",
			"endpoint":   "/api/users/123",
			"timestamp":  "2024-01-01T12:00:00Z",
			"request_id": "req-abc-123",
		}).
		Build()
}
