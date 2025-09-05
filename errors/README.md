# Errors 包使用指南

## 概述

`errors` 包是 YaBase 项目的统一错误处理模块，提供了结构化的错误类型定义和便捷的错误处理函数。该包遵循 Go 语言的错误处理最佳实践，支持错误包装、类型检查和错误转换。

## 核心结构

### AppError 类型

```go
type AppError struct {
    Type    string `json:"type"`    // 错误类型标识
    Message string `json:"message"` // 错误消息
    Cause   error  `json:"-"`       // 原始错误（不序列化）
}
```

## 基础用法

### 1. 创建新错误

```go
// 使用预定义错误
err := errors.ErrNotFound

// 创建自定义错误
err := errors.New("CUSTOM_ERROR", "自定义错误消息")

// 包装现有错误
err := errors.Wrap(originalErr, "WRAP_ERROR", "包装后的错误消息")
```

### 2. 错误类型检查

```go
// 检查是否为特定类型的错误
if errors.Is(err, errors.ErrNotFound) {
    // 处理未找到错误
}

// 获取错误类型
errorType := errors.GetType(err)
```

### 3. 错误信息获取

```go
// 获取错误消息
message := err.Error()

// 获取原始错误
if appErr, ok := err.(*errors.AppError); ok {
    originalErr := appErr.Cause
}
```

## 错误类型分类

### 通用错误 (common.go)

| 错误类型 | 常量 | 预定义实例 | 描述 |
|---------|------|-----------|------|
| 记录不存在 | `TypeNotFound` | `ErrNotFound` | 数据库记录未找到 |
| 数据重复 | `TypeDuplicateKey` | `ErrDuplicateKey` | 唯一约束冲突 |
| 数据库错误 | `TypeDatabaseError` | `ErrDatabaseError` | 数据库操作失败 |
| 内部错误 | `TypeInternalError` | `ErrInternalError` | 系统内部错误 |
| 验证错误 | `TypeValidation` | `ErrValidation` | 数据验证失败 |
| 输入错误 | `TypeInvalidInput` | `ErrInvalidInput` | 输入参数错误 |

### 认证授权错误 (errors_auth.go)

| 错误类型 | 常量 | 预定义实例 | 描述 |
|---------|------|-----------|------|
| 凭据无效 | `TypeInvalidCredentials` | `ErrInvalidCredentials` | 用户名或密码错误 |
| 令牌过期 | `TypeTokenExpired` | `ErrTokenExpired` | 登录已过期 |
| 令牌无效 | `TypeTokenInvalid` | `ErrTokenInvalid` | 无效的登录凭证 |
| 未授权 | `TypeUnauthorized` | `ErrUnauthorized` | 请先登录 |
| 权限不足 | `TypePermissionDenied` | `ErrPermissionDenied` | 权限不足 |

### 用户管理错误 (errors_user.go)

| 错误类型 | 常量 | 预定义实例 | 描述 |
|---------|------|-----------|------|
| 用户不存在 | `TypeUserNotFound` | `ErrUserNotFound` | 用户不存在 |
| 用户已存在 | `TypeUserExists` | `ErrUserExists` | 用户已存在 |
| 用户被禁用 | `TypeUserDisabled` | `ErrUserDisabled` | 账户已被禁用 |
| 操作失败 | `TypeOperationFailed` | `ErrCreateFailed` | 创建/更新/删除失败 |

### OAuth 错误 (errors_oauth.go)

| 错误类型 | 常量 | 预定义实例 | 描述 |
|---------|------|-----------|------|
| 客户端无效 | `TypeInvalidClient` | `ErrInvalidClient` | 客户端认证失败 |
| 授权类型不支持 | `TypeUnsupportedGrantType` | `ErrUnsupportedGrantType` | 不支持的授权类型 |
| 授权码无效 | `TypeInvalidGrant` | `ErrInvalidGrant` | 授权码无效或已过期 |
| 权限范围无效 | `TypeInvalidScope` | `ErrInvalidScope` | 请求的权限范围无效 |

## 使用示例

### 1. 数据库操作错误处理

```go
func GetUserByID(id uint) (*User, error) {
    var user User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, errors.FromDatabaseError(err)
    }
    return &user, nil
}

// 使用示例
user, err := GetUserByID(123)
if err != nil {
    if errors.Is(err, errors.ErrNotFound) {
        // 处理用户不存在的情况
        return c.JSON(404, gin.H{"error": "用户不存在"})
    }
    // 处理其他数据库错误
    return c.JSON(500, gin.H{"error": "服务器内部错误"})
}
```

### 2. 用户认证错误处理

