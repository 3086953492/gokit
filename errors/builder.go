package errors

import (
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/3086953492/YaBase/logger"
)

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
	// 去除包路径，只保留函数名
	// 例如：github.com/user/repo/service.CreateExpense -> CreateExpense
	parts := strings.Split(fullName, "/")
	lastName := parts[len(parts)-1]

	// 去除包名，只保留函数名
	if idx := strings.LastIndex(lastName, "."); idx != -1 {
		return lastName[idx+1:]
	}

	return lastName
}

// logError 记录错误日志
func logError(appErr *AppError) {
	// 获取调用函数名（跳过 logError -> Log -> 用户代码，所以 skip=3）
	funcName := getCallerFunctionName(3)

	// 获取 operation
	operation := typeToOperation[appErr.Type]
	if operation == "" {
		operation = "unknown"
	}

	// 构建 zap 字段
	baseFields := []zap.Field{
		zap.String("function", funcName),
		zap.String("operation", operation),
		zap.String("type", appErr.Type),
		zap.Time("timestamp", time.Now()),
	}

	// 添加原始错误
	if appErr.Cause != nil {
		baseFields = append(baseFields, zap.Error(appErr.Cause))
	}

	// 添加上下文字段
	for key, value := range appErr.Fields {
		baseFields = append(baseFields, zap.Any(key, value))
	}

	// 记录日志
	logErrorToLogger(appErr.Message, baseFields)
}

// logErrorToLogger 实际记录日志
func logErrorToLogger(message string, fields []zap.Field) {
	defer func() {
		if r := recover(); r != nil {
			// 如果 logger 未初始化，忽略错误
		}
	}()

	logger.Error(message, fields...)
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
