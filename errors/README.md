# Errors 包使用指南

## 概述

`errors` 包是 YaBase 项目的统一错误处理模块，提供了基于链式 API 的现代化错误构建系统。通过流畅的接口设计，支持错误包装、上下文字段和类型检查。

## 核心理念

- **链式调用**：使用流畅的 Builder 模式构建错误
- **类型安全**：预定义的错误类型保证一致性
- **丰富上下文**：支持附加任意字段用于调试和日志记录
- **错误包装**：完整保留原始错误信息

## 核心结构

### AppError 类型

```go
type AppError struct {
    Type    string                 `json:"type"`             // 错误类型标识
    Message string                 `json:"message"`          // 错误消息
    Cause   error                  `json:"-"`                // 原始错误（不序列化）
    Fields  map[string]interface{} `json:"fields,omitempty"` // 上下文字段
}
```

## 快速开始

### 1. 基础用法

```go
// 简单创建错误
func GetUser(id uint) (*User, error) {
    var user User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, errors.NotFound().
            WithMessage("用户不存在").
            Build()
    }
    return &user, nil
}
```

### 2. 添加上下文字段

```go
// 添加调试信息
func GetUser(id uint) (*User, error) {
    var user User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, errors.NotFound().
            WithMessage("用户不存在").
            WithField("user_id", id).
            WithField("operation", "GetUser").
            Build()
    }
    return &user, nil
}
```

### 3. 包装现有错误

```go
// 保留原始错误信息
func GetUser(id uint) (*User, error) {
    var user User
    err := db.First(&user, id).Error
    if err != nil {
        return nil, errors.Database().
            WithMessage("查询用户失败").
            WithCause(err).
            WithField("user_id", id).
            Build()
    }
    return &user, nil
}
```

### 4. 批量添加字段

```go
func ProcessPayment(orderID string, amount float64, userID uint) error {
    err := payment.Charge(orderID, amount)
    if err != nil {
        return errors.Internal().
            WithMessage("支付处理失败").
            WithCause(err).
            WithFields(map[string]interface{}{
                "order_id": orderID,
                "amount":   amount,
                "user_id":  userID,
                "timestamp": time.Now(),
            }).
            Build()
    }
    return nil
}
```

## 标准库函数重新导出

为了方便使用，本包重新导出了 Go 标准库 `errors` 包的常用函数，这样项目中只需要 `import "github.com/3086953492/YaBase/errors"` 一个包即可。

### 可用的标准库函数

| 函数 | 说明 | 使用场景 |
|------|------|---------|
| `NewSimple(text string) error` | 创建简单错误 | 不需要结构化信息的简单错误 |
| `As(err error, target any) bool` | 类型断言 | 从错误链中提取特定类型的错误 |
| `Is(err, target error) bool` | 错误比较 | 检查错误链中是否包含目标错误 |
| `Unwrap(err error) error` | 解包错误 | 获取错误链中的下一个错误 |
| `Join(errs ...error) error` | 合并错误 | 将多个错误组合成一个错误（Go 1.20+）|

### 标准库函数使用示例

```go
import (
    "io"
    "github.com/3086953492/YaBase/errors"
)

// 创建简单错误
err := errors.NewSimple("配置文件不存在")

// 错误比较（标准库 Is）
if errors.Is(err, io.EOF) {
    // 处理 EOF 错误
}

// 类型断言（标准库 As）
var pathErr *fs.PathError
if errors.As(err, &pathErr) {
    fmt.Printf("路径错误: %s\n", pathErr.Path)
}

// 解包错误
unwrapped := errors.Unwrap(err)

// 合并多个错误
multiErr := errors.Join(err1, err2, err3)
```

### AppError 类型检查

对于我们自定义的 `AppError` 类型，使用 `IsAppError()` 函数进行类型检查：

```go
// 创建 AppError
notFoundErr := errors.NotFound().Build()

// 检查是否为特定类型的 AppError
if errors.IsAppError(err, notFoundErr) {
    // 这是一个 NotFound 类型的 AppError
}

// 或者直接比较错误类型字符串
if errors.GetType(err) == errors.TypeNotFound {
    // 处理 NotFound 错误
}
```

## 预定义错误类型

### 错误构造器

| 构造器 | HTTP 状态码 | 错误类型常量 | 默认消息 | 使用场景 |
|--------|------------|-------------|---------|---------|
| `NotFound()` | 404 | `TypeNotFound` | "记录不存在" | 数据库记录未找到 |
| `InvalidInput()` | 400 | `TypeInvalidInput` | "输入参数错误" | 请求参数不合法 |
| `Unauthorized()` | 401 | `TypeUnauthorized` | "未授权" | 未登录或令牌无效 |
| `Forbidden()` | 403 | `TypeForbidden` | "权限不足" | 已登录但无权限 |
| `Duplicate()` | 409 | `TypeDuplicate` | "数据已存在" | 唯一约束冲突 |
| `Internal()` | 500 | `TypeInternal` | "服务器内部错误" | 系统内部错误 |
| `Database()` | 500 | `TypeDatabase` | "数据库操作失败" | 数据库操作错误 |
| `Validation()` | 422 | `TypeValidation` | "数据验证失败" | 数据验证不通过 |

