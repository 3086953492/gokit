package errors

// ================================
// 认证和授权相关错误
// ================================

// 认证相关错误类型常量
const (
	TypeInvalidCredentials   = "INVALID_CREDENTIALS"
	TypeTokenExpired         = "TOKEN_EXPIRED"
	TypeTokenInvalid         = "TOKEN_INVALID"
	TypeUnauthorized         = "UNAUTHORIZED"
	TypePermissionDenied     = "PERMISSION_DENIED"
	TypeTokenRefreshTooEarly = "TOKEN_REFRESH_TOO_EARLY"
	TypeTokenGenerateFailed  = "TOKEN_GENERATE_FAILED"
)

// 认证相关错误实例
var (
	ErrInvalidCredentials   = New(TypeInvalidCredentials, "用户名或密码错误")
	ErrTokenExpired         = New(TypeTokenExpired, "登录已过期，请重新登录")
	ErrTokenInvalid         = New(TypeTokenInvalid, "无效的登录凭证")
	ErrUnauthorized         = New(TypeUnauthorized, "请先登录")
	ErrPermissionDenied     = New(TypePermissionDenied, "权限不足")
	ErrTokenRefreshTooEarly = New(TypeTokenRefreshTooEarly, "刷新令牌过于频繁")
	ErrTokenGenerateFailed  = New(TypeTokenGenerateFailed, "生成token失败")
)

// 认证错误工厂函数
func NewInvalidCredentialsError(reason string) *AppError {
	if reason != "" {
		return New(TypeInvalidCredentials, "登录失败: "+reason)
	}
	return ErrInvalidCredentials
}

func NewPermissionError(resource string) *AppError {
	return New(TypePermissionDenied, "无权限访问: "+resource)
}
