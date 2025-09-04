package errors

import "fmt"

// ================================
// 用户管理相关错误
// ================================

// 用户相关错误类型常量
const (
	TypeUserNotFound            = "USER_NOT_FOUND"
	TypeUserExists              = "USER_EXISTS"
	TypeUserDisabled            = "USER_DISABLED"
	TypeOperationFailed         = "OPERATION_FAILED"
	TypeUsernameOrPasswordError = "USERNAME_OR_PASSWORD_ERROR"
	TypeUsernameExists          = "USERNAME_EXISTS"
	TypeUserRegisterFailed      = "USER_REGISTER_FAILED"
	TypeUserListNotFound        = "USER_LIST_NOT_FOUND"
)

// 用户相关错误实例
var (
	ErrUserNotFound            = New(TypeUserNotFound, "用户不存在")
	ErrUserExists              = New(TypeUserExists, "用户已存在")
	ErrUserDisabled            = New(TypeUserDisabled, "账户已被禁用")
	ErrCreateFailed            = New(TypeOperationFailed, "创建失败")
	ErrUpdateFailed            = New(TypeOperationFailed, "更新失败")
	ErrDeleteFailed            = New(TypeOperationFailed, "删除失败")
	ErrUsernameOrPasswordError = New(TypeUsernameOrPasswordError, "用户名或密码错误")
	ErrUsernameExists          = New(TypeUsernameExists, "用户名已存在")
	ErrUserRegisterFailed      = New(TypeUserRegisterFailed, "用户注册失败")
	ErrUserListNotFound        = New(TypeUserListNotFound, "用户列表不存在")
)

// 用户错误工厂函数
func NewUserNotFoundError(userID uint) *AppError {
	return New(TypeUserNotFound, fmt.Sprintf("用户不存在: ID=%d", userID))
}

func NewValidationError(field, message string) *AppError {
	return New(TypeInvalidInput, fmt.Sprintf("%s: %s", field, message))
}