```go
func Login(username, password string) (*User, error) {
    var user User
    err := db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.IsNotFoundError(err) {
            return nil, errors.ErrInvalidCredentials
        }
        return nil, errors.FromDatabaseError(err)
    }
    
    if !checkPassword(password, user.Password) {
        return nil, errors.NewInvalidCredentialsError("密码错误")
    }
    
    if user.Status == "disabled" {
        return nil, errors.ErrUserDisabled
    }
    
    return &user, nil
}

// 使用示例
user, err := Login("admin", "password")
if err != nil {
    switch errors.GetType(err) {
    case errors.TypeInvalidCredentials:
        return c.JSON(401, gin.H{"error": "用户名或密码错误"})
    case errors.TypeUserDisabled:
        return c.JSON(403, gin.H{"error": "账户已被禁用"})
    default:
        return c.JSON(500, gin.H{"error": "登录失败"})
    }
}
```

### 3. 数据验证错误处理

```go
func CreateUser(userData CreateUserRequest) (*User, error) {
    // 验证用户名
    if userData.Username == "" {
        return nil, errors.NewValidationError("username", "用户名不能为空")
    }
    
    // 检查用户名是否已存在
    var existingUser User
    err := db.Where("username = ?", userData.Username).First(&existingUser).Error
    if err == nil {
        return nil, errors.ErrUsernameExists
    }
    if !errors.IsNotFoundError(err) {
        return nil, errors.FromDatabaseError(err)
    }
    
    // 创建用户
    user := &User{
        Username: userData.Username,
        Email:    userData.Email,
    }
    
    err = db.Create(user).Error
    if err != nil {
        if errors.IsDuplicateError(err) {
            return nil, errors.ErrUserExists
        }
        return nil, errors.FromDatabaseError(err)
    }
    
    return user, nil
}
```

### 4. 错误包装和链式处理

```go
func ProcessUserData(userID uint, data UserData) error {
    // 获取用户
    user, err := GetUserByID(userID)
    if err != nil {
        return errors.Wrap(err, "USER_PROCESS_FAILED", "获取用户信息失败")
    }
    
    // 验证数据
    if err := validateUserData(data); err != nil {
        return errors.Wrap(err, "VALIDATION_FAILED", "用户数据验证失败")
    }
    
    // 更新用户
    err = db.Model(user).Updates(data).Error
    if err != nil {
        return errors.Wrap(err, "UPDATE_FAILED", "更新用户数据失败")
    }
    
    return nil
}

// 使用示例
err := ProcessUserData(123, userData)
if err != nil {
    // 获取错误类型
    errorType := errors.GetType(err)
    
    // 检查是否为特定错误
    if errors.Is(err, errors.ErrNotFound) {
        // 处理用户不存在
    }
    
    // 记录错误日志
    log.Printf("处理用户数据失败: %v", err)
}
```

### 5. 中间件中的错误处理

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            var appErr *errors.AppError
            if errors.As(err.Err, &appErr) {
                // 根据错误类型返回不同的HTTP状态码
                switch appErr.Type {
                case errors.TypeNotFound:
                    c.JSON(404, gin.H{"error": appErr.Message})
                case errors.TypeUnauthorized:
                    c.JSON(401, gin.H{"error": appErr.Message})
                case errors.TypePermissionDenied:
                    c.JSON(403, gin.H{"error": appErr.Message})
                case errors.TypeValidation:
                    c.JSON(400, gin.H{"error": appErr.Message})
                default:
                    c.JSON(500, gin.H{"error": "服务器内部错误"})
                }
            } else {
                c.JSON(500, gin.H{"error": "服务器内部错误"})
            }
        }
    }
}
```

## 最佳实践

### 1. 错误类型选择

- 使用预定义的错误类型和实例，保持一致性
- 为特定业务场景创建专门的错误类型
- 避免创建过于通用的错误类型

### 2. 错误包装

- 在调用链中适当包装错误，保留原始错误信息
- 使用有意义的错误消息，便于调试
- 避免过度包装，保持错误信息的简洁性

### 3. 错误处理

- 在适当的层级处理错误，不要忽略错误
- 使用 `errors.Is()` 和 `errors.GetType()` 进行错误类型检查
- 为不同的错误类型提供不同的处理逻辑

### 4. 日志记录

- 记录错误的完整信息，包括原始错误
- 使用结构化日志，便于错误追踪
- 区分可恢复和不可恢复的错误

## 扩展指南

### 添加新的错误类型

1. 在相应的文件中定义错误类型常量
2. 创建预定义的错误实例
3. 添加必要的工厂函数
4. 更新文档

### 示例：添加文件操作错误

```go
// 在 errors_file.go 中
const (
    TypeFileNotFound    = "FILE_NOT_FOUND"
    TypeFileReadError   = "FILE_READ_ERROR"
    TypeFileWriteError  = "FILE_WRITE_ERROR"
)

var (
    ErrFileNotFound   = New(TypeFileNotFound, "文件不存在")
    ErrFileReadError  = New(TypeFileReadError, "文件读取失败")
    ErrFileWriteError = New(TypeFileWriteError, "文件写入失败")
)

func NewFileNotFoundError(filename string) *AppError {
    return New(TypeFileNotFound, fmt.Sprintf("文件不存在: %s", filename))
}
```

这个错误处理包为 YaBase 项目提供了统一、结构化的错误处理机制，有助于提高代码的可维护性和调试效率。