### 链式方法

| 方法 | 说明 | 示例 |
|------|------|------|
| `WithMessage(msg string)` | 自定义错误消息 | `.WithMessage("用户不存在")` |
| `WithCause(err error)` | 包装原始错误 | `.WithCause(dbErr)` |
| `WithField(key, value)` | 添加单个上下文字段 | `.WithField("user_id", 123)` |
| `WithFields(map[string]interface{})` | 批量添加字段 | `.WithFields(contextData)` |
| `Build()` | 构建最终错误 | `.Build()` |

## 完整使用示例

### 1. 用户认证场景

```go
func Login(username, password string) (*User, error) {
    var user User
    err := db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.IsNotFoundError(err) {
            return nil, errors.Unauthorized().
                WithMessage("用户名或密码错误").
                Build()
        }
        return nil, errors.Database().
            WithMessage("查询用户失败").
            WithCause(err).
            Build()
    }
    
    if !checkPassword(password, user.Password) {
        return nil, errors.Unauthorized().
            WithMessage("用户名或密码错误").
            WithField("username", username).
            Build()
    }
    
    if user.Status == "disabled" {
        return nil, errors.Forbidden().
            WithMessage("账户已被禁用").
            WithField("user_id", user.ID).
            WithField("username", username).
            Build()
    }
    
    return &user, nil
}
```

### 2. 数据创建场景

```go
func CreateUser(req CreateUserRequest) (*User, error) {
    // 验证输入
    if req.Username == "" {
        return nil, errors.InvalidInput().
            WithMessage("用户名不能为空").
            WithField("field", "username").
            Build()
    }
    
    // 检查用户名是否存在
    var existing User
    err := db.Where("username = ?", req.Username).First(&existing).Error
    if err == nil {
        return nil, errors.Duplicate().
            WithMessage("用户名已存在").
            WithField("username", req.Username).
            Build()
    }
    if !errors.IsNotFoundError(err) {
        return nil, errors.Database().
            WithMessage("检查用户名失败").
            WithCause(err).
            Build()
    }
    
    // 创建用户
    user := &User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    err = db.Create(user).Error
    if err != nil {
        return nil, errors.Database().
            WithMessage("创建用户失败").
            WithCause(err).
            WithField("username", req.Username).
            Build()
    }
    
    return user, nil
}
```

### 3. 权限检查场景

```go
func DeleteUser(operatorID, targetUserID uint) error {
    // 检查操作者权限
    var operator User
    err := db.First(&operator, operatorID).Error
    if err != nil {
        return errors.NotFound().
            WithMessage("操作者不存在").
            WithField("operator_id", operatorID).
            Build()
    }
    
    if operator.Role != "admin" {
        return errors.Forbidden().
            WithMessage("只有管理员可以删除用户").
            WithField("operator_id", operatorID).
            WithField("operator_role", operator.Role).
            Build()
    }
    
    // 执行删除
    err = db.Delete(&User{}, targetUserID).Error
    if err != nil {
        return errors.Database().
            WithMessage("删除用户失败").
            WithCause(err).
            WithField("target_user_id", targetUserID).
            Build()
    }
    
    return nil
}
```

### 4. 数据库错误自动转换

```go
func GetUserByID(id uint) (*User, error) {
    var user User
    err := db.First(&user, id).Error
    if err != nil {
        // 自动将 GORM 错误转换为 AppError
        return nil, errors.FromDatabaseError(err)
    }
    return &user, nil
}
```

## 错误处理最佳实践

### 1. 在 Gin 中间件中统一处理

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var appErr *errors.AppError
            if errors.As(err, &appErr) {
                // 根据错误类型返回相应的 HTTP 状态码
                statusCode := getHTTPStatusCode(appErr.Type)
                
                response := gin.H{
                    "error": appErr.Message,
                    "type":  appErr.Type,
                }
                
                // 在开发环境返回详细信息
                if config.IsDevelopment() {
                    response["fields"] = appErr.Fields
                    if appErr.Cause != nil {
                        response["cause"] = appErr.Cause.Error()
                    }
                }
                
                c.JSON(statusCode, response)
            } else {
                c.JSON(500, gin.H{"error": "服务器内部错误"})
            }
        }
    }
}

