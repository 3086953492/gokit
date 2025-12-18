package errors

import (
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger 日志接口，用于解耦 errors 包与具体日志实现
type Logger interface {
	Error(msg string, kv ...any)
}

// nopLogger 默认空日志实现
type nopLogger struct{}

func (nopLogger) Error(string, ...any) {}

var (
	globalLogger Logger = nopLogger{}
	loggerMu     sync.RWMutex
)

// SetLogger 设置全局日志记录器
// 若传入 nil，则使用 nop logger（不记录任何日志）
func SetLogger(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if l == nil {
		globalLogger = nopLogger{}
	} else {
		globalLogger = l
	}
}

// getLogger 获取当前日志记录器
func getLogger() Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return globalLogger
}

// typeToOperation 错误类型到操作名的映射
var typeToOperation = map[string]string{
	TypeDatabase:     "database",
	TypeValidation:   "validation",
	TypeNotFound:     "query",
	TypeInvalidInput: "input",
	TypeUnauthorized: "auth",
	TypeForbidden:    "permission",
	TypeDuplicate:    "duplicate",
	TypeInternal:     "internal",
}

// ErrorBuilder 错误构造器，用于链式构建错误
type ErrorBuilder struct {
	errType string
	message string
	cause   error
	fields  map[string]any
}

// Msg 设置错误消息
func (b *ErrorBuilder) Msg(msg string) *ErrorBuilder {
	b.message = msg
	return b
}

// Err 设置原始错误（错误包装）
func (b *ErrorBuilder) Err(err error) *ErrorBuilder {
	b.cause = err
	return b
}

// Field 添加单个上下文字段（可链式调用多次）
func (b *ErrorBuilder) Field(key string, value any) *ErrorBuilder {
	if b.fields == nil {
		b.fields = make(map[string]any)
	}
	b.fields[key] = value
	return b
}

// Build 构建最终的 AppError
func (b *ErrorBuilder) Build() *AppError {
	return &AppError{
		Type:    b.errType,
		Message: b.message,
		Cause:   b.cause,
		Fields:  b.fields,
	}
}

// Log 构建 AppError 并自动记录日志
func (b *ErrorBuilder) Log() *AppError {
	appErr := b.Build()
	logError(appErr)
	return appErr
}

// getCallerFunctionName 获取调用者函数名
func getCallerFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, "/")
	lastName := parts[len(parts)-1]

	if idx := strings.LastIndex(lastName, "."); idx != -1 {
		return lastName[idx+1:]
	}

	return lastName
}

// logError 记录错误日志
func logError(appErr *AppError) {
	funcName := getCallerFunctionName(3)

	operation := typeToOperation[appErr.Type]
	if operation == "" {
		operation = "unknown"
	}

	// 构建 kv 字段
	kv := []any{
		"function", funcName,
		"operation", operation,
		"type", appErr.Type,
		"timestamp", time.Now(),
	}

	if appErr.Cause != nil {
		kv = append(kv, "error", appErr.Cause)
	}

	for key, value := range appErr.Fields {
		kv = append(kv, key, value)
	}

	getLogger().Error(appErr.Message, kv...)
}

// newBuilder 创建新的错误构造器
func newBuilder(errType, defaultMessage string) *ErrorBuilder {
	return &ErrorBuilder{
		errType: errType,
		message: defaultMessage,
	}
}

// NotFound 创建 404 类型错误（记录不存在）
func NotFound() *ErrorBuilder {
	return newBuilder(TypeNotFound, "记录不存在")
}

// InvalidInput 创建 400 类型错误（输入参数错误）
func InvalidInput() *ErrorBuilder {
	return newBuilder(TypeInvalidInput, "输入参数错误")
}

// Unauthorized 创建 401 类型错误（未授权）
func Unauthorized() *ErrorBuilder {
	return newBuilder(TypeUnauthorized, "未授权")
}

// Forbidden 创建 403 类型错误（权限不足）
func Forbidden() *ErrorBuilder {
	return newBuilder(TypeForbidden, "权限不足")
}

// Duplicate 创建 409 类型错误（数据重复）
func Duplicate() *ErrorBuilder {
	return newBuilder(TypeDuplicate, "数据已存在")
}

// Internal 创建 500 类型错误（内部错误）
func Internal() *ErrorBuilder {
	return newBuilder(TypeInternal, "服务器内部错误")
}

// Database 创建数据库错误
func Database() *ErrorBuilder {
	return newBuilder(TypeDatabase, "数据库操作失败")
}

// Validation 创建 422 类型错误（验证失败）
func Validation() *ErrorBuilder {
	return newBuilder(TypeValidation, "数据验证失败")
}
