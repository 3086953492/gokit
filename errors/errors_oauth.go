package errors

const (
	TypeInvalidClient           = "INVALID_CLIENT"
	TypeUnsupportedGrantType    = "UNSUPPORTED_GRANT_TYPE"
	TypeInvalidGrant            = "INVALID_GRANT"
	TypeInvalidScope            = "INVALID_SCOPE"
	TypeInvalidRedirectURI      = "INVALID_REDIRECT_URI"
	TypeAccessDenied            = "ACCESS_DENIED"
	TypeUnsupportedResponseType = "UNSUPPORTED_RESPONSE_TYPE"
	TypeInvalidToken            = "INVALID_TOKEN"
	TypeOAuthTokenExpired       = "TOKEN_EXPIRED"
	TypeInsufficientScope       = "INSUFFICIENT_SCOPE"
	TypeClientNotFound          = "CLIENT_NOT_FOUND"
	TypeUpdateClientFailed      = "UPDATE_CLIENT_FAILED"
	TypeDeleteClientFailed      = "DELETE_CLIENT_FAILED"
	TypeInvalidCodeOrState      = "INVALID_CODE_OR_STATE"
	TypeUnauthorized            = "UNAUTHORIZED"
)

var (
	ErrInvalidClient           = New(TypeInvalidClient, "客户端认证失败")
	ErrUnsupportedGrantType    = New(TypeUnsupportedGrantType, "不支持的授权类型")
	ErrInvalidGrant            = New(TypeInvalidGrant, "授权码无效或已过期")
	ErrInvalidScope            = New(TypeInvalidScope, "请求的权限范围无效")
	ErrInvalidRedirectURI      = New(TypeInvalidRedirectURI, "重定向URI无效")
	ErrAccessDenied            = New(TypeAccessDenied, "用户拒绝授权")
	ErrUnsupportedResponseType = New(TypeUnsupportedResponseType, "不支持的响应类型")
	ErrInvalidToken            = New(TypeInvalidToken, "访问令牌无效")
	ErrOAuthTokenExpired       = New(TypeOAuthTokenExpired, "访问令牌已过期")
	ErrInsufficientScope       = New(TypeInsufficientScope, "权限范围不足")
	ErrClientNotFound          = New(TypeClientNotFound, "客户端不存在")
	ErrUpdateClientFailed      = New(TypeUpdateClientFailed, "更新客户端失败")
	ErrDeleteClientFailed      = New(TypeDeleteClientFailed, "删除客户端失败")
	ErrInvalidCodeOrState      = New(TypeInvalidCodeOrState, "授权码或状态无效")
	ErrUnauthorized            = New(TypeUnauthorized, "未授权")
)