func getHTTPStatusCode(errType string) int {
    switch errType {
    case errors.TypeNotFound:
        return 404
    case errors.TypeInvalidInput:
        return 400
    case errors.TypeUnauthorized:
        return 401
    case errors.TypeForbidden:
        return 403
    case errors.TypeDuplicate:
        return 409
    case errors.TypeValidation:
        return 422
    default:
        return 500
    }
}
```

### 2. 日志记录

```go
func LogError(err error) {
    if appErr, ok := err.(*errors.AppError); ok {
        // 结构化日志
        logger.Error("应用错误",
            zap.String("type", appErr.Type),
            zap.String("message", appErr.Message),
            zap.Any("fields", appErr.Fields),
            zap.Error(appErr.Cause),
        )
    } else {
        logger.Error("未知错误", zap.Error(err))
    }
}
```

### 3. 错误类型检查

```go
func HandleError(err error) {
    // 获取错误类型
    errType := errors.GetType(err)
    
    switch errType {
    case errors.TypeNotFound:
        // 处理未找到错误
    case errors.TypeUnauthorized:
        // 处理未授权错误
    case errors.TypeForbidden:
        // 处理权限不足错误
    default:
        // 处理其他错误
    }
}
```

### 4. 获取错误上下文

```go
func ProcessError(err error) {
    // 获取所有字段
    fields := errors.GetFields(err)
    
    // 或者获取特定字段
    if appErr, ok := err.(*errors.AppError); ok {
        if userID, ok := appErr.GetField("user_id"); ok {
            log.Printf("错误涉及用户: %v", userID)
        }
    }
}
```

## 工具函数

### 错误类型检查

```go
// 检查是否为 GORM 的 RecordNotFound 错误
errors.IsNotFoundError(err)

// 检查是否为数据库唯一约束错误
errors.IsDuplicateError(err)

// 获取错误类型字符串
errors.GetType(err)

// 获取错误的所有上下文字段
errors.GetFields(err)
```

### 数据库错误转换

```go
// 自动将数据库错误转换为相应的 AppError
appErr := errors.FromDatabaseError(dbErr)
// - gorm.ErrRecordNotFound -> NotFound
// - 唯一约束错误 -> Duplicate
// - 其他数据库错误 -> Database
```

## 迁移指南

### 从旧 API 迁移

#### 旧方式
```go
// 使用预定义错误
return errors.ErrUserNotFound

// 使用工厂函数
return errors.NewUserNotFoundError(userID)

// 创建新错误
return errors.New("USER_NOT_FOUND", "用户不存在")

// 包装错误
return errors.Wrap(err, "DATABASE_ERROR", "数据库操作失败")
```

#### 新方式
```go
// 基础用法
return errors.NotFound().
    WithMessage("用户不存在").
    Build()

// 带上下文
return errors.NotFound().
    WithMessage("用户不存在").
    WithField("user_id", userID).
    Build()

// 包装错误
return errors.Database().
    WithMessage("数据库操作失败").
    WithCause(err).
    Build()
```

## 设计原则

1. **明确性优于简洁性**：宁愿多写几行代码，也要清楚表达错误信息
2. **上下文优先**：尽可能添加有助于调试的上下文字段
3. **保留原始错误**：使用 `WithCause()` 包装错误，不要丢弃原始信息
4. **统一处理**：在中间件层统一处理错误到 HTTP 响应的转换
5. **结构化日志**：充分利用错误的结构化信息进行日志记录

## 常见模式

### 模式 1：简单场景（使用默认消息）

```go
if user == nil {
    return errors.NotFound().Build()
}
```

### 模式 2：自定义消息

```go
if user == nil {
    return errors.NotFound().
        WithMessage("指定用户不存在").
        Build()
}
```

### 模式 3：带关键字段

```go
if user == nil {
    return errors.NotFound().
        WithMessage("用户不存在").
        WithField("user_id", userID).
        Build()
}
```

### 模式 4：包装数据库错误

```go
err := db.First(&user, id).Error
if err != nil {
    return errors.FromDatabaseError(err)
}
```

### 模式 5：完整错误信息

```go
return errors.Internal().
    WithMessage("处理订单失败").
    WithCause(originalErr).
    WithFields(map[string]interface{}{
        "order_id":   orderID,
        "user_id":    userID,
        "step":       "payment",
        "timestamp":  time.Now(),
    }).
    Build()
```

## 总结

新的错误处理包提供了：

- ✅ **简洁的 API**：链式调用，清晰易读
- ✅ **类型安全**：预定义的错误类型
- ✅ **丰富上下文**：支持任意字段
- ✅ **完整追踪**：保留原始错误链
- ✅ **统一处理**：便于中间件和日志集成

遵循本指南，可以构建健壮、易于调试的错误处理体系。
