package errors

import "maps"

// ErrorBuilder 错误构造器，用于链式构建错误
type ErrorBuilder struct {
	errType string
	message string
	cause   error
	fields  map[string]any
}

// WithMessage 设置错误消息
func (b *ErrorBuilder) WithMessage(msg string) *ErrorBuilder {
	b.message = msg
	return b
}

// WithCause 设置原始错误（错误包装）
func (b *ErrorBuilder) WithCause(err error) *ErrorBuilder {
	b.cause = err
	return b
}

// WithField 添加单个上下文字段
func (b *ErrorBuilder) WithField(key string, value any) *ErrorBuilder {
	if b.fields == nil {
		b.fields = make(map[string]any)
	}
	b.fields[key] = value
	return b
}

// WithFields 批量添加上下文字段
func (b *ErrorBuilder) WithFields(fields map[string]any) *ErrorBuilder {
	if b.fields == nil {
		b.fields = make(map[string]any)
	}
	maps.Copy(b.fields, fields)
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
